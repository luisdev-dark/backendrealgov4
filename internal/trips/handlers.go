package trips

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luis/ride-mvp/internal/db"
	"github.com/luis/ride-mvp/internal/httpx"
	"github.com/luis/ride-mvp/internal/models"
)

// HardcodedPassengerID is the fixed user ID for MVP testing
const HardcodedPassengerID = "11111111-1111-1111-1111-111111111111"

// CreateTrip handles POST /api/trips
// Creates a new trip request with hardcoded passenger ID
func CreateTrip(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.RouteID == "" {
		httpx.Error(w, http.StatusBadRequest, "route_id is required")
		return
	}

	// Validate payment method
	validPaymentMethods := map[string]bool{
		"cash": true,
		"yape": true,
		"plin": true,
	}
	if !validPaymentMethods[req.PaymentMethod] {
		httpx.Error(w, http.StatusBadRequest, "payment_method must be one of: cash, yape, plin")
		return
	}

	// Validate pickup != dropoff if both are provided
	if req.PickupStopID != nil && req.DropoffStopID != nil {
		if *req.PickupStopID == *req.DropoffStopID {
			httpx.Error(w, http.StatusBadRequest, "pickup_stop_id and dropoff_stop_id must be different")
			return
		}
	}

	ctx := r.Context()
	pool, err := db.GetPool(ctx)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Database connection failed")
		return
	}

	// Verify route exists and get base price
	var basePriceCents int
	var currency string
	routeQuery := `
		SELECT base_price_cents, currency
		FROM app.routes
		WHERE id = $1 AND is_active = true
	`
	err = pool.QueryRow(ctx, routeQuery, req.RouteID).Scan(&basePriceCents, &currency)
	if err != nil {
		if err.Error() == "no rows in result set" {
			httpx.Error(w, http.StatusNotFound, "Route not found")
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "Failed to verify route")
		return
	}

	// Generate new trip ID
	tripID := uuid.New().String()

	// Insert trip
	insertQuery := `
		INSERT INTO app.trips (
			id, route_id, passenger_id, pickup_stop_id, dropoff_stop_id,
			status, payment_method, price_cents, currency
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, status, created_at
	`

	var trip models.Trip
	err = pool.QueryRow(
		ctx,
		insertQuery,
		tripID,
		req.RouteID,
		HardcodedPassengerID,
		req.PickupStopID,
		req.DropoffStopID,
		"requested",
		req.PaymentMethod,
		basePriceCents,
		currency,
	).Scan(&trip.ID, &trip.Status, &trip.CreatedAt)

	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create trip")
		return
	}

	// Build response
	trip.RouteID = req.RouteID
	trip.PassengerID = HardcodedPassengerID
	trip.PickupStopID = req.PickupStopID
	trip.DropoffStopID = req.DropoffStopID
	trip.PaymentMethod = req.PaymentMethod
	trip.PriceCents = basePriceCents
	trip.Currency = currency

	httpx.JSON(w, http.StatusCreated, trip)
}

// GetTripByID handles GET /api/trips/{id}
// Returns trip details by ID
func GetTripByID(w http.ResponseWriter, r *http.Request) {
	tripID := chi.URLParam(r, "id")
	if tripID == "" {
		httpx.Error(w, http.StatusBadRequest, "Trip ID is required")
		return
	}

	ctx := r.Context()
	pool, err := db.GetPool(ctx)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Database connection failed")
		return
	}

	query := `
		SELECT id, route_id, passenger_id, pickup_stop_id, dropoff_stop_id,
		       status, payment_method, price_cents, currency, created_at, updated_at
		FROM app.trips
		WHERE id = $1
	`

	var trip models.Trip
	err = pool.QueryRow(ctx, query, tripID).Scan(
		&trip.ID,
		&trip.RouteID,
		&trip.PassengerID,
		&trip.PickupStopID,
		&trip.DropoffStopID,
		&trip.Status,
		&trip.PaymentMethod,
		&trip.PriceCents,
		&trip.Currency,
		&trip.CreatedAt,
		&trip.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			httpx.Error(w, http.StatusNotFound, "Trip not found")
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "Failed to fetch trip")
		return
	}

	httpx.JSON(w, http.StatusOK, trip)
}

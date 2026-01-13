package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luis/ride-mvp/internal/db"
	"github.com/luis/ride-mvp/internal/httpx"
	"github.com/luis/ride-mvp/internal/models"
)

// GetRoutes handles GET /api/routes
// Returns a list of all active routes
func GetRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	pool, err := db.GetPool(ctx)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Database connection failed")
		return
	}

	query := `
		SELECT id, name, origin_name, destination_name, base_price_cents, currency
		FROM app.routes
		WHERE is_active = true
		ORDER BY name
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to fetch routes")
		return
	}
	defer rows.Close()

	var routes []models.RouteSummary
	for rows.Next() {
		var route models.RouteSummary
		err := rows.Scan(
			&route.ID,
			&route.Name,
			&route.OriginName,
			&route.DestinationName,
			&route.BasePriceCents,
			&route.Currency,
		)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "Failed to parse route data")
			return
		}
		routes = append(routes, route)
	}

	if err := rows.Err(); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Error reading routes")
		return
	}

	// Return empty array if no routes found
	if routes == nil {
		routes = []models.RouteSummary{}
	}

	httpx.JSON(w, http.StatusOK, routes)
}

// GetRouteByID handles GET /api/routes/{id}
// Returns route details with all stops ordered by stop_order
func GetRouteByID(w http.ResponseWriter, r *http.Request) {
	routeID := chi.URLParam(r, "id")
	if routeID == "" {
		httpx.Error(w, http.StatusBadRequest, "Route ID is required")
		return
	}

	ctx := context.Background()
	pool, err := db.GetPool(ctx)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Database connection failed")
		return
	}

	// Fetch route details
	routeQuery := `
		SELECT id, name, origin_name, origin_lat, origin_lon,
		       destination_name, destination_lat, destination_lon,
		       base_price_cents, currency, is_active, created_at, updated_at
		FROM app.routes
		WHERE id = $1 AND is_active = true
	`

	var route models.Route
	err = pool.QueryRow(ctx, routeQuery, routeID).Scan(
		&route.ID,
		&route.Name,
		&route.OriginName,
		&route.OriginLat,
		&route.OriginLon,
		&route.DestinationName,
		&route.DestinationLat,
		&route.DestinationLon,
		&route.BasePriceCents,
		&route.Currency,
		&route.IsActive,
		&route.CreatedAt,
		&route.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			httpx.Error(w, http.StatusNotFound, "Route not found")
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "Failed to fetch route")
		return
	}

	// Fetch stops for this route
	stopsQuery := `
		SELECT id, name, stop_order
		FROM app.route_stops
		WHERE route_id = $1 AND is_active = true
		ORDER BY stop_order
	`

	rows, err := pool.Query(ctx, stopsQuery, routeID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to fetch stops")
		return
	}
	defer rows.Close()

	var stops []models.RouteStopSummary
	for rows.Next() {
		var stop models.RouteStopSummary
		err := rows.Scan(&stop.ID, &stop.Name, &stop.StopOrder)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "Failed to parse stop data")
			return
		}
		stops = append(stops, stop)
	}

	if err := rows.Err(); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Error reading stops")
		return
	}

	// Return empty array if no stops found
	if stops == nil {
		stops = []models.RouteStopSummary{}
	}

	response := models.RouteDetailResponse{
		Route: route,
		Stops: stops,
	}

	httpx.JSON(w, http.StatusOK, response)
}

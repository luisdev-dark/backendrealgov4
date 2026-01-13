package models

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	FullName  string    `json:"full_name"`
	PhoneE164 string    `json:"phone_e164"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Route represents a transportation route from origin to destination
type Route struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	OriginName      string    `json:"origin_name"`
	OriginLat       float64   `json:"origin_lat"`
	OriginLon       float64   `json:"origin_lon"`
	DestinationName string    `json:"destination_name"`
	DestinationLat  float64   `json:"destination_lat"`
	DestinationLon  float64   `json:"destination_lon"`
	BasePriceCents  int       `json:"base_price_cents"`
	Currency        string    `json:"currency"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RouteSummary represents a simplified route for listing
type RouteSummary struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	OriginName      string `json:"origin_name"`
	DestinationName string `json:"destination_name"`
	BasePriceCents  int    `json:"base_price_cents"`
	Currency        string `json:"currency"`
}

// RouteStop represents a stop along a route
type RouteStop struct {
	ID        string    `json:"id"`
	RouteID   string    `json:"route_id"`
	StopOrder int       `json:"stop_order"`
	Name      string    `json:"name"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RouteStopSummary represents a simplified stop for route details
type RouteStopSummary struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StopOrder int    `json:"stop_order"`
}

// Trip represents a passenger trip request
type Trip struct {
	ID            string    `json:"id"`
	RouteID       string    `json:"route_id"`
	PassengerID   string    `json:"passenger_id"`
	PickupStopID  *string   `json:"pickup_stop_id"`
	DropoffStopID *string   `json:"dropoff_stop_id"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	PriceCents    int       `json:"price_cents"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateTripRequest represents the request body for creating a trip
type CreateTripRequest struct {
	RouteID       string  `json:"route_id"`
	PickupStopID  *string `json:"pickup_stop_id"`
	DropoffStopID *string `json:"dropoff_stop_id"`
	PaymentMethod string  `json:"payment_method"`
}

// RouteDetailResponse represents the response for route details with stops
type RouteDetailResponse struct {
	Route Route              `json:"route"`
	Stops []RouteStopSummary `json:"stops"`
}

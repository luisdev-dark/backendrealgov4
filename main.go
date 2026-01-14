package main

import (
	"net/http"

	"api/internal/httpx"
	"api/internal/routes"
	"api/internal/trips"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Handler is the main entry point for Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	// Create router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "realgo-mvp",
		})
	})

	// Routes endpoints
	router.Get("/routes", routes.GetRoutes)
	router.Get("/routes/{id}", routes.GetRouteByID)

	// Trips endpoints
	router.Post("/trips", trips.CreateTrip)
	router.Get("/trips/{id}", trips.GetTripByID)

	// Serve the request
	router.ServeHTTP(w, r)
}

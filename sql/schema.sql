-- RealGo MVP Database Schema for PostgreSQL (Neon)
-- Simple schema with app namespace, no RLS

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create app schema
CREATE SCHEMA IF NOT EXISTS app;

-- User roles enum
CREATE TYPE app.user_role AS ENUM ('passenger', 'driver', 'admin');

-- Trip status enum
CREATE TYPE app.trip_status AS ENUM ('requested', 'accepted', 'in_progress', 'completed', 'cancelled');

-- Payment method enum
CREATE TYPE app.payment_method AS ENUM ('cash', 'yape', 'plin');

-- Users table
CREATE TABLE app.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role app.user_role NOT NULL DEFAULT 'passenger',
    full_name VARCHAR(255) NOT NULL,
    phone_e164 VARCHAR(20) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Routes table
CREATE TABLE app.routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    origin_name VARCHAR(255) NOT NULL,
    origin_lat DECIMAL(10, 8) NOT NULL,
    origin_lon DECIMAL(11, 8) NOT NULL,
    destination_name VARCHAR(255) NOT NULL,
    destination_lat DECIMAL(10, 8) NOT NULL,
    destination_lon DECIMAL(11, 8) NOT NULL,
    base_price_cents INTEGER NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'PEN',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Route stops table
CREATE TABLE app.route_stops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES app.routes(id) ON DELETE CASCADE,
    stop_order INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    lat DECIMAL(10, 8) NOT NULL,
    lon DECIMAL(11, 8) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_route_stop_order UNIQUE (route_id, stop_order)
);

-- Trips table
CREATE TABLE app.trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES app.routes(id),
    passenger_id UUID NOT NULL REFERENCES app.users(id),
    pickup_stop_id UUID REFERENCES app.route_stops(id),
    dropoff_stop_id UUID REFERENCES app.route_stops(id),
    status app.trip_status NOT NULL DEFAULT 'requested',
    payment_method app.payment_method NOT NULL,
    price_cents INTEGER NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'PEN',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT different_pickup_dropoff CHECK (
        pickup_stop_id IS NULL OR 
        dropoff_stop_id IS NULL OR 
        pickup_stop_id != dropoff_stop_id
    )
);

-- Indexes for better query performance
CREATE INDEX idx_users_phone ON app.users(phone_e164);
CREATE INDEX idx_routes_active ON app.routes(is_active);
CREATE INDEX idx_route_stops_route_id ON app.route_stops(route_id);
CREATE INDEX idx_route_stops_order ON app.route_stops(route_id, stop_order);
CREATE INDEX idx_trips_passenger ON app.trips(passenger_id);
CREATE INDEX idx_trips_route ON app.trips(route_id);
CREATE INDEX idx_trips_status ON app.trips(status);
CREATE INDEX idx_trips_created_at ON app.trips(created_at DESC);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION app.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON app.users
    FOR EACH ROW
    EXECUTE FUNCTION app.update_updated_at_column();

CREATE TRIGGER update_routes_updated_at
    BEFORE UPDATE ON app.routes
    FOR EACH ROW
    EXECUTE FUNCTION app.update_updated_at_column();

CREATE TRIGGER update_route_stops_updated_at
    BEFORE UPDATE ON app.route_stops
    FOR EACH ROW
    EXECUTE FUNCTION app.update_updated_at_column();

CREATE TRIGGER update_trips_updated_at
    BEFORE UPDATE ON app.trips
    FOR EACH ROW
    EXECUTE FUNCTION app.update_updated_at_column();

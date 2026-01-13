-- RealGo MVP Seed Data
-- Fixed UUIDs for testing

-- Insert test passenger user
-- UUID: 11111111-1111-1111-1111-111111111111
INSERT INTO app.users (id, role, full_name, phone_e164) VALUES
('11111111-1111-1111-1111-111111111111', 'passenger', 'Test Passenger', '+51987654321');

-- Insert routes
-- Route 1 UUID: 22222222-2222-2222-2222-222222222222 (Centro → Norte)
-- Route 2 UUID: 33333333-3333-3333-3333-333333333333 (Sur → Este)
INSERT INTO app.routes (id, name, origin_name, origin_lat, origin_lon, destination_name, destination_lat, destination_lon, base_price_cents, currency, is_active) VALUES
('22222222-2222-2222-2222-222222222222', 'Ruta Centro - Norte', 'Plaza de Armas', -12.046374, -77.042793, 'Terminal Norte', -12.001234, -77.051234, 1550, 'PEN', true),
('33333333-3333-3333-3333-333333333333', 'Ruta Sur - Este', 'Estación Sur', -12.091234, -77.032123, 'Comercial Este', -12.071234, -77.001234, 1200, 'PEN', true);

-- Insert stops for Route 1 (Centro → Norte)
-- Stop 1: 44444444-4444-4444-4444-444444444444
-- Stop 2: 55555555-5555-5555-5555-555555555555
-- Stop 3: 66666666-6666-6666-6666-666666666666
INSERT INTO app.route_stops (id, route_id, stop_order, name, lat, lon, is_active) VALUES
('44444444-4444-4444-4444-444444444444', '22222222-2222-2222-2222-222222222222', 1, 'Plaza de Armas', -12.046374, -77.042793, true),
('55555555-5555-5555-5555-555555555555', '22222222-2222-2222-2222-222222222222', 2, 'Plaza Mayor', -12.031234, -77.045123, true),
('66666666-6666-6666-6666-666666666666', '22222222-2222-2222-2222-222222222222', 3, 'Terminal Norte', -12.001234, -77.051234, true);

-- Insert stops for Route 2 (Sur → Este)
-- Stop 1: 77777777-7777-7777-7777-777777777777
-- Stop 2: 88888888-8888-8888-8888-888888888888
-- Stop 3: 99999999-9999-9999-9999-999999999999
INSERT INTO app.route_stops (id, route_id, stop_order, name, lat, lon, is_active) VALUES
('77777777-7777-7777-7777-777777777777', '33333333-3333-3333-3333-333333333333', 1, 'Estación Sur', -12.091234, -77.032123, true),
('88888888-8888-8888-8888-888888888888', '33333333-3333-3333-3333-333333333333', 2, 'Puente Central', -12.081234, -77.021234, true),
('99999999-9999-9999-9999-999999999999', '33333333-3333-3333-3333-333333333333', 3, 'Comercial Este', -12.071234, -77.001234, true);

-- Summary of UUIDs for testing:
-- User (Passenger): 11111111-1111-1111-1111-111111111111
-- Route 1 (Centro-Norte): 22222222-2222-2222-2222-222222222222
-- Route 2 (Sur-Este): 33333333-3333-3333-3333-333333333333
-- Route 1 Stops: 44444444..., 55555555..., 66666666...
-- Route 2 Stops: 77777777..., 88888888..., 99999999...

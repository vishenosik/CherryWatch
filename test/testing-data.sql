INSERT INTO endpoints (id, service_name, url, interval) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'User Service', 'https://api.example.com/users/health', 30),
('6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'Payment Service', 'https://payments.example.com/status', 60),
('6ba7b811-9dad-11d1-80b4-00c04fd430c9', 'Inventory Service', 'https://inventory.example.com/healthcheck', 45),
('6ba7b812-9dad-11d1-80b4-00c04fd430c0', 'Email Service', 'https://mail.example.com/ping', 120),
('6ba7b813-9dad-11d1-80b4-00c04fd430c1', 'Auth Service', 'https://auth.example.com/health', 90),
('6ba7b814-9dad-11d1-80b4-00c04fd430c2', 'Analytics Service', 'https://analytics.example.com/status', 180),
('6ba7b815-9dad-11d1-80b4-00c04fd430c3', 'Notification Service', 'https://notify.example.com/health', 60),
('6ba7b816-9dad-11d1-80b4-00c04fd430c4', 'Content Service', 'https://content.example.com/ping', 30),
('6ba7b817-9dad-11d1-80b4-00c04fd430c5', 'Search Service', 'https://search.example.com/healthcheck', 45),
('6ba7b818-9dad-11d1-80b4-00c04fd430c6', 'Gateway Service', 'https://gateway.example.com/status', 120);

-- Insert success codes for each endpoint
INSERT INTO endpoint_success_codes (endpoint_id, code) VALUES
-- User Service
('550e8400-e29b-41d4-a716-446655440000', 200),
('550e8400-e29b-41d4-a716-446655440000', 201),
('550e8400-e29b-41d4-a716-446655440000', 204),
-- Payment Service
('6ba7b810-9dad-11d1-80b4-00c04fd430c8', 200),
('6ba7b810-9dad-11d1-80b4-00c04fd430c8', 202),
-- Inventory Service
('6ba7b811-9dad-11d1-80b4-00c04fd430c9', 200),
('6ba7b811-9dad-11d1-80b4-00c04fd430c9', 304),
-- Email Service
('6ba7b812-9dad-11d1-80b4-00c04fd430c0', 200),
-- Auth Service
('6ba7b813-9dad-11d1-80b4-00c04fd430c1', 200),
('6ba7b813-9dad-11d1-80b4-00c04fd430c1', 201),
('6ba7b813-9dad-11d1-80b4-00c04fd430c1', 401), -- Special case for auth
-- Analytics Service
('6ba7b814-9dad-11d1-80b4-00c04fd430c2', 200),
-- Notification Service
('6ba7b815-9dad-11d1-80b4-00c04fd430c3', 200),
('6ba7b815-9dad-11d1-80b4-00c04fd430c3', 202),
-- Content Service
('6ba7b816-9dad-11d1-80b4-00c04fd430c4', 200),
('6ba7b816-9dad-11d1-80b4-00c04fd430c4', 204),
-- Search Service
('6ba7b817-9dad-11d1-80b4-00c04fd430c5', 200),
('6ba7b817-9dad-11d1-80b4-00c04fd430c5', 206),
-- Gateway Service
('6ba7b818-9dad-11d1-80b4-00c04fd430c6', 200),
('6ba7b818-9dad-11d1-80b4-00c04fd430c6', 301),
('6ba7b818-9dad-11d1-80b4-00c04fd430c6', 302);

-- Insert notification services for each endpoint
INSERT INTO endpoint_notification_services (endpoint_id, service_name) VALUES
-- User Service
('550e8400-e29b-41d4-a716-446655440000', 'slack'),
('550e8400-e29b-41d4-a716-446655440000', 'email'),
('550e8400-e29b-41d4-a716-446655440000', 'sms'),
-- Payment Service
('6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'slack'),
('6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'pagerduty'),
-- Inventory Service
('6ba7b811-9dad-11d1-80b4-00c04fd430c9', 'email'),
-- Email Service
('6ba7b812-9dad-11d1-80b4-00c04fd430c0', 'slack'),
('6ba7b812-9dad-11d1-80b4-00c04fd430c0', 'email'),
-- Auth Service
('6ba7b813-9dad-11d1-80b4-00c04fd430c1', 'pagerduty'),
-- Analytics Service
('6ba7b814-9dad-11d1-80b4-00c04fd430c2', 'slack'),
-- Notification Service
('6ba7b815-9dad-11d1-80b4-00c04fd430c3', 'email'),
-- Content Service
('6ba7b816-9dad-11d1-80b4-00c04fd430c4', 'slack'),
('6ba7b816-9dad-11d1-80b4-00c04fd430c4', 'sms'),
-- Search Service
('6ba7b817-9dad-11d1-80b4-00c04fd430c5', 'slack'),
-- Gateway Service
('6ba7b818-9dad-11d1-80b4-00c04fd430c6', 'slack'),
('6ba7b818-9dad-11d1-80b4-00c04fd430c6', 'email'),
('6ba7b818-9dad-11d1-80b4-00c04fd430c6', 'pagerduty');
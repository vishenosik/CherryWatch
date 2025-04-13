-- +goose Up

-- Main endpoint table
CREATE TABLE endpoints (
    id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(64) NOT NULL,
    url TEXT NOT NULL,
    interval INTEGER NOT NULL,
    CONSTRAINT uk_endpoints_service_name UNIQUE (service_name),
    CONSTRAINT uk_endpoints_url UNIQUE (url)
);

-- Table for success codes (many-to-one relationship with endpoints)
CREATE TABLE endpoint_success_codes (
    endpoint_id VARCHAR(36) NOT NULL,
    code INTEGER NOT NULL,
    PRIMARY KEY (endpoint_id, code),
    FOREIGN KEY (endpoint_id) REFERENCES endpoints(id) ON DELETE CASCADE
);

-- Table for notification services (many-to-one relationship with endpoints)
CREATE TABLE endpoint_notification_services (
    endpoint_id VARCHAR(36) NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    PRIMARY KEY (endpoint_id, service_name),
    FOREIGN KEY (endpoint_id) REFERENCES endpoints(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE IF EXISTS endpoints;
DROP TABLE IF EXISTS endpoint_success_codes;
DROP TABLE IF EXISTS endpoint_notification_services;
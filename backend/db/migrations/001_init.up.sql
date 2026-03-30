BEGIN;

-- Required for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number TEXT NOT NULL UNIQUE,
    username TEXT UNIQUE,
    password TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

CREATE TABLE IF NOT EXISTS packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    package_type VARCHAR(20) NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    duration_minutes BIGINT,
    data_limit_mb BIGINT,
    speed_limit_up BIGINT,
    speed_limit_down BIGINT,
    validity_days BIGINT NOT NULL DEFAULT 30,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_packages_deleted_at ON packages (deleted_at);

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    package_id UUID NOT NULL REFERENCES packages(id),
    amount NUMERIC(12,2) NOT NULL,
    phone_number TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    mpesa_checkout_id TEXT UNIQUE,
    mpesa_receipt_number TEXT,
    mpesa_transaction_id TEXT,
    transaction_date TIMESTAMPTZ,
    result_code BIGINT,
    result_desc TEXT,
    callback_received BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_payments_deleted_at ON payments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments (user_id);
CREATE INDEX IF NOT EXISTS idx_payments_package_id ON payments (package_id);

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    package_id UUID NOT NULL REFERENCES packages(id),
    payment_id UUID REFERENCES payments(id),
    session_id TEXT UNIQUE,
    username TEXT NOT NULL,
    nas_ip_address TEXT,
    nas_port_id TEXT,
    ip_address TEXT,
    mac_address TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL,
    data_used_bytes BIGINT NOT NULL DEFAULT 0,
    data_limit_bytes BIGINT,
    time_used_minutes BIGINT NOT NULL DEFAULT 0,
    time_limit_minutes BIGINT,
    last_update_time TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sessions_deleted_at ON sessions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_package_id ON sessions (package_id);
CREATE INDEX IF NOT EXISTS idx_sessions_payment_id ON sessions (payment_id);
CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions (status);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions (expires_at);

COMMIT;

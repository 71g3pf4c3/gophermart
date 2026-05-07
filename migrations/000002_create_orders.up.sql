CREATE TABLE IF NOT EXISTS orders (
    number VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'NEW',
    accrual NUMERIC(15,2),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_uploaded ON orders(user_id, uploaded_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_pending ON orders(status) WHERE status IN ('NEW', 'PROCESSING');

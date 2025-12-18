CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE operators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL, 
    full_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_operators_username ON operators(username);

CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    address TEXT,
    postcode TEXT,
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    blocked_reason TEXT
);

CREATE INDEX idx_customers_name ON customers(full_name);
CREATE INDEX idx_customers_postcode ON customers(postcode);

CREATE TABLE customer_ledger_accounts (
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    ledger_id INT NOT NULL, -- ISO 4217
    tb_account_id UUID NOT NULL UNIQUE, 
    PRIMARY KEY (customer_id, ledger_id)
);

CREATE INDEX idx_ledger_accounts_tb_id ON customer_ledger_accounts(tb_account_id);

CREATE TABLE fx_trades (
    tb_pending_id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id),
    operator_id UUID NOT NULL REFERENCES operators(id),

    exchange_rate NUMERIC(18, 9) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes TEXT,
    is_unusual BOOLEAN NOT NULL DEFAULT FALSE,
    unusual_reason TEXT
);

CREATE INDEX idx_fx_trades_customer_id ON fx_trades(customer_id);
CREATE INDEX idx_fx_trades_operator_id ON fx_trades(operator_id);
CREATE INDEX idx_fx_trades_created_at ON fx_trades(created_at);
CREATE INDEX idx_fx_trades_unusual ON fx_trades(is_unusual) WHERE is_unusual IS TRUE;
CREATE INDEX idx_fx_trades_compliance_scan ON fx_trades(customer_id, created_at);

CREATE OR REPLACE FUNCTION fn_block_unusual_customers()
RETURNS TRIGGER AS $$
BEGIN
    IF (NEW.is_unusual = TRUE) THEN
        UPDATE customers
        SET is_blocked = TRUE,
            blocked_reason = 'Automatically blocked; an operator marked transaction ' || NEW.tb_pending_id || ' as unusual'
        WHERE id = NEW.customer_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_block_unusual_customer
AFTER INSERT OR UPDATE OF is_unusual ON fx_trades
FOR EACH ROW
EXECUTE FUNCTION fn_block_unusual_customers();

CREATE TABLE eod_closings (
    tb_pending_id UUID PRIMARY KEY, 
    ledger_id INT NOT NULL,
    operator_id UUID NOT NULL REFERENCES operators(id),
    closed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_eod_closings_date ON eod_closings(closed_at DESC);
CREATE INDEX idx_eod_closings_ledger ON eod_closings(ledger_id);
CREATE INDEX idx_eod_closings_operator ON eod_closings(operator_id, closed_at);

CREATE TABLE sod_openings (
    tb_pending_id UUID PRIMARY KEY, 
    ledger_id INT NOT NULL,
    operator_id UUID NOT NULL REFERENCES operators(id),
    opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sod_openings_ledger ON sod_openings(ledger_id);
CREATE INDEX idx_sod_openings_date ON sod_openings(opened_at DESC);

CREATE TABLE ledger_adjustments (
    tb_transfer_id UUID PRIMARY KEY, 
    ledger_id INT NOT NULL,
    operator_id UUID NOT NULL REFERENCES operators(id),
    adjustment_type TEXT NOT NULL CHECK (adjustment_type IN ('OVER', 'SHORT')),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ledger_adjustments_ledger ON ledger_adjustments(ledger_id);
CREATE INDEX idx_ledger_adjustments_operator ON ledger_adjustments(operator_id);
CREATE INDEX idx_ledger_adjustments_date ON ledger_adjustments(created_at DESC);

CREATE TABLE midday_counts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ledger_id INT NOT NULL,
    operator_id UUID NOT NULL REFERENCES operators(id),
    ledger_balance_at_count BIGINT NOT NULL, 
    physical_count_recorded BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_midday_counts_ledger_date ON midday_counts(ledger_id, created_at DESC);
CREATE INDEX idx_midday_counts_operator ON midday_counts(operator_id);

CREATE TABLE IF NOT EXISTS mfa_settings (
    id SERIAL PRIMARY KEY,
    enforce_2fa BOOLEAN NOT NULL DEFAULT FALSE,
    updated_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Insert default row
INSERT INTO mfa_settings (enforce_2fa) VALUES (FALSE) ON CONFLICT DO NOTHING;

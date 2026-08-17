-- Create shared expense tables

DO $$ BEGIN
    CREATE TYPE split_bill_status AS ENUM ('OPEN', 'PARTIALLY_SETTLED', 'SETTLED');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS split_bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id BIGINT NOT NULL REFERENCES users(id),
    payer_id BIGINT NOT NULL REFERENCES users(id),
    category_id BIGINT NOT NULL REFERENCES categories(id),
    title VARCHAR(255) NOT NULL,
    total_amount DECIMAL(15, 2) NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    status split_bill_status DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS split_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    split_bill_id UUID NOT NULL REFERENCES split_bills(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id), -- Nullable for Shadow Users
    shadow_name VARCHAR(255),            -- Used if user_id is null
    share_amount DECIMAL(15, 2) NOT NULL,
    is_paid BOOLEAN DEFAULT FALSE,
    settled_at TIMESTAMP WITH TIME ZONE,
    transaction_id BIGINT REFERENCES transactions(id), -- Link to the auto-generated income/expense
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Performance Indexes
CREATE INDEX IF NOT EXISTS idx_split_bills_payer_id ON split_bills(payer_id);
CREATE INDEX IF NOT EXISTS idx_split_participants_user_id ON split_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_split_participants_bill_id ON split_participants(split_bill_id);

-- Add customization (split method), tax, discount, and per-item support to split bills

ALTER TABLE split_bills
    ADD COLUMN IF NOT EXISTS subtotal DECIMAL(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS tax_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS split_method VARCHAR(20) NOT NULL DEFAULT 'EQUAL';

-- Backfill subtotal for existing rows (previously the total was the base)
UPDATE split_bills SET subtotal = total_amount WHERE subtotal = 0;

ALTER TABLE split_participants
    ADD COLUMN IF NOT EXISTS base_share DECIMAL(15, 2) NOT NULL DEFAULT 0;

-- Backfill base_share from the existing final share for old rows
UPDATE split_participants SET base_share = share_amount WHERE base_share = 0;

-- Line items (used by the ITEM split method)
CREATE TABLE IF NOT EXISTS split_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    split_bill_id UUID NOT NULL REFERENCES split_bills(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(15, 2) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Which participants share which item
CREATE TABLE IF NOT EXISTS split_item_participants (
    item_id UUID NOT NULL REFERENCES split_items(id) ON DELETE CASCADE,
    participant_id UUID NOT NULL REFERENCES split_participants(id) ON DELETE CASCADE,
    PRIMARY KEY (item_id, participant_id)
);

CREATE INDEX IF NOT EXISTS idx_split_items_bill_id ON split_items(split_bill_id);
CREATE INDEX IF NOT EXISTS idx_split_item_participants_item_id ON split_item_participants(item_id);

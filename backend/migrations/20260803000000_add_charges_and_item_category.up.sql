-- Add service charge / other charge to split bills, and per-item category.

ALTER TABLE split_bills
    ADD COLUMN IF NOT EXISTS service_charge DECIMAL(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS other_charge DECIMAL(15, 2) NOT NULL DEFAULT 0;

-- Per-item category (nullable; falls back to the bill's category). SET NULL on
-- category delete so a category can still be removed.
ALTER TABLE split_items
    ADD COLUMN IF NOT EXISTS category_id BIGINT REFERENCES categories(id) ON DELETE SET NULL;

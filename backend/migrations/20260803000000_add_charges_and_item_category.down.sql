-- Revert charges and per-item category.

ALTER TABLE split_items
    DROP COLUMN IF EXISTS category_id;

ALTER TABLE split_bills
    DROP COLUMN IF EXISTS service_charge,
    DROP COLUMN IF EXISTS other_charge;

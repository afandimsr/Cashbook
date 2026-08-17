-- Revert split customization additions

DROP INDEX IF EXISTS idx_split_item_participants_item_id;
DROP INDEX IF EXISTS idx_split_items_bill_id;

DROP TABLE IF EXISTS split_item_participants;
DROP TABLE IF EXISTS split_items;

ALTER TABLE split_participants
    DROP COLUMN IF EXISTS base_share;

ALTER TABLE split_bills
    DROP COLUMN IF EXISTS subtotal,
    DROP COLUMN IF EXISTS tax_amount,
    DROP COLUMN IF EXISTS discount_amount,
    DROP COLUMN IF EXISTS split_method;

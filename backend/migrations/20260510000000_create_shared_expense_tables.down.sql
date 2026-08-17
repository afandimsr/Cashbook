-- Drop shared expense tables

DROP INDEX IF EXISTS idx_split_participants_bill_id;
DROP INDEX IF EXISTS idx_split_participants_user_id;
DROP INDEX IF EXISTS idx_split_bills_payer_id;

DROP TABLE IF EXISTS split_participants;
DROP TABLE IF EXISTS split_bills;

DROP TYPE IF EXISTS split_bill_status;

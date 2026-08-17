-- Allow deleting a transaction that a split participant references.
-- The original FK had no ON DELETE action (RESTRICT), which blocked deleting
-- any transaction linked from a settled split participant. Switch to SET NULL.

ALTER TABLE split_participants
    DROP CONSTRAINT IF EXISTS split_participants_transaction_id_fkey;

ALTER TABLE split_participants
    ADD CONSTRAINT split_participants_transaction_id_fkey
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE SET NULL;

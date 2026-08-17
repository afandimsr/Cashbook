-- Revert to the original FK without an ON DELETE action.

ALTER TABLE split_participants
    DROP CONSTRAINT IF EXISTS split_participants_transaction_id_fkey;

ALTER TABLE split_participants
    ADD CONSTRAINT split_participants_transaction_id_fkey
    FOREIGN KEY (transaction_id) REFERENCES transactions(id);

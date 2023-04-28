CREATE TABLE IF NOT EXISTS account_event (
  account_id uuid,
  event_id timeuuid,
  payload blob,
  PRIMARY KEY (account_id, event_id)
) WITH COMPACT STORAGE;

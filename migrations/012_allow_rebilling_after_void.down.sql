DROP INDEX IF EXISTS bills_order_id_active;
ALTER TABLE bills ADD CONSTRAINT bills_order_id_key UNIQUE (order_id);

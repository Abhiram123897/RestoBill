ALTER TABLE bills DROP CONSTRAINT bills_order_id_key;
CREATE UNIQUE INDEX bills_order_id_active ON bills (order_id) WHERE status != 'void';

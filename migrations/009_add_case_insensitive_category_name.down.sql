DROP INDEX IF EXISTS categories_name_unique;
ALTER TABLE categories ADD CONSTRAINT categories_name_key UNIQUE (name);

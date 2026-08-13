CREATE UNIQUE INDEX menu_items_name_category_unique ON menu_items (LOWER(name), category_id);

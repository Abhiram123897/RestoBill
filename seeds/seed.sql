-- Seed data for RestoBill
-- Run: psql -U postgres -d restobill -f seeds/seed.sql

-- Categories
INSERT INTO categories (name, display_order) VALUES ('Starters', 1);
INSERT INTO categories (name, display_order) VALUES ('Main Course', 2);
INSERT INTO categories (name, display_order) VALUES ('Beverages', 3);
INSERT INTO categories (name, display_order) VALUES ('Desserts', 4);

-- Menu Items (Starters)
INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Paneer Tikka', 250.00, true FROM categories WHERE name = 'Starters';

INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Chicken 65', 300.00, true FROM categories WHERE name = 'Starters';

INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Spring Rolls', 180.00, true FROM categories WHERE name = 'Starters';

-- Menu Items (Main Course)
INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Butter Chicken', 350.00, true FROM categories WHERE name = 'Main Course';

INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Dal Makhani', 220.00, true FROM categories WHERE name = 'Main Course';

INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Veg Biryani', 280.00, true FROM categories WHERE name = 'Main Course';

INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Chicken Biryani', 320.00, true FROM categories WHERE name = 'Main Course';

-- Menu Items (Beverages)
INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Coke', 60.00, true FROM categories WHERE name = 'Beverages';

INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Fresh Lime Soda', 80.00, true FROM categories WHERE name = 'Beverages';

INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Mango Lassi', 100.00, true FROM categories WHERE name = 'Beverages';

-- Menu Items (Desserts)
INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Gulab Jamun', 120.00, true FROM categories WHERE name = 'Desserts';

INSERT INTO menu_items (category_id, name, price, is_available)
SELECT id, 'Ice Cream', 150.00, true FROM categories WHERE name = 'Desserts';

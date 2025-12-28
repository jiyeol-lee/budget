-- Migration: 2025-12-07-001
-- Description: Insert expected expenses for weekly and monthly grocery purchases

-- ============================================================================
-- Weekly Expected Expenses
-- ============================================================================
INSERT INTO expected_expenses (item_name, source, expected_amount, expense_type) VALUES
    ('Milk', 'Publix', 4.00, 'weekly'),
    ('Creamer', 'Costco', 8.00, 'weekly'),
    ('Eggs', 'Costco', 15.00, 'weekly'),
    ('Onions', 'Costco', 10.00, 'weekly'),
    ('Potatoes', 'Costco', 10.00, 'weekly'),
    ('Green Onions', 'Costco', 7.00, 'weekly'),
    ('Zuccini', 'H-Mart', 5.00, 'weekly'),
    ('Mushrooms', 'H-Mart', 20.00, 'weekly'),
    ('Chicken', 'Costco', 25.00, 'weekly'),
    ('Beef', 'H-Mart', 30.00, 'weekly'),
    ('Yogurt (Ian)', 'H-Mart', 5.00, 'weekly'),
    ('Squid banchan', 'H-Mart', 10.00, 'weekly'),
    ('Formula', 'Walmart', 40.00, 'weekly'),
    ('Chives', 'H-Mart', 9.00, 'weekly'),
    ('Baby Cereal', 'Walmart', 4.99, 'weekly'),
    ('Hoppang', 'H-Mart', 7.98, 'weekly'),
    ('Rice Cakes', 'H-Mart', 4.99, 'weekly');

-- ============================================================================
-- Monthly Expected Expenses
-- ============================================================================
INSERT INTO expected_expenses (item_name, source, expected_amount, expense_type) VALUES
    ('Coffee', 'Costco', 44.99, 'monthly'),
    ('Protein Shake', 'Costco', 36.99, 'monthly'),
    ('Diapers', 'Costco', 49.99, 'monthly'),
    ('Baby Wipes', 'Costco', 22.99, 'monthly'),
    ('Toilet Paper', 'Costco', 24.99, 'monthly'),
    ('Paper Towels', 'Costco', 23.99, 'monthly'),
    ('Tissues', 'Costco', 19.99, 'monthly'),
    ('Applesauce', 'Costco', 13.99, 'monthly'),
    ('Yogurt Melts', 'Costco', 18.99, 'monthly'),
    ('Shin Ramen', 'Costco', 19.99, 'monthly'),
    ('Vegetable Oil', 'Costco', 15.49, 'monthly'),
    ('Spagetti Noodles', 'Costco', 13.99, 'monthly'),
    ('Marinara', 'Costco', 14.49, 'monthly'),
    ('Alfredo', 'Costco', 10.99, 'monthly'),
    ('Dish Soap', 'Costco', 8.00, 'monthly'),
    ('Rice', 'H-Mart', 16.99, 'monthly'),
    ('Seaweed', 'H-Mart', 20.00, 'monthly'),
    ('Tofu', 'H-Mart', 10.00, 'monthly'),
    ('Kimchi', 'H-Mart', 30.00, 'monthly'),
    ('Panko', 'H-Mart', 5.00, 'monthly'),
    ('K-Pancake Mix', 'H-Mart', 5.00, 'monthly'),
    ('K-Frying Mix', 'H-Mart', 5.00, 'monthly'),
    ('Rice cakes', 'H-Mart', 7.00, 'monthly'),
    ('Soy Sauce', 'H-Mart', 5.00, 'monthly'),
    ('Gochujang', 'H-Mart', 10.00, 'monthly'),
    ('Soybean Paste', 'H-Mart', 6.00, 'monthly'),
    ('Grapeseed Oil', 'H-Mart', 16.00, 'monthly'),
    ('Dried Seaweed', 'H-Mart', 6.99, 'monthly'),
    ('Bottle Soap', 'Target', 10.00, 'monthly');

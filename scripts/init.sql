CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    price DECIMAL
);

-- Заполним тестовыми данными (например, 1000 записей)
INSERT INTO products (name, description, price)
SELECT 
    'Product ' || i, 
    'Description for product ' || i, 
    random() * 1000
FROM generate_series(1, 1000) s(i);

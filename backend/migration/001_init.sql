CREATE TABLE users (
    id           SERIAL PRIMARY KEY,
    username     VARCHAR(100) UNIQUE NOT NULL,
    password     TEXT NOT NULL,
    phone_number VARCHAR(20),
    street       TEXT,
    city         VARCHAR(100),
    state        VARCHAR(100),
    country      VARCHAR(100),
    zip_code     VARCHAR(20),
    role         VARCHAR(20) NOT NULL DEFAULT 'customer'
);

CREATE TABLE phones (
    id          SERIAL PRIMARY KEY,
    brand       VARCHAR(100) NOT NULL,
    model       VARCHAR(100) NOT NULL,
    price       NUMERIC(10,2) NOT NULL,
    stock       INT NOT NULL DEFAULT 0,
    description TEXT,
    image_url   TEXT
);

CREATE TABLE carts (
    id      SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    status  VARCHAR(20) NOT NULL DEFAULT 'active'
);

CREATE TABLE cart_items (
    id       SERIAL PRIMARY KEY,
    cart_id  INT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    phone_id INT NOT NULL REFERENCES phones(id),
    quantity INT NOT NULL,
    price    NUMERIC(10,2) NOT NULL,
    UNIQUE (cart_id, phone_id)
);

CREATE TABLE orders (
    id         TEXT PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users(id),
    status     VARCHAR(20) NOT NULL DEFAULT 'succeeded',
    total      NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id         SERIAL PRIMARY KEY,
    order_id   TEXT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    phone_id   INT NOT NULL REFERENCES phones(id),
    phone_name TEXT NOT NULL DEFAULT '',
    quantity   INT NOT NULL,
    price      NUMERIC(10,2) NOT NULL
);

CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,  -- e.g., 'phone'
    entity_id INT NOT NULL,            -- ID of the phone being changed
    action VARCHAR(50) NOT NULL,       -- 'price_update', 'stock_addition', 'created'
    old_value JSONB,                   -- The value before the change
    new_value JSONB,                   -- The value after the change
    changed_by INT,                    -- UserID of the admin who did it
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast searching by phone
CREATE INDEX idx_audit_entity ON audit_logs (entity_type, entity_id);

-- ================================================
-- Seed phones (prices in SGD)
-- ================================================

INSERT INTO phones (brand, model, price, stock, description, image_url) VALUES

-- Apple
('Apple', 'iPhone Air',           1499.00, 30, 'Sleek and lightweight iPhone with A18 chip and stunning OLED display.',                       '/src/assets/images/Apple-iPhone-Air.jpg'),
('Apple', 'iPhone 17',            1299.00, 40, 'The standard iPhone 17 with A19 chip, improved camera system and all-day battery life.',      '/src/assets/images/iphone-17.jpg'),
('Apple', 'iPhone 17 Pro',        1649.00, 25, 'Pro-grade titanium iPhone with A19 Pro chip, ProMotion display and advanced camera controls.','/src/assets/images/Apple-iPhone17-Pro.jpg'),
('Apple', 'iPhone 17 Pro Max',    1949.00, 20, 'The ultimate iPhone with the largest Pro display, longest battery life and A19 Pro chip.',    '/src/assets/images/Apple-iPhone17-Pro-Max.jpg'),
('Apple', 'iPhone 17e',           1099.00, 50, 'Affordable iPhone 17e with A18 chip, compact design and essential Pro features.',             '/src/assets/images/Apple-iPhone17e.jpg'),

-- Samsung
('Samsung', 'Galaxy S25',         1199.00, 35, 'Flagship Samsung with Snapdragon 8 Elite, refined design and AI-powered camera.',             '/src/assets/images/samsung-s25.jpg'),
('Samsung', 'Galaxy S25 Plus',    1399.00, 30, 'Larger Galaxy S25 with bigger battery, enhanced display and Snapdragon 8 Elite.',             '/src/assets/images/samsung-s25-plus.jpg'),
('Samsung', 'Galaxy S25 Ultra',   1799.00, 20, 'The most powerful Galaxy with built-in S Pen, 200MP camera and titanium frame.',              '/src/assets/images/samsung-s25-ultra.jpg'),
('Samsung', 'Galaxy S25 FE',      899.00,  45, 'Fan Edition Galaxy S25 with flagship features at a more accessible price.',                  '/src/assets/images/samsung-s25-fe.jpg'),
('Samsung', 'Galaxy S26',         1249.00, 30, 'Next-gen Galaxy S26 with improved AI features, brighter display and faster charging.',        '/src/assets/images/samsung-s26.jpg'),
('Samsung', 'Galaxy S26 Plus',    1499.00, 25, 'Galaxy S26 Plus with larger screen, extended battery and premium build quality.',             '/src/assets/images/samsung-s26-plus.jpg'),
('Samsung', 'Galaxy S26 Ultra',   1849.00, 15, 'Top-of-the-line Galaxy with advanced S Pen, pro camera system and titanium chassis.',         '/src/assets/images/samsung-s26-ultra.jpg'),
('Samsung', 'Galaxy Z Flip 7',    1499.00, 20, 'Compact foldable phone with larger cover screen, Snapdragon 8 Elite and refined hinge.',     '/src/assets/images/samsung-galaxy-z-flip-7.jpg'),

-- Google
('Google', 'Pixel 10',            1099.00, 35, 'Pure Android experience with Google Tensor G5, exceptional computational photography.',       '/src/assets/images/Google-Pixel-10.jpg'),
('Google', 'Pixel 10 Pro',        1349.00, 25, 'Pro Pixel with upgraded telephoto camera, refined design and advanced AI features.',          '/src/assets/images/Google-Pixel-10-Pro.jpg'),
('Google', 'Pixel 10 Pro XL',     1499.00, 20, 'Largest Pixel with expansive display, massive battery and full Pro camera system.',           '/src/assets/images/Google-Pixel-10-Pro-ProXL.jpg'),
('Google', 'Pixel 10 Pro Fold',   2499.00, 10, 'Google foldable with inner and outer displays, Tensor G5 and premium build.',                '/src/assets/images/Google-Pixel-10-Pro-Fold.jpg');
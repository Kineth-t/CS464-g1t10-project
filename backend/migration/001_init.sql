CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(100) UNIQUE NOT NULL,
    password    TEXT NOT NULL,
    phone_number VARCHAR(20),
    street      TEXT,
    city        VARCHAR(100),
    state       VARCHAR(100),
    country     VARCHAR(100),
    zip_code    VARCHAR(20),
    role        VARCHAR(20) NOT NULL DEFAULT 'customer'
);

CREATE TABLE phones (
    id          SERIAL PRIMARY KEY,
    brand       VARCHAR(100) NOT NULL,
    model       VARCHAR(100) NOT NULL,
    price       NUMERIC(10,2) NOT NULL,
    stock       INT NOT NULL DEFAULT 0,
    description TEXT
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
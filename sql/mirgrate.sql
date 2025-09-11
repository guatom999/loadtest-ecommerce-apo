CREATE EXTENSION
IF NOT EXISTS "uuid-ossp";


CREATE TABLE
IF NOT EXISTS users
(
id UUID PRIMARY KEY DEFAULT uuid_generate_v4
(),
email TEXT UNIQUE NOT NULL,
password_hash TEXT NOT NULL,
first_name TEXT NOT NULL,
last_name TEXT NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT now
()
);


CREATE TABLE
IF NOT EXISTS products
(
id UUID PRIMARY KEY DEFAULT uuid_generate_v4
(),
name TEXT NOT NULL,
description TEXT,
price NUMERIC
(12,2) NOT NULL CHECK
(price >= 0),
expires_at DATE,
created_at TIMESTAMPTZ NOT NULL DEFAULT now
()
);


-- helpful index for listing
CREATE INDEX
IF NOT EXISTS idx_products_created_at ON products
(created_at DESC);
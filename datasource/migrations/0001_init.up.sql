-- Enable UUID extension (optional, but good practice)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Persons table (no dependencies)
CREATE TABLE persons (
    id            BIGSERIAL PRIMARY KEY,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    email         VARCHAR(100) UNIQUE NOT NULL,
    person_number VARCHAR(100) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Customers table (depends on persons)
CREATE TABLE customers (
    id           BIGSERIAL PRIMARY KEY,
    username     VARCHAR(100) NOT NULL,
    email        VARCHAR(100) UNIQUE NOT NULL,
    person_id    BIGINT UNIQUE REFERENCES persons(id),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Accounts table (depends on customers)
CREATE TABLE accounts (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(50) UNIQUE NOT NULL,
    email           VARCHAR(100) UNIQUE NOT NULL,
    phone           VARCHAR(20) NULL,
    customer_id     BIGINT NULL REFERENCES customers(id),
    active          BOOLEAN DEFAULT TRUE,
    balance         DECIMAL(10, 2) DEFAULT 0.00,
    currency        VARCHAR(3) DEFAULT 'USD',
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Cards table (depends on accounts)
CREATE TABLE cards (
    id           BIGSERIAL PRIMARY KEY,
    account_id   BIGINT NOT NULL REFERENCES accounts(id),
    last_four    VARCHAR(4) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at   TIMESTAMP DEFAULT (CURRENT_TIMESTAMP + INTERVAL '2 year'),
    cvv          VARCHAR(3) NOT NULL,
    active       BOOLEAN DEFAULT TRUE
);

-- 5. Transactions table (depends on accounts and cards)
CREATE TABLE transactions (
    id           BIGSERIAL PRIMARY KEY,
    account_id   BIGINT NOT NULL REFERENCES accounts(id),
    card_id      BIGINT NULL REFERENCES cards(id),
    direction    VARCHAR(10) NOT NULL CHECK (direction IN ('in', 'out')),
    amount       DECIMAL(10, 2) NOT NULL,
    status       VARCHAR(20) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

-- 6. Addresses table (depends on customers)
CREATE TABLE addresses (
    id           BIGSERIAL PRIMARY KEY,
    street       VARCHAR(255) NOT NULL,
    city         VARCHAR(100) NOT NULL,
    state        VARCHAR(100) NOT NULL,
    zip_code     VARCHAR(20) NOT NULL,
    country      VARCHAR(100) NOT NULL,
    customer_id  BIGINT NULL REFERENCES customers(id),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
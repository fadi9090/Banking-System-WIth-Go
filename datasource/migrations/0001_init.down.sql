-- Drop tables in reverse order
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS persons;

-- Drop extension if used
DROP EXTENSION IF EXISTS "uuid-ossp";
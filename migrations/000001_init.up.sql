CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE ROLE AS ENUM ('admin', 'support', 'regular');

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  last_name VARCHAR,
  first_name VARCHAR,
  middle_name VARCHAR,
  email VARCHAR UNIQUE NOT NULL,
  hashed_password VARCHAR NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS roles (
  user_id UUID,
  role ROLE NOT NULL DEFAULT 'regular'
);
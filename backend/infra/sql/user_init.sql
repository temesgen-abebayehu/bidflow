-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Companies Table
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    logo_url TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    founded_date DATE,
    area VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'BIDDER', -- SELLER, BIDDER, ADMIN
    company_id UUID REFERENCES companies(id) ON DELETE SET NULL,
    is_verified BOOLEAN DEFAULT FALSE, -- For seller verification
    is_active BOOLEAN DEFAULT FALSE,   -- For email verification
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    two_factor_secret TEXT,            -- TOTP Secret
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
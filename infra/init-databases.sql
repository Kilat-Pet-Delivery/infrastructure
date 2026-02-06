-- Kilat Pet Runner: Database Initialization Script
-- This script runs on first PostgreSQL startup via docker-entrypoint-initdb.d

-- Create separate databases for each bounded context
CREATE DATABASE kilat_booking;
CREATE DATABASE kilat_payment;
CREATE DATABASE kilat_runner;
CREATE DATABASE kilat_identity;
CREATE DATABASE kilat_tracking;

-- Enable extensions in each database

\c kilat_identity
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c kilat_runner
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c kilat_booking
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c kilat_payment
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c kilat_tracking
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

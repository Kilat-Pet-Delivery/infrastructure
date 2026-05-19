-- Kilat Pet Runner: Database Initialization Script
-- This script runs on first PostgreSQL startup via docker-entrypoint-initdb.d

-- Create separate databases for each bounded context
CREATE DATABASE kilat_booking;
CREATE DATABASE kilat_payment;
CREATE DATABASE kilat_runner;
CREATE DATABASE kilat_identity;
CREATE DATABASE kilat_tracking;
CREATE DATABASE kilat_notification;
CREATE DATABASE kilat_review;
CREATE DATABASE service_chat;
CREATE DATABASE service_incident;
CREATE DATABASE service_loyalty;
CREATE DATABASE service_zones;

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

\c kilat_notification
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c kilat_review
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c service_chat
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c service_incident
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c service_loyalty
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\c service_zones
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

import { Knex } from "knex";

export async function up(knex: Knex): Promise<void> {
  return knex.schema.raw(`
    BEGIN;

    SET statement_timeout = 0;
    SET lock_timeout = 0;
    SET idle_in_transaction_session_timeout = 0;
    SET client_encoding = 'UTF8';
    SET standard_conforming_strings = on;
    SET check_function_bodies = false;
    SET xmloption = content;
    SET client_min_messages = warning;
    SET row_security = off;
    
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
    COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';
    
    CREATE EXTENSION IF NOT EXISTS "pg_trgm" WITH SCHEMA public;

    -- CREATE EXTENSION IF NOT EXISTS postgis;
    -- CREATE EXTENSION postgis_topology;
    
    SET default_tablespace = '';
    COMMIT;
  `);
}

export async function down(knex: Knex): Promise<void> {}

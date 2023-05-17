import { Knex } from "knex";

export async function up(knex: Knex): Promise<void> {
    return knex.schema.raw(`
        CREATE INDEX post_search_idx
        ON posts USING GIN (title gin_trgm_ops, "content" gin_trgm_ops);
    `);
  }
  
  export async function down(knex: Knex): Promise<void> {
    return knex.schema.raw(`
        DROP INDEX IF EXISTS post_search_idx;
    `);
  }
  
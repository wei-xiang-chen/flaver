import { Knex } from "knex";

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("posts", (tableBuilder) => {
    tableBuilder.bigIncrements("id", { primaryKey: true });
    tableBuilder.string("title").notNullable();
    tableBuilder.string("content").notNullable();
    tableBuilder.specificType("img_urls", "text[]").notNullable();
    tableBuilder.integer("rating").notNullable();
    tableBuilder.specificType("topic_ids", "integer[]").notNullable();
    tableBuilder
      .uuid("owner_uid")
      .notNullable()
      .references("uid")
      .inTable("users");
    tableBuilder
      .string("google_place_id")
      .notNullable()
      .references("place_id")
      .inTable("restaurants");
    tableBuilder.string("instantly_role").notNullable();
    tableBuilder
      .timestamp("created_at")
      .defaultTo(knex.raw("CURRENT_TIMESTAMP"));
    tableBuilder
      .timestamp("updated_at")
      .defaultTo(knex.raw("CURRENT_TIMESTAMP"));
    tableBuilder.timestamp("deleted_at");
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists("posts");
}

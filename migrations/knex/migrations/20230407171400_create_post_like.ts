import { Knex } from "knex";

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("post_likes", (tableBuilder) => {
    tableBuilder
      .uuid("user_uid")
      .notNullable()
      .references("uid")
      .inTable("users");
    tableBuilder
      .bigInteger("post_id")
      .notNullable()
      .references("id")
      .inTable("posts");
    tableBuilder
      .timestamp("created_at")
      .defaultTo(knex.raw("CURRENT_TIMESTAMP"));

    tableBuilder.primary(["user_uid", "post_id"]);
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists("post_likes");
}

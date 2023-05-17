import { Knex } from "knex";

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("user_topics", (tableBuilder) => {
    tableBuilder
      .uuid("user_uid")
      .notNullable()
      .references("uid")
      .inTable("users");
    tableBuilder
      .bigInteger("topic_id")
      .notNullable()
      .references("id")
      .inTable("topics");
    tableBuilder
      .timestamp("created_at")
      .defaultTo(knex.raw("CURRENT_TIMESTAMP"));

    tableBuilder.primary(["user_uid", "topic_id"]);
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists("user_topics");
}

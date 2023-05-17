import { Knex } from "knex";

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("users", (tableBuilder) => {
    tableBuilder
      .uuid("uid", { primaryKey: true })
      .defaultTo(knex.raw("uuid_generate_v4()"));
    tableBuilder.string("nickname").notNullable();
    tableBuilder.string("google_id_token").unique();
    tableBuilder.string("apple_id_token").unique();
    tableBuilder.string("avatar_img_url").notNullable();
    tableBuilder.timestamp("birthday");
    tableBuilder.string("email").notNullable().unique();
    tableBuilder.string("role").notNullable();
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
  return knex.schema.dropTableIfExists("users");
}

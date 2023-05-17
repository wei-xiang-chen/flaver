import { Knex } from "knex";

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("topic_groups", (tableBuilder) => {
    tableBuilder.bigIncrements("id", { primaryKey: true });
    tableBuilder.string("name").notNullable().unique();
    tableBuilder.string("img_url").notNullable();
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists("topic_groups");
}

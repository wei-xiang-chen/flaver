import { Knex } from "knex";

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("topics", (tableBuilder) => {
    tableBuilder.bigIncrements("id", { primaryKey: true });
    tableBuilder.string("name").notNullable().unique();
    tableBuilder.string("img_url").notNullable();
    tableBuilder
      .bigInteger("topic_group_id")
      .notNullable()
      .references("id")
      .inTable("topic_groups");
  });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists("topics");
}

import { Knex } from "knex";

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("restaurants", (tableBuilder) => {
    tableBuilder.string("place_id").primary();
    tableBuilder.string("name").notNullable();
    tableBuilder.string("phone");
    tableBuilder.string("address").nullable();
    // tableBuilder.specificType('location', 'GEOGRAPHY(POINT, 4326)').nullable();
    tableBuilder.float("lat").notNullable();
    tableBuilder.float("lng").notNullable();
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
  return knex.schema.dropTableIfExists("restaurants");
}

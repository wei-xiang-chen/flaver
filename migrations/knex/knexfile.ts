import type { Knex } from "knex";

const knexConfig: { [key: string]: Knex.Config } = {
  local: {
    client: "postgresql",
    connection: {
      host: "127.0.0.1",
      user: "flaver",
      password: "!QAZ2wsx",
      port: 5432,
      database: "flaver",
    },
    seeds: {
      directory: "./seeds/settings_local",
    },
  },
  server: {},
  dev: {},
};

module.exports = knexConfig;

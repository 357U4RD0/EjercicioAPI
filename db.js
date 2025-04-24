const { Pool } = require('pg');

const pool = new Pool({
  user: 'postgres',
  host: 'localhost',
  database: 'Incidentes',
  password: 'P0S7GR3SQ1',
  port: 5432,
});

module.exports = pool;
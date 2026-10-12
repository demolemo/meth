-- Applied on every Open, so a fresh database and an existing one converge on
-- the same shape. Every statement is IF NOT EXISTS: this is not a migration
-- runner and it will not alter a table that is already there, so a column
-- added below never reaches a database that predates it.

CREATE TABLE IF NOT EXISTS metrics (
  id   INTEGER PRIMARY KEY,
  name TEXT,
  unit TEXT NOT NULL DEFAULT ''
);

-- name is the main handle for a metric, ids are hard to keep track of.
-- NOCASE so `weight` and `Weight` are the same metric. Names are dotted
-- paths, `sleep.rem`, and this index is the only thing that makes a family
-- unambiguous - nothing in the schema knows that `sleep.` is a prefix.
CREATE UNIQUE INDEX IF NOT EXISTS metrics_name ON metrics(name COLLATE NOCASE);

CREATE TABLE IF NOT EXISTS metric_values (
  id         INTEGER PRIMARY KEY,
  metric_id  INTEGER NOT NULL REFERENCES metrics(id),
  value      FLOAT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- A reading is one number for one metric at one instant, and this is that
-- identity made enforceable. created_at is always written in UTC (see
-- sqlite/value.go): go-sqlite3 stores a time as text with its offset, so two
-- spellings of one instant would otherwise both fit through the index.
CREATE UNIQUE INDEX IF NOT EXISTS metric_values_reading
  ON metric_values(metric_id, created_at);

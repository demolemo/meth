-- Applied on every Open, so a fresh database and an existing one converge on
-- the same shape. Every statement is IF NOT EXISTS: this is not a migration
-- runner and it will not alter a table that is already there.

CREATE TABLE IF NOT EXISTS metrics (
  id   INTEGER PRIMARY KEY,
  name TEXT
);

-- name is the main handle for a metric, ids are hard to keep track of.
-- NOCASE so `weight` and `Weight` are the same metric.
CREATE UNIQUE INDEX IF NOT EXISTS metrics_name ON metrics(name COLLATE NOCASE);

CREATE TABLE IF NOT EXISTS metric_values (
  id         INTEGER PRIMARY KEY,
  metric_id  INTEGER NOT NULL REFERENCES metrics(id),
  value      FLOAT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

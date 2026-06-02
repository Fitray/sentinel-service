ALTER TABLE app.requests
DROP COLUMN scale;

DROP INDEX IF EXISTS app.requests_unique_query_idx;
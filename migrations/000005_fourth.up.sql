ALTER TABLE app.requests
ADD COLUMN scale INTEGER NOT NULL DEFAULT 10;

CREATE UNIQUE INDEX requests_unique_query_idx
ON app.requests (
    user_id,
    city,
    date_from,
    date_to,
    bands,
    scale,
    dimensions
);
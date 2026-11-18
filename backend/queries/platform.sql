-- name: DatabaseNow :one
-- Used by the readiness check: proves the pool can run a query, not only open a connection.
SELECT now()::timestamptz AS now;

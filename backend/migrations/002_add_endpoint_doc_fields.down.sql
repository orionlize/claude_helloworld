-- Rollback: Remove the added documentation fields from endpoints table

DROP INDEX IF EXISTS idx_endpoints_request_params;
DROP INDEX IF EXISTS idx_endpoints_response_params;

ALTER TABLE endpoints
DROP COLUMN IF EXISTS request_params,
DROP COLUMN IF EXISTS request_body,
DROP COLUMN IF EXISTS response_params,
DROP COLUMN IF EXISTS response_body;

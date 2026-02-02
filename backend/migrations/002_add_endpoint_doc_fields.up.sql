-- Add request_params, request_body, response_params, response_body fields to endpoints table
-- These fields store detailed API documentation information

ALTER TABLE endpoints
ADD COLUMN IF NOT EXISTS request_params JSONB DEFAULT '[]',
ADD COLUMN IF NOT EXISTS request_body JSONB,
ADD COLUMN IF NOT EXISTS response_params JSONB DEFAULT '[]',
ADD COLUMN IF NOT EXISTS response_body JSONB;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_endpoints_request_params ON endpoints USING gin(request_params);
CREATE INDEX IF NOT EXISTS idx_endpoints_response_params ON endpoints USING gin(response_params);

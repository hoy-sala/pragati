ALTER TABLE refresh_tokens ADD COLUMN IF NOT EXISTS lookup_hash VARCHAR(64) UNIQUE;

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_lookup_hash ON refresh_tokens(lookup_hash);

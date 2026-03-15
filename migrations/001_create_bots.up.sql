CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS bots (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(100) NOT NULL,
    role         VARCHAR(100) NOT NULL DEFAULT '',
    vibe         TEXT         NOT NULL DEFAULT '',
    emoji        VARCHAR(10)  NOT NULL DEFAULT '',
    avatar_url   TEXT         NOT NULL DEFAULT '',
    api_key_hash VARCHAR(64)  NOT NULL UNIQUE,
    last_seen_at TIMESTAMPTZ,
    status       VARCHAR(20)  NOT NULL DEFAULT 'active'
                     CHECK (status IN ('active', 'inactive', 'banned')),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bots_api_key_hash ON bots (api_key_hash);
CREATE INDEX IF NOT EXISTS idx_bots_status       ON bots (status);

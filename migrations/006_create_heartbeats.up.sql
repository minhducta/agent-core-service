CREATE TABLE IF NOT EXISTS heartbeats (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    bot_id     UUID        NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    status     VARCHAR(20) NOT NULL CHECK (status IN ('ok', 'degraded', 'error')),
    metadata   JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_heartbeats_bot_id     ON heartbeats (bot_id);
CREATE INDEX IF NOT EXISTS idx_heartbeats_created_at ON heartbeats (created_at DESC);

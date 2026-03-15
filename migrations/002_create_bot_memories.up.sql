CREATE TABLE IF NOT EXISTS bot_memories (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    bot_id      UUID        NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    type        VARCHAR(30) NOT NULL
                    CHECK (type IN ('fact', 'preference', 'instruction', 'context')),
    content     TEXT        NOT NULL,
    tags        TEXT[]      NOT NULL DEFAULT '{}',
    importance  INTEGER     NOT NULL DEFAULT 0,
    expires_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bot_memories_bot_id     ON bot_memories (bot_id);
CREATE INDEX IF NOT EXISTS idx_bot_memories_expires_at ON bot_memories (expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_bot_memories_importance ON bot_memories (importance DESC);

CREATE TABLE IF NOT EXISTS bot_policies (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    bot_id     UUID        NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    action     VARCHAR(100) NOT NULL,
    effect     VARCHAR(10)  NOT NULL CHECK (effect IN ('ALLOW', 'DENY')),
    conditions JSONB,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bot_policies_bot_id ON bot_policies (bot_id);
CREATE INDEX IF NOT EXISTS idx_bot_policies_action ON bot_policies (action);

CREATE TABLE IF NOT EXISTS todos (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    title         VARCHAR(255) NOT NULL,
    description   TEXT         NOT NULL DEFAULT '',
    status        VARCHAR(20)  NOT NULL DEFAULT 'pending'
                      CHECK (status IN ('pending', 'in_progress', 'done', 'cancelled')),
    priority      VARCHAR(10)  NOT NULL DEFAULT 'medium'
                      CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
    result        TEXT         NOT NULL DEFAULT '',
    due_date      TIMESTAMPTZ,
    assigned_to   UUID         REFERENCES bots(id) ON DELETE SET NULL,
    assigned_by   UUID         REFERENCES bots(id) ON DELETE SET NULL,
    dependency_id UUID         REFERENCES todos(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_todos_assigned_to ON todos (assigned_to);
CREATE INDEX IF NOT EXISTS idx_todos_assigned_by ON todos (assigned_by);
CREATE INDEX IF NOT EXISTS idx_todos_status       ON todos (status);
CREATE INDEX IF NOT EXISTS idx_todos_priority     ON todos (priority);

CREATE TABLE IF NOT EXISTS todo_checklist_items (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    todo_id     UUID        NOT NULL REFERENCES todos(id) ON DELETE CASCADE,
    content     TEXT        NOT NULL,
    is_checked  BOOLEAN     NOT NULL DEFAULT FALSE,
    order_index INTEGER     NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_checklist_items_todo_id ON todo_checklist_items (todo_id);

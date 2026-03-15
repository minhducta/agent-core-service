---
applyTo: "migrations/**"
---

# Migration Instructions

Using `golang-migrate` format: `NNN_description.up.sql` / `NNN_description.down.sql`.

## Rules

1. **Sequential numbering**: `001`, `002`, `003`...
2. **Idempotent**: use `CREATE TABLE IF NOT EXISTS`, `CREATE INDEX IF NOT EXISTS`.
3. **`down.sql`**: must perfectly reverse the `up.sql` — `DROP TABLE IF EXISTS`.
4. **UUID primary keys**: use `gen_random_uuid()` as default.
5. **Timestamps**: use `TIMESTAMPTZ NOT NULL DEFAULT NOW()` for `created_at`/`updated_at`.
6. **Foreign keys**: include `ON DELETE CASCADE` where child records should follow parent.
7. **Indexes**: add on foreign keys and frequently-queried columns.

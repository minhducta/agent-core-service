---
applyTo: "migrations/**"
---

# Migrations Instructions

The `migrations/` directory contains SQL migration files compatible with [golang-migrate](https://github.com/golang-migrate/migrate).

## File Naming Convention

```
NNN_short_description.up.sql    # Forward migration
NNN_short_description.down.sql  # Rollback migration
```

- `NNN` is a **zero-padded 3-digit sequential number** (e.g., `001`, `002`, `003`).
- `short_description` uses `snake_case` and describes *what* the migration does.
- Always create **both** `up` and `down` files for every migration.

## Current Schema

| Table | Migration | Key Columns |
|---|---|---|
| `bots` | `001` | `id`, `name`, `role`, `vibe`, `emoji`, `avatar_url`, `api_key_hash`, `last_seen_at`, `status`, `created_at`, `updated_at` |
| `bot_memories` | `002` | `id`, `bot_id`, `type`, `content`, `tags`, `importance`, `expires_at`, `created_at`, `updated_at` |
| `bot_skills` | `003` | `id`, `bot_id`, `name`, `description`, `usage_guide`, `created_at`, `updated_at` |
| `bot_policies` | `004` | `id`, `bot_id`, `action`, `effect`, `conditions`, `created_at`, `updated_at` |
| `todos` | `005` | `id`, `title`, `description`, `status`, `priority`, `result`, `due_date`, `assigned_to`, `assigned_by`, `dependency_id`, `created_at`, `updated_at` |
| `heartbeats` | `006` | `id`, `bot_id`, `status`, `metadata`, `created_at` |

## Up Migration Rules

1. **One logical change per migration** — don't bundle unrelated schema changes.
2. Use `CREATE TABLE IF NOT EXISTS` and `ADD COLUMN IF NOT EXISTS` to make migrations idempotent where possible.
3. Every new table must include:
   - `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
   - `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
   - `updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
4. Add indexes for **foreign keys** and **frequently filtered columns** (e.g., `bot_id`, `status`, `assigned_to`).
5. Use `TIMESTAMPTZ` (not `TIMESTAMP`) for all time columns — timezone-aware.
6. Use `TEXT` for variable-length strings; use `VARCHAR(N)` only when there is a hard business constraint on length.
7. Enum-like columns use `TEXT` with a `CHECK` constraint rather than PostgreSQL `ENUM` types — easier to extend without a full migration.
8. Boolean flags (e.g., `is_checked`) use `BOOLEAN NOT NULL DEFAULT FALSE`.

## Down Migration Rules

1. The `down` migration must **exactly undo** the `up` migration.
2. Use `DROP TABLE IF EXISTS` / `DROP COLUMN IF EXISTS`.
3. Drop indexes and constraints added in the `up` migration.

## Column Naming — Match Go Struct Tags

Column names **must match** the `db:"xxx"` struct tags in `internal/domain/`. When adding a new column:
1. Add the column to the migration SQL.
2. Add the corresponding field with a `db:` tag to the Go struct.
3. Update the repository SQL queries.

## Running Migrations

Migrations run **automatically** on application startup via `pkg/migration/migration.go`. The migration runner is called in `cmd/api/main.go` after the database connection is established.

Configuration in `config/config.yaml`:

```yaml
migration:
  enabled: true
  path: "migrations"
```

> Never edit an already-applied migration. Always add a new migration to make schema changes.

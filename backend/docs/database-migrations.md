# Database Migrations

This directory contains SQL migration files for the Budget Tracker database schema.

## Overview

The migration system uses a file-based approach where each `.sql` file represents a database schema change. All migration files are embedded into the compiled binary using Go's `//go:embed` directive, ensuring migrations are always available at runtime without external file dependencies.

Migrations are executed automatically when the application starts, applying any pending changes in version order.

## File Naming Convention

### Format

```
YYYY-MM-DD-NNN.sql
```

| Part   | Description                                           | Example      |
| ------ | ----------------------------------------------------- | ------------ |
| `YYYY` | Year when the migration was created                   | `2025`       |
| `MM`   | Month (zero-padded)                                   | `11`, `12`   |
| `DD`   | Day (zero-padded)                                     | `07`, `29`   |
| `NNN`  | Sequence number for same-day migrations (zero-padded) | `001`, `002` |

### Examples

| Filename             | Meaning                                       |
| -------------------- | --------------------------------------------- |
| `2025-11-29-001.sql` | First migration created on November 29, 2025  |
| `2025-11-29-002.sql` | Second migration created on November 29, 2025 |
| `2025-12-07-001.sql` | First migration created on December 7, 2025   |

## How to Add a New Migration

1. **Create a new `.sql` file** with the correct naming format:

   ```bash
   # Example: Creating a migration on December 7, 2025
   touch 2025-12-07-001.sql
   ```

2. **Write your SQL statements** in the file:

   ```sql
   -- Migration: 2025-12-07-001
   -- Description: Add user preferences table

   CREATE TABLE IF NOT EXISTS user_preferences (
       id INTEGER PRIMARY KEY AUTOINCREMENT,
       preference_key TEXT NOT NULL UNIQUE,
       preference_value TEXT,
       created_at DATETIME DEFAULT CURRENT_TIMESTAMP
   );
   ```

3. **Test locally** by running the application and verifying:
   - The migration applies without errors
   - The schema changes are correct
   - Existing functionality still works

4. **Commit with your PR** - the migration file will be embedded in the binary during build.

## Version Calculation

The filename is converted to an integer version number for ordering:

### Formula

```
version = year × 10,000,000 + month × 100,000 + day × 1,000 + sequence
```

### Examples

| Filename             | Version Number |
| -------------------- | -------------- |
| `2025-11-29-001.sql` | `20251129001`  |
| `2025-12-07-001.sql` | `20251207001`  |
| `2025-12-07-002.sql` | `20251207002`  |

This format ensures migrations are always sorted chronologically.

## Best Practices

### Do

- **One migration per PR/feature** - Keep migrations focused and reviewable
- **Use `IF NOT EXISTS` / `IF EXISTS`** - Makes migrations safer to re-run
- **Keep migrations idempotent when possible** - Running twice should have the same effect as running once
- **Add descriptive comments** - Explain what the migration does and why
- **Test on a copy of production data** - Before deploying, verify migrations work with real data

### Don't

- **Never modify existing migration files** - Once committed and deployed, migrations are immutable
- **Don't delete migration files** - They're part of the schema history
- **Avoid destructive operations without backups** - `DROP TABLE`, `DELETE`, etc. should be used carefully

## Backward Compatibility

### Legacy Version Mapping

Before the file-based migration system, migrations used simple integer versions (1, 2, 3). For databases created with the old system:

| Legacy Version | New Version | Description             |
| -------------- | ----------- | ----------------------- |
| 1              | 20251129001 | budget_limits table     |
| 2              | 20251129001 | expected_expenses table |
| 3              | 20251129001 | actual_expenses table   |

This mapping is handled automatically. If your database has legacy versions 1, 2, or 3 in `schema_migrations`, the system recognizes that `2025-11-29-001.sql` has already been applied.

## Example Migration File

```sql
-- Migration: 2025-12-07-001
-- Description: Add receipt tracking support

CREATE TABLE IF NOT EXISTS receipts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK(status IN ('pending', 'processing', 'completed', 'failed')),
    uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_receipts_status ON receipts(status);
```

## Troubleshooting

### Migration Failed

Migrations run in transactions. If a migration fails, the transaction is rolled back. Fix the SQL and restart the application.

### Checking Applied Migrations

```sql
SELECT version, description, applied_at
FROM schema_migrations
ORDER BY version;
```

# Bootstrapping sqlc for your project

This guide explains how to set up sqlc to generate type-safe Go code from your SQL queries.

## Prerequisites

- Go installed (>= 1.16)
- PostgreSQL (or your chosen database)
- [sqlc](https://github.com/sqlc-dev/sqlc) CLI

## 1. Install sqlc

Refer to [official documentation for installation instructions.](https://docs.sqlc.dev/en/latest/overview/install.html)

Or download a release from: https://github.com/sqlc-dev/sqlc/releases

Or directly to your `GOBIN`: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

## 2. Initialize sqlc configuration

In your project root, run:
```bash
sqlc init
```

This creates a `sqlc.yaml` file. You can also start from our example config:

```yaml
version: "2"
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "postgresql"
    gen:
      go:
        out: "internal/database"
        package: "database" # optional: default is directory name
```

## 3. Organize your SQL files

- Place DDL files (schema/migrations) in `sql/schema`, e.g., `001_users.sql`.
- Place query files in `sql/queries`, e.g., `users.sql`.
- Use named queries in comments:

```sql
-- name: GetUserByID :one
SELECT id, email, created_at
FROM users
WHERE id = $1;
```

## 4. Generate Go code

Run:
```bash
sqlc generate
```
This reads your SQL files and regenerates Go code under the path specified in `sqlc.yaml`.

## 5. Use the generated code

```go
import (
    "context"
    "database/sql"
    "yourmodule/internal/database"
)

func main() {
    conn, err := sql.Open("postgres", "postgresql://user:pass@localhost:5432/dbname?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    queries := database.New(conn)
    user, err := queries.GetUserByID(context.Background(), userID)
    // ...
}
```

## 6. Integrate into your workflow

- Re-run `sqlc generate` whenever you update schema or queries.
- Optionally, add `sqlc generate` to your CI/CD pipeline to catch mismatches early.

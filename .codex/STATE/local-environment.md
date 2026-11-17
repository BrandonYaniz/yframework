# Local Environment

Recorded 2026-11-03. This file is local-only (`.codex/` is gitignored).

## MariaDB (development)

| Setting | Value |
|---|---|
| Server | MariaDB 13.0.2 (Homebrew), listening on `127.0.0.1:3306` and `::1` |
| Database | `devframework` |
| User | `devframework@localhost` |
| Password | `YFRAMEWORK_DB_PASSWORD` in `.env.local` (gitignored) |
| Grants | `ALL PRIVILEGES ON devframework.*` |
| Charset | `utf8mb4` / `utf8mb4_general_ci` |
| sql_mode | `STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION` |
| Time zone | `SYSTEM` |

Verified on 2026-11-03: login succeeded, `devframework` exists, and a
create/drop table probe confirmed write access. The client is
`/opt/homebrew/bin/mariadb`.

## How tasks should use it

- Load `.env.local`; never hardcode credentials in source, tests, docs, commits,
  or log output.
- This server is the integration target from task 009 (connection config) and
  task 018 (migration, lock, and transaction suite) onward.
- Tests that mutate schema must create and drop their own objects. Do not drop,
  truncate, or reset the `devframework` database.
- A separate `devframework_test` database is recommended before task 018. It has
  not been created or granted yet; the current grant covers only `devframework`.
- The `mariadb` client is on `PATH`. Connectivity check:
  `mariadb -h 127.0.0.1 -u devframework -p devframework -e "SELECT 1"`.

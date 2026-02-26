Add SQL migration files here using this naming pattern:

- `0001_create_users_table.up.sql`
- `0001_create_users_table.down.sql`

On backend startup, `up` migrations are run automatically via `golang-migrate`.
You can run `down` migrations with `database.RollbackLastMigration(...)`.

CLI commands:

- `go run . migrate up`
- `go run . migrate down`

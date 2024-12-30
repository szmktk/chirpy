# Goose Migrations

Create a separate database user:

```
CREATE USER omg WITH PASSWORD 'wtf';
GRANT ALL PRIVILEGES ON DATABASE chirpy TO omg;
\du
```

To install `goose`:

```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

☝️ This is the method suggested by package maintainers, but I prefer to install with Homebrew anyway.

To run the migration use the connection string in schema directory (migration file must be there):

```
cd sql/schema
goose postgres "postgres://omg:wtf@localhost:5432/chirpy" up
goose postgres "postgres://omg:wtf@localhost:5432/chirpy" down
```

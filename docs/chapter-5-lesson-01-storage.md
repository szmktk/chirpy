# Storage

Make sure `psql` command is available:

```
brew install libpq
brew link --force libpq  # necessary to make psql accessible globally
psql --version
```

Export env vars to be able to connect by issuing `psql` with no parameters:

```
export PGHOST=localhost
export PGPORT=5432
export PGUSER=omg
export PGPASSWORD=wtf
export PGDATABASE=chirpy
```

While in `psql` shell, create a new database:

```
CREATE DATABASE chirpy;
```

Connect to the newly created database:

```
\c chirpy
```

Set the user password (necessary on Linux only):

```
ALTER USER postgres PASSWORD 'postgres';
```

One final verification:

```
SELECT version();
```

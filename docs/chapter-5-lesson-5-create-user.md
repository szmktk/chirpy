# Create User

Create a compose stack with a proper version of PostgreSQL & Adminer:

```
docker compose up -d
```

Make sure to have the following env vars in `.env` file at this point:

```
DB_URL="postgres://omg:wtf@localhost:5432/chirpy?sslmode=disable"
POSTGRES_DB=chirpy
POSTGRES_PASSWORD=wtf
POSTGRES_SERVER=localhost
POSTGRES_USER=omg
PLATFORM=dev
```

`POSTGRES_SERVER` points to `localhost` because we want to use Adminer from host machine 🙂

Using docker compose is convenient, as we do not have to manually create a database or admin user. The image will do it automatically given the right environment variables.

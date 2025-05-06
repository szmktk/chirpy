# Chirpy

A Go-based web service that provides API endpoints for user management, authentication, and chirps (posts).
Inspired by the [Learn HTTP Servers](https://www.boot.dev/courses/learn-http-servers) course on [boot.dev](https://boot.dev), and extended to suit my preferences.

![](assets/logo.png)


## Prerequisites

- Go 1.23 or later
- PostgreSQL database


## Configuration

The service requires the following environment variables:

- `DB_URL`: PostgreSQL connection string (e.g., `postgres://username:password@localhost:5432/chirpy`)
- `PORT`: Port number for the server (default: 8080)
- `JWT_SECRET`: Secret key for JWT token generation
- `POLKA_API_KEY`: API key for Polka webhook integration


## Installation

1. Clone the repository:
```bash
git clone https://github.com/szmktk/chirpy.git
cd chirpy
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your environment variables (you can create a `.env` file):
```bash
DB_URL=postgres://username:password@localhost:5432/chirpy
PORT=8080
JWT_SECRET=your-secret-key
POLKA_API_KEY=your-polka-api-key
```


## Running the Service

1. Start the service:
```bash
go run .
```

The server will start on the configured port (default: 8080).


## API Endpoints

### Authentication
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token

### Users
- `POST /api/users` - Create new user
- `PUT /api/users` - Update user (requires authentication)

### Chirps
- `POST /api/chirps` - Create new chirp (requires authentication)
- `GET /api/chirps` - Get all chirps
- `GET /api/chirps/{chirpID}` - Get specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete chirp (requires authentication)

### Admin
- `GET /admin/metrics` - Get server metrics
- `POST /admin/reset` - Reset server metrics

### Other
- `GET /api/healthz` - Health check endpoint
- `POST /api/polka/webhooks` - Polka webhook endpoint


## Development

### Running Tests

The project uses standard Go tooling. You can run tests using:

```bash
go test ./...
```

### Local Development

For a "hot-reload" feature during local development, you can use [`entr`](https://github.com/eradman/entr):

```bash
find . -type f -name '*.go' | entr -r go run .
```

When using this approach you need only the dependencies specified in [docker-compose.yaml](./docker-compose.yaml) file:

```
docker compose -f docker-compose.yaml up -d
```

### Running the Backend in a Docker Container

To run the backend component in a Docker container, simply run:

```
docker compose up -d
```

This command will also take [docker-compose.override.yaml](./docker-compose.override.yaml) file into account, which contains the backend service specification.

# chirpy

Implemented during [Learn HTTP Servers](https://www.boot.dev/courses/learn-http-servers) course on [boot.dev](https://boot.dev)


## pricing

As of 2024-10-07:

MONTHLY MEMBERSHIP

~~$49 USD~~

zł87/ month (zł1044 / year)

You want to learn a specific skill, or complete a single course


# todos


- [ ] 500 api responses should not return details to the caller, log them with level ERROR instead
- [ ] ask AI how to deduplicate the parsing & validation logic in HTTP handler functions
- [ ] add unit tests for internal/auth module
- [ ] handler_users_create.go: handle case when trying to create a user with the same email
    - [ ] differentiate between client & server errors here
- [ ] add integration tests for endpoint (one case per happy path and one per each error)
- [ ] try SQLC's `emit_json_tags: true` option instead of mapping JSON fields in code
- [ ] this error should complain about bad request (empty json payload)

    ```
    ❯ http post http://localhost:8080/api/users

    HTTP/1.1 500 Internal Server Error
    Access-Control-Allow-Origin: *
    Content-Length: 41
    Content-Type: application/json
    Date: Mon, 30 Dec 2024 11:04:24 GMT

    {
        "error": "Error decoding JSON body: EOF"
    }
    ```
- [ ] …

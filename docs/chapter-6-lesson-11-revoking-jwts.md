# Revoking JWTs

One of the main benefits of JWTs is that they're _stateless_. The server doesn't need to keep track of which users are logged in via JWT. The server just needs to issue a JWT to a user and the user can use that JWT to authenticate themselves. Statelessness is _fast and scalable_ because your server doesn't need to consult a database to see if a user is currently logged in.

However, that same benefit poses a potential problem. JWTs can't be revoked. If a user's JWT is stolen, there's no easy way to stop the JWT from being used. JWTs are just a signed string of text.

The JWTs we've been using so far are more specifically _access tokens_. Access tokens are used to authenticate a user to a server, and they provide _access_ to protected resources. Access tokens are:

- Stateless
- Short-lived (15m-24h)
- Irrevocable

They _must_ be short-lived because they can't be revoked. The shorter the lifespan, the more secure they are. Trouble is, this can create a poor user experience. We don't want users to have to log in every 15 minutes.

## A Solution: Refresh Tokens

Refresh tokens don't provide access to resources directly, but they can be used to get new access tokens. Refresh tokens are much longer lived, and importantly, they _can_ be revoked. They are:

- Stateful
- Long-lived (24h-60d)
- Revocable

Now we get the best of both worlds! Our endpoints and servers that provide access to protected resources can use access tokens, which are fast, stateless, simple, and scalable. On the other hand, refresh tokens are used to keep users logged in for longer periods of time, and they can be revoked if a user's access token is compromised.

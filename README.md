# Token Authentication

Token Authentication middleware for [Fiber](https://github.com/gofiber/fiber) that provides a basic token authentication. It calls the next handler for valid token and [401 Unauthorized](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401) for a missing or an invalid token.

### How to use

- [In memory](#in-memory)
- [Databases](#databases)
- [Redis](#redis)

### IN MEMORY

Example for in memory token storage:

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/klipitkas/tokenauth"
)

func main() {
	app := fiber.New()

	app.Use(tokenauth.New(tokenauth.Config{
		Tokens: map[string]tokenauth.Claims{
			"token":  {"user": "john", "id": "42"},
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(tokenauth.Claims)
		return c.SendString("Hello, " + claims["user"] + " 👋!")
	})

	_ = app.Listen(":3000")
}
```

Trying to access the route without a token:

```shell
$ curl http://localhost:3000
Unauthorized
```

Trying to access the route after providing a token:

```shell
$ curl -H 'Authorization: Bearer token' http://localhost:3000
Hello, john 👋!
```

### DATABASES

- Coming soon!

### REDIS

- Coming soon!

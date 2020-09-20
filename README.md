# Token Authentication

Token Authentication middleware for [Fiber](https://github.com/gofiber/fiber) that provides a basic token authentication. It calls the next handler for valid token and [401 Unauthorized](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401) or a custom response for missing or invalid token.

### Table of Contents

- [Signatures](#signatures)
- [Examples](#examples)
- [Config](#config)
- [Default Config](#default-config)

### Signatures

```go
func New(config Config) fiber.Handler
```

### Examples

Import the middleware package as shown below:

```go
import (
  "github.com/klipitkas/tokenauth"
)
```

After you initiate your Fiber app, you can use the following possibilities:

```go
claims := Claims{"user": "john", "email": "john@example.com", "id": "42"}

// Provide a minimal config
app.Use(tokenauth.New(tokenauth.Config{
	Tokens: map[string]string{
		"1HTWgKFX6zaCb5pwpH4RKJz7": claims,
	},
}))

// Or extend your config for customization
app.Use(tokenauth.New(tokenauth.Config{
	Tokens: map[string]string{
		"1HTWgKFX6zaCb5pwpH4RKJz7": claims,
	},
	Realm: "Forbidden",
	Authorizer: func(token string) Claims {
		if token == "1HTWgKFX6zaCb5pwpH4RKJz7" {
			return claims
		}
		return nil
	},
	Unauthorized: func(c *fiber.Ctx) error {
		return c.SendFile("./unauthorized.html")
	},
	ContextKey: "_claims"
}))
```

### Config

```go
// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Tokens defines the mappings between the tokens and their claims.
	//
	// Required. Default: map[string]Claims{}
	Tokens map[string]Claims

	// Realm is a string to define realm attribute of BasicAuth.
	// the realm identifies the system to authenticate against
	// and can be used by clients to save credentials
	//
	// Optional. Default: "Restricted".
	Realm string

	// Authorizer defines a function you can pass
	// to check the token however you want.
	// It will be called with the provided token
	// and is expected to return the identifier of
	// what the token matched with or empty string
	// if it did not match with anything.
	//
	// Optional. Default: nil.
	Authorizer func(string) Claims

	// Unauthorized defines the response body for unauthorized responses.
	// By default it will return with a 401 Unauthorized and the correct WWW-Auth header
	//
	// Optional. Default: nil
	Unauthorized fiber.Handler

	// ContextKey is the key to store user information to locals.
	//
	// Optional. Default: "claims"
	ContextKey string
}
```

### Default Config

```go
// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:         nil,
	Tokens:       map[string]Claims{},
	Realm:        "Restricted",
	Authorizer:   nil,
	Unauthorized: nil,
	ContextKey:   "claims",
}
```

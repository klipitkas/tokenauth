package tokenauth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Claims represent the available claims that are connected with a token.
type Claims map[string]string

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

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

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:         nil,
	Authorizer:   nil,
	Unauthorized: nil,
	ContextKey:   "claims",
}

// New creates a new middleware handler
func New(config Config) fiber.Handler {
	cfg := config
	// Set default values
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}
	if cfg.Authorizer == nil {
		cfg.Authorizer = func(token string) Claims {
			if len(token) == 0 {
				return nil
			}
			return Claims{}
		}
	}
	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = ConfigDefault.ContextKey
	}
	// Return new handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}
		// Get authorization header
		auth := c.Get(fiber.HeaderAuthorization)
		// Check if header is valid
		if len(auth) > 6 && strings.ToLower(auth[:6]) == "bearer" {
			// Get the token
			token := auth[7:]
			if len(token) == 0 {
				return cfg.Unauthorized(c)
			}
			claims := cfg.Authorizer(token)
			if claims != nil {
				c.Locals(cfg.ContextKey, claims)
				return c.Next()
			}
		}
		// Authentication failed
		return cfg.Unauthorized(c)
	}
}

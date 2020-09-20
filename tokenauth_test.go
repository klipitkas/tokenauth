package tokenauth

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// go test -run Test_TokenAuth_Next
func Test_TokenAuth_Next(t *testing.T) {
	app := fiber.New()
	app.Use(New(Config{
		Next: func(_ *fiber.Ctx) bool {
			return true
		},
	}))

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_Middleware_TokenAuth(t *testing.T) {
	app := fiber.New()

	cfg := Config{
		Tokens: map[string]Claims{"1HTWgKFX6zaCb5pwpH4RKJz7": {"user": "john", "email": "john@example.com", "id": "42"}},
	}

	app.Use(New(cfg))
	app.Get("/testauth", func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(Claims)
		bytes, _ := json.Marshal(claims)
		return c.Send(bytes)
	})

	tests := []struct {
		url        string
		statusCode int
		token      string
		claims     Claims
	}{
		{
			url:        "/testauth",
			statusCode: 200,
			token:      "1HTWgKFX6zaCb5pwpH4RKJz7",
			claims:     Claims{"user": "john", "email": "john@example.com", "id": "42"},
		},
		{
			url:        "/testauth",
			statusCode: 401,
			token:      "",
			claims:     nil,
		},
		{
			url:        "/testauth",
			statusCode: 401,
			token:      "123456",
			claims:     nil,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/testauth", nil)
		req.Header.Add("Authorization", "Bearer "+tt.token)
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err)

		body, err := ioutil.ReadAll(resp.Body)

		utils.AssertEqual(t, nil, err)
		utils.AssertEqual(t, tt.statusCode, resp.StatusCode)

		// Only check body if statusCode is 200
		if tt.statusCode == 200 {
			c := Claims{}
			_ = json.Unmarshal(body, &c)
			utils.AssertEqual(t, tt.claims, c)
		}
	}
}

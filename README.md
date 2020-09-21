# Token Authentication

Token Authentication middleware for [Fiber](https://github.com/gofiber/fiber) that provides a basic token authentication. It calls the next handler for valid token and [401 Unauthorized](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401) for a missing or an invalid token.

### How to use

- [In memory](#in-memory)
- [Databases](#databases)
  - [Postgres](#postgres)
- [Redis](#redis)

### IN MEMORY

Example for in memory token storage:

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/klipitkas/tokenauth"
)

var Tokens = map[string]tokenauth.Claims{
	"token": {"user": "john", "email": "john@example.com", "id": "42"},
}

func main() {
	app := fiber.New()

	app.Use(tokenauth.New(tokenauth.Config{
		Authorizer: func(s string) tokenauth.Claims {
			claims, exist := Tokens[s]
			if !exist {
				return nil
			}
			return claims
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(tokenauth.Claims)
		return c.SendString("Hello, " + claims["user"] + " 👋!")
	})

	_ = app.Listen(":3000")
}
```

Try to access the route without a token:

```shell
$ curl http://localhost:3000
Unauthorized
```

Try to access the route after providing a token:

```shell
$ curl -H 'Authorization: Bearer token' http://localhost:3000
Hello, john 👋!
```

### DATABASES

#### POSTGRES

Use docker to launch a new ephemeral postgres database:

```shell
$ docker run -p 5432:5432 -e POSTGRES_PASSWORD=tokenauth -e POSTGRES_USER=tokenauth -e POSTGRES_DB=tokenauth -d postgres
```

Try the following example:

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/klipitkas/tokenauth"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "tokenauth"
	pass   = "tokenauth"
	dbname = "tokenauth"
)

var sqlMigration string = `
DROP TABLE IF EXISTS user_tokens;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
	id INT,
	email TEXT,
	PRIMARY KEY(id)
);

CREATE TABLE user_tokens (
	id INT,
	user_id INT,
	token TEXT,
	PRIMARY KEY(id),
	UNIQUE(token),
	CONSTRAINT fk_customer FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);`

var sqlSeeder string = `
	INSERT INTO users VALUES (1, 'john@example.com');
	INSERT INTO users VALUES (2, 'jim@example.com');

	INSERT INTO user_tokens VALUES (1, 1, 'token');
	INSERT INTO user_tokens VALUES (2, 2, 'fgu9KILmznhLtQgmr3');
`

func main() {
	app := fiber.New()

	// Connect to the database.
	db, err := dbConnect(host, user, pass, dbname, port)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	// Run the migrations.
	if _, err = db.Exec(sqlMigration); err != nil {
		log.Fatalf("Failed to run the migrations: %v", err)
	}

	// Run the seeders.
	if _, err = db.Exec(sqlSeeder); err != nil {
		log.Fatalf("Failed to run the migrations: %v", err)
	}

	// Use our custom authorizer to verify tokens from the DB.
	app.Use(tokenauth.New(tokenauth.Config{Authorizer: dbTokenAuthorizer}))

	// Our protected route
	app.Get("/", func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(tokenauth.Claims)
		return c.SendString("Hello, user with ID: " + claims["user_id"] + " 👋!")
	})

	_ = app.Listen(":3000")

	defer db.Close()
}

func dbConnect(host, user, pass, dbname string, port int) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)

	// Validate the connection arguments
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping the db: %v", err)
	}

	return db, nil
}

func dbTokenAuthorizer(token string) tokenauth.Claims {
	// Connect to the database.
	db, err := dbConnect(host, user, pass, dbname, port)
	if err != nil {
		fmt.Printf("connect to db: %v", err)
		return nil
	}
	defer db.Close()

	// Get the user id if the token is valid.
	// You can do more things here, such as JOIN statements with the users
	// table and get more information for the user.
	rows, err := db.Query("SELECT user_id FROM user_tokens WHERE token = $1", token)
	if err != nil || rows.Err() != nil {
		fmt.Printf("run query for tokens: %v", err)
		return nil
	}

	// Set a default empty user ID.
	userID := ""
	for rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			fmt.Printf("get token rows: %v", err)
			return nil
		}
	}

	// Check if we got a valid (non-zero length) user id.
	if len(userID) > 0 {
		return tokenauth.Claims{"user_id": userID}
	}

	return nil
}
```

Try to access the route without a token:

```shell
$ curl http://localhost:3000
Unauthorized
```

Try to access the route after providing a token:

```shell
$ curl -H 'Authorization: Bearer token' http://localhost:3000
Hello, user with ID: 1 👋!
```

```shell
$ curl -H 'Authorization: Bearer fgu9KILmznhLtQgmr3' http://localhost:3000
Hello, user with ID: 2 👋!
```

### REDIS

Use docker to launch a new ephemeral redis instance:

```shell
$ docker run -p 6379:6379 --name redis -d redis
```

Try the following example:

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/klipitkas/tokenauth"

	"context"

	"github.com/go-redis/redis/v8"
)

func main() {
	app := fiber.New()

	// Use our custom authorizer to verify tokens from Redis.
	app.Use(tokenauth.New(tokenauth.Config{Authorizer: redisTokenAuthorizer}))

	// Our protected route
	app.Get("/", func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(tokenauth.Claims)
		return c.SendString("Hello, user with details: " + claims["user"] + " !👋")
	})

	_ = app.Listen(":3000")
}

func redisTokenAuthorizer(token string) tokenauth.Claims {
	// Create a new redis client and connect to Redis instance.
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redis.Close()

	// Test if connection to Redis works.
	if _, err := redis.Ping(context.Background()).Result(); err != nil {
		return nil
	}

	// Fetch the user details using the token as a key.
	user, err := redis.Get(context.Background(), token).Result()
	if err != nil {
		return nil
	}

	if len(user) > 0 {
		return tokenauth.Claims{"user": user}
	}

	return nil
}
```

Now you can create a new token using `redis-cli`:

```shell
$  docker exec -it redis redis-cli
127.0.0.1:6379> set token test
OK
```

Try to access the route without a token:

```shell
$ curl http://localhost:3000
Unauthorized
```

Try to access the route after providing a token:

```shell
$ curl -H 'Authorization: Bearer token' http://localhost:3000
Hello, user with details: test 👋!
```

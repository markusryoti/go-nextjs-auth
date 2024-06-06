package main

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/markusryoti/next-js-go-auth/internal/auth"
)

type RegisterCmd struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginCmd struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/register", func(c fiber.Ctx) error {
		cmd := new(RegisterCmd)

		err := c.Bind().Body(cmd)
		if err != nil {
			return err
		}

		tokenRes, err := auth.RegisterUser(cmd.Email, cmd.Password)
		if err != nil {
			return err
		}

		return c.JSON(tokenRes)
	})

	app.Post("/login", func(c fiber.Ctx) error {
		cmd := new(LoginCmd)

		err := c.Bind().Body(cmd)
		if err != nil {
			return err
		}

		tokenRes, err := auth.LoginUser(cmd.Email, cmd.Password)
		if err != nil {
			return err
		}

		return c.JSON(tokenRes)
	})

	app.Get("/current-user", func(c fiber.Ctx) error {
		sessionId, err := getSessionId(c)
		if err != nil {
			return err
		}

		currentUserResponse, err := auth.GetCurrentUser(sessionId)
		if err != nil {
			return err
		}

		return c.JSON(currentUserResponse)
	}, authenticated)

	app.Get("/logout", func(c fiber.Ctx) error {
		sessionId, err := getSessionId(c)
		if err != nil {
			return err
		}

		auth.Logout(sessionId)

		return c.SendString("ok")
	}, authenticated)

	log.Fatal(app.Listen(":5000"))
}

func authenticated(c fiber.Ctx) error {
	headers := c.GetReqHeaders()

	authHeaders, ok := headers["Authorization"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("no authorization headers")
	}

	if len(authHeaders) != 1 {
		return c.Status(fiber.StatusBadRequest).SendString("bad number of auth headers")
	}

	authHeader := authHeaders[0]

	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid auth header")
	}

	token := parts[1]

	claims, err := auth.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	c.Locals("claims", claims)

	return c.Next()
}

func getSessionId(c fiber.Ctx) (string, error) {
	claims, ok := c.Locals("claims").(auth.MyCustomClaims)
	if !ok {
		return "", errors.New("couldn't convert to custom claims")
	}

	return claims.Session, nil
}

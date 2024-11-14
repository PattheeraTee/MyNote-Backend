package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydb"       // as defined in docker-compose.yml
)

func checkMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64) // เก็บ user_id จากโทเค็น (float64 เพราะ JSON decode เป็น float)

	// เก็บ user_id ใน Context เพื่อใช้ในแฮนด์เลอร์อื่น ๆ
	c.Locals("user_id", uint(userID))

	return c.Next()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	db := setupDatabase()
	app := fiber.New()

	app.Post("/register", func(c *fiber.Ctx) error {
		return registerHandler(db, c)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return loginHandler(db, c)
	})

	app.Use(checkMiddleware)

	app.Get("/user/:id", func(c *fiber.Ctx) error {
		return getUserHandler(db, c)
	})

	app.Put("/user/:id", func(c *fiber.Ctx) error {
		return changeUsernameHandler(db, c)
	})

	app.Post("/forgot-password", func(c *fiber.Ctx) error {
		return forgotPasswordHandler(db, c)
	})

	app.Post("/reset-password", func(c *fiber.Ctx) error {
		return resetPasswordHandler(db, c)
	})

	app.Listen(":8000")
}

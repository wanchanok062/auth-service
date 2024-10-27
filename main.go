package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/wanchanok6698/web-auth/api/v1/routes"
	"github.com/wanchanok6698/web-auth/config"
)

func main() {
	config.ConnectDB()
	if config.DB == nil {
		log.Fatal("MongoDB connection failed")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load environment variables: ", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	app := fiber.New()
	app.Use(logger.New())

	app.Use(limiter.New(limiter.Config{
		Max:        15,
		Expiration: 30 * time.Second,
	}))

	prefix := os.Getenv("API_PREFIX")
	routes.AuthRoutes(prefix, app)

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

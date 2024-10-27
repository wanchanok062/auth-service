package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/wanchanok6698/web-auth/api/middleware"
	"github.com/wanchanok6698/web-auth/api/v1/controllers"
	"github.com/wanchanok6698/web-auth/api/v1/models"
	"github.com/wanchanok6698/web-auth/api/v1/services"
)

func AuthRoutes(prefix string, app *fiber.App) {
	authService, err := services.NewAuthService()
	if err != nil {
		log.Fatalf("Failed to initialize auth  service: %v", err)
	}

	authController := controllers.NewAuthController(*authService)

	routeGrop := app.Group(prefix)
	routeGrop.Get("user/:id", authController.GetUserByID)
	routeGrop.Post("register", middleware.ValidateData(&models.RegisterRequest{}), authController.RegisterUser)
}

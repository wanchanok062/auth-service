package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/wanchanok6698/web-auth/api/v1/models"
	"github.com/wanchanok6698/web-auth/api/v1/services"
	"github.com/wanchanok6698/web-auth/util"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (ac *AuthController) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return util.HandleError(c, "ไม่พบ ID parameter", "ต้องมี ID ", fiber.StatusBadRequest)
	}

	resultChan := make(chan *models.GetUser)
	errorChan := make(chan error)

	go func() {
		user, err := ac.service.GetUserByID(context.Background(), id)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- user
	}()

	select {
	case user := <-resultChan:
		return util.HandleSuccess(c, "ดึงข้อมูลสำเร็จ", user)
	case err := <-errorChan:
		return util.HandleError(c, "ดึงข้อมูล user ไม่สำเร็จ", err.Error(), fiber.StatusInternalServerError)
	}
}

func (ac *AuthController) RegisterUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return util.HandleError(c, "Invalid request data", err.Error(), fiber.StatusBadRequest)
	}

	userId, token, err := ac.service.RegisterUser(context.Background(), user)
	if err != nil {
		if err.Error() == "username already taken" {
			return util.HandleError(c, "register ไม่สำเร็จ", err.Error(), fiber.StatusConflict)
		}
		return util.HandleError(c, "register ไม่สำเร็จ", err.Error(), fiber.StatusInternalServerError)
	}

	responseData := struct {
		UserId string `json:"userId"`
		Token  string `json:"token"`
	}{
		UserId: userId,
		Token:  token,
	}

	return util.HandleSuccess(c, "registered สำเร็จ", responseData)
}

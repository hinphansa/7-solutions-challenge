package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hinphansa/7-solutions-challenge/internal/ports"
	"github.com/hinphansa/7-solutions-challenge/pkg/logger"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	log         logger.Logger
	authsvc     ports.AuthService
	userHandler *UserHandler
}

func NewAuthHandler(log logger.Logger, authService ports.AuthService, userHandler *UserHandler) *AuthHandler {
	log = log.WithFields(logrus.Fields{
		"module": "auth-handler",
	})
	return &AuthHandler{
		log:         log,
		authsvc:     authService,
		userHandler: userHandler,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Login
// @Summary Login
// @Description Login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := MustValid(c, &req); err != nil {
		return err
	}

	token, err := h.authsvc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

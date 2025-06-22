package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hinphansa/7-solutions-challenge/internal/domain"
	"github.com/hinphansa/7-solutions-challenge/internal/ports"
	"github.com/hinphansa/7-solutions-challenge/pkg/logger"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
)

/* -------------------------------------------------------------------------- */
/*                                 UserHandler                                */
/* -------------------------------------------------------------------------- */

type UserHandler struct {
	log     logger.Logger
	usersvc ports.UserService
}

func NewUserHandler(log logger.Logger, userService ports.UserService) *UserHandler {
	log = log.WithFields(logrus.Fields{
		"module": "user-handler",
	})
	return &UserHandler{log: log, usersvc: userService}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=3"`
}

// Register
// @Summary Register a new user
// @Description Register a new user
// @Tags user
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := MustValid(c, &req); err != nil {
		return err
	}

	id, err := h.usersvc.Register(c.Context(), &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}

// GetUser by id
// @Summary Get user by id
// @Description Get user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	bsonId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.usersvc.GetByID(c.Context(), bsonId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

type UpdateUserRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Name  string `json:"name" validate:"omitempty,min=3"`
}

// UpdateUser by id
// @Summary Update user by id
// @Description Update user's email and name by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "Update request"
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	var (
		req    UpdateUserRequest
		bsonId bson.ObjectID
	)

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	bsonId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := MustValid(c, &req); err != nil {
		return err
	}

	if err := h.usersvc.Update(c.Context(), bsonId, &domain.User{
		Email: req.Email,
		Name:  req.Name,
	}); err != nil {
		h.log.Errorf("Failed to update user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// DeleteUser by id
// @Summary Delete user by id
// @Description Delete user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	bsonId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := h.usersvc.Delete(c.Context(), bsonId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// ListAllUsers
// @Summary List all users
// @Description List all users
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} domain.User
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid offset",
		})
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit",
		})
	}

	var users []domain.User
	if offset == 0 && limit == 0 {
		users, err = h.usersvc.GetAll(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to list users",
			})
		}
	} else {
		users, err = h.usersvc.List(c.Context(), &ports.Pagination{
			Offset: int64(offset),
			Limit:  int64(limit),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to list users",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
	})
}

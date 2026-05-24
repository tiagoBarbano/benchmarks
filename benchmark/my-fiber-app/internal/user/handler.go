package user

import (
	dto "my-fiber-app/internal/user/dto"
	"my-fiber-app/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func RegisterRoutes(router fiber.Router) {
	repo := NewRepository()
	service := NewService(repo)
	handler := NewHandler(service)

	group := router.Group("/users")
	group.Post("/", handler.CreateUser)
	group.Get("/", handler.GetUsers)
	group.Get("/:id", handler.GetUserByID)
	group.Put("/:id", handler.UpdateUser)
	group.Delete("/:id", handler.DeleteUser)
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	req, err := utils.ValidateRequestBody[dto.CreateUserRequest](c)
	if err != nil {
		return err // já retorna o JSON com os erros
	}

	user := User{
		Name:  req.Name,
		Email: req.Email,
		CEP:   req.CEP,
	}

	created, err := h.service.Create(c.UserContext(), &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *Handler) GetUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAll(c.UserContext())
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(users)
}

func (h *Handler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.service.GetByID(c.UserContext(), id)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	if user == nil {
		return fiber.ErrNotFound
	}
	return c.JSON(user)
}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	if err := c.BodyParser(&user); err != nil {
		return fiber.ErrBadRequest
	}
	updated, err := h.service.Update(c.UserContext(), id, &user)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(updated)
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.Delete(c.UserContext(), id)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(fiber.StatusNoContent)
}

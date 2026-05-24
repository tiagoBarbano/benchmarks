package utils

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validator = validator.New()

type ValidationErrorDetail struct {
	Loc  []string `json:"loc"`
	Msg  string   `json:"msg"`
	Type string   `json:"type"`
}

func ValidateRequestBody[T any](c *fiber.Ctx) (*T, error) {
	var payload T
	if err := c.BodyParser(&payload); err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": []ValidationErrorDetail{
				{
					Loc:  []string{"body"},
					Msg:  "JSON inválido",
					Type: "value_error.json",
				},
			},
		})
	}

	if err := Validator.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var details []ValidationErrorDetail

		for _, e := range validationErrors {
			field, _ := reflect.TypeOf(payload).FieldByName(e.StructField())
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = e.Field()
			}

			details = append(details, ValidationErrorDetail{
				Loc:  []string{"body", jsonTag},
				Msg:  friendlyMessage(e),
				Type: "value_error." + e.Tag(),
			})
		}

		return nil, c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"detail": details,
		})
	}

	return &payload, nil
}

func friendlyMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "campo obrigatório"
	case "email":
		return "email inválido"
	case "min":
		return "mínimo de " + e.Param() + " caracteres"
	case "len":
		return "deve ter exatamente " + e.Param() + " caracteres"
	case "numeric":
		return "deve conter apenas números"
	default:
		return "valor inválido"
	}
}

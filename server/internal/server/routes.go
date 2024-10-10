package server

import (
	"net/http"
	"server/internal/email"
	"server/internal/types"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	api := s.App.Group("/api")
	api.Get("/", s.HelloWorldHandler)
	api.Get("/send-email", s.sendEmail)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) sendEmail(c *fiber.Ctx) error {
	var emailToSend types.Email
	if err := c.BodyParser(&emailToSend); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := email.SendEmail(emailToSend)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Email sent"})
}

func (s *FiberServer) login(c *fiber.Ctx) error {
	return c.SendString("Login")
}

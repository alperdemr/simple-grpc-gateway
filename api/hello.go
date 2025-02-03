package api

import "github.com/gofiber/fiber/v2"

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Message string `json:"message"`
}

func (s *Server) Hello(c *fiber.Ctx) error {
	req := new(HelloRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(HelloResponse{
		Message: "Hello, " + req.Name,
	})
}


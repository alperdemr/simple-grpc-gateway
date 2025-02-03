package api

import (
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App

}

func NewServer() *Server {
	server := &Server{
		app: fiber.New(),
	}
	server.SetupRoutes()
	return server
}


func (s *Server) SetupRoutes() {
	v1 := s.app.Group("/api/v1")
	v1.Post("/hello",s.Hello)
}


func (s *Server) Start(address string) {
	s.app.Listen(address)
}
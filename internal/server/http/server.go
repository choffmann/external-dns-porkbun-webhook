package http

import (
	"context"
	"fmt"

	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	cfg *config.Config
	svc *service.Service
}

func NewServer(cfg *config.Config, svc *service.Service) *Server {
	return &Server{
		cfg: cfg,
		svc: svc,
	}
}

func (s *Server) Run(ctx context.Context) error {
	app := fiber.New()
	app.Mount("/", s.router())

	app.Use(s.healthCheck())

	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down server...")
		app.Shutdown()
	}()

	return app.Listen("127.0.0.1:8888")
}

package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

func (s *Server) healthCheck() func(*fiber.Ctx) error {
	return healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint:  "/healthz",
		ReadinessEndpoint: "/ready",
	})
}

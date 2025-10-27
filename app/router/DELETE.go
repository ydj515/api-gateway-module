package router

import (
	"api-gateway-module/app/client"
	"api-gateway-module/config"

	"github.com/gofiber/fiber/v2"
)

type delete struct {
	cfg    config.Router
	client client.HttpClient
}

func AddDelete(cfg config.Router, client client.HttpClient) func(c *fiber.Ctx) error {
	r := delete{cfg: cfg, client: client}
	return r.handleRequest
}

func (r delete) handleRequest(c *fiber.Ctx) error {
	apiResult, err := r.client.DELETE(r.cfg.Path, r.cfg, c.Request().Body())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(apiResult)
}

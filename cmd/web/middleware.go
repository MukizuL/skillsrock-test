package main

import (
	"github.com/gofiber/fiber/v2"
)

// logRequest middleware is used to lag every request
func (app *application) logRequest(c *fiber.Ctx) error {
	var (
		ip     = c.IP()
		proto  = c.Protocol()
		method = c.Method()
		uri    = c.BaseURL()
	)

	app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

	return c.Next()
}

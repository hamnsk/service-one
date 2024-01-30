package main

import "github.com/gofiber/fiber/v2"

type app struct {
	httpServer *fiber.App
	tp         AppTracer
}

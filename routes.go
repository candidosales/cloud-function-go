package main

import (
  "fmt"
  "github.com/gofiber/fiber"
)

func  (app *AppConfig) setupRoutes() {
  app.fiber.Get("/ping", app.ping)
  app.fiber.Post("/", app.createCF)
}

func (app *AppConfig)  ping(c *fiber.Ctx) {
  c.JSON(fiber.Map{"pong": "ok"})
}

func (app *AppConfig) createCF(c *fiber.Ctx) {
  createCFRequest := &CreateCFRequest{}

  if err := c.BodyParser(&createCFRequest); err != nil {
    fmt.Printf("err[%#v] \n", err)
  }

  result, err := app.createCloudFunction(createCFRequest)

  if err != nil {
    c.Status(503).Send(err)
    return
  }

  c.JSON(fiber.Map{"message": result })
}

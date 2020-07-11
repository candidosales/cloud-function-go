package main

import (
	"fmt"

	"github.com/gofiber/fiber"
)

func (app *AppConfig) setupRoutes() {
	app.fiber.Get("/ping", app.ping)
	app.fiber.Post("/", app.createCF)
	app.fiber.Get("/", app.getCF)
	app.fiber.Delete("/", app.deleteCF)
	app.fiber.Post("/policy", app.setPolicy)
}

func (app *AppConfig) ping(c *fiber.Ctx) {
	if err := c.JSON(fiber.Map{"pong": "ok"}); err != nil {
		c.Next(err)
	}
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

	if err := c.JSON(fiber.Map{"message": result}); err != nil {
		c.Next(err)
	}
}

func (app *AppConfig) getCF(c *fiber.Ctx) {
	functionID := c.Query("function_id")

	result, err := app.getCloudFunction(functionID)
	if err != nil {
		c.Status(400).Send(err)
		return
	}

	if err := c.JSON(fiber.Map{"cloudFunction": result}); err != nil {
		c.Next(err)
	}
}

func (app *AppConfig) deleteCF(c *fiber.Ctx) {
	functionID := c.Query("function_id")

	result, err := app.deleteCloudFunction(functionID)
	if err != nil {
		c.Status(400).Send(err)
		return
	}

	if err := c.JSON(fiber.Map{"cloudFunction": result}); err != nil {
		c.Next(err)
	}
}

func (app *AppConfig) setPolicy(c *fiber.Ctx) {
	getCFRequest := &GetCFRequest{}

	if err := c.BodyParser(&getCFRequest); err != nil {
		fmt.Printf("err[%#v] \n", err)
	}

	result, err := app.setIAMPolicy(getCFRequest.FunctionID)

	if err != nil {
		c.Status(400).Send(err)
		return
	}

	if err := c.JSON(fiber.Map{"policy": result}); err != nil {
		c.Next(err)
	}
}

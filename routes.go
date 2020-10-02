package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (app *AppConfig) setupRoutes() {
	app.fiber.Get("/ping", app.ping)
	app.fiber.Post("/", app.createCF)
	app.fiber.Get("/", app.getCF)
	app.fiber.Delete("/", app.deleteCF)
	app.fiber.Post("/policy", app.setPolicy)
}

func (app *AppConfig) ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"pong": "ok"})
}

func (app *AppConfig) createCF(c *fiber.Ctx) error {
	createCFRequest := &CreateCFRequest{}

	if err := c.BodyParser(&createCFRequest); err != nil {
		fmt.Printf("err[%#v] \n", err)
	}

	result, err := app.createCloudFunction(createCFRequest)
	if err != nil {
		return c.Status(503).Send([]byte(err.Error()))
	}

	return c.JSON(fiber.Map{"message": result})
}

func (app *AppConfig) getCF(c *fiber.Ctx) error {
	functionID := c.Query("function_id")

	result, err := app.getCloudFunction(functionID)
	if err != nil {
		return c.Status(400).Send([]byte(err.Error()))
	}

	return c.JSON(fiber.Map{"cloudFunction": result})
}

func (app *AppConfig) deleteCF(c *fiber.Ctx) error {
	functionID := c.Query("function_id")

	result, err := app.deleteCloudFunction(functionID)
	if err != nil {
		return c.Status(400).Send([]byte(err.Error()))
	}

	return c.JSON(fiber.Map{"cloudFunction": result})
}

func (app *AppConfig) setPolicy(c *fiber.Ctx) error {
	getCFRequest := &GetCFRequest{}

	if err := c.BodyParser(&getCFRequest); err != nil {
		fmt.Printf("err[%#v] \n", err)
	}

	result, err := app.setIAMPolicy(getCFRequest.FunctionID)
	if err != nil {
		return c.Status(400).Send([]byte(err.Error()))
	}

	return c.JSON(fiber.Map{"policy": result})
}

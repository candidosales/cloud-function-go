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
	if err := c.JSON(fiber.Map{"pong": "ok"}); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusForbidden)
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

	if err := c.JSON(fiber.Map{"message": result}); err != nil {
		return err
	}

	return err
}

func (app *AppConfig) getCF(c *fiber.Ctx) error {
	functionID := c.Query("function_id")

	result, err := app.getCloudFunction(functionID)
	if err != nil {
		return c.Status(400).Send([]byte(err.Error()))
	}

	if err := c.JSON(fiber.Map{"cloudFunction": result}); err != nil {
		return err
	}

	return err
}

func (app *AppConfig) deleteCF(c *fiber.Ctx) error {
	functionID := c.Query("function_id")

	result, err := app.deleteCloudFunction(functionID)
	if err != nil {
		return c.Status(400).Send([]byte(err.Error()))
	}

	if err := c.JSON(fiber.Map{"cloudFunction": result}); err != nil {
		return err
	}

	return err
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

	if err := c.JSON(fiber.Map{"policy": result}); err != nil {
		return err
	}

	return err
}

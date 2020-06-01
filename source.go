package main

import (
	"github.com/gofiber/fiber"
	"google.golang.org/api/cloudfunctions/v1"
)

type AppConfig struct {
	projectID      string
	region         string
	cloudfunctions *cloudfunctions.Service
	fiber          *fiber.App
}

type CreateCFRequest struct {
	Name  string         `json:"name"`
	Files []FilesRequest `json:"files"`
}

type FilesRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type GetCFRequest struct {
	FunctionID string `json:"function_id"`
}

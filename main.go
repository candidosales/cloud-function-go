package main

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/logger"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/option"
)

const (
	projectID = "vendasta-hackathon"
	region    = "us-central1"
)

var (
	ErrFunctionIDIsEmpty        = errors.New("function_id is empty")
	ErrGoogleCredentialsIsEmpty = errors.New("GOOGLE_APPLICATION_CREDENTIALS variable is empty")
	ErrProjectIDIsInvalid = errors.New("projectID constant is invalid. Please set the projectID")
)

func main() {
	ctx := context.Background()

	credentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentials == "" {
		log.Fatal(ErrGoogleCredentialsIsEmpty)
	}

	if projectID == "YOUR-PROJECT" {
		log.Fatal(ErrProjectIDIsInvalid)
	}

	cfService, err := cloudfunctions.NewService(ctx, option.WithCredentialsFile(credentials))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cloud Function Service started ... ")

	appFiber := fiber.New()
	appFiber.Use(logger.New())
	appFiber.Use(cors.New())

	appConfig := AppConfig{
		projectID:      projectID,
		region:         region,
		cloudfunctions: cfService,
		fiber:          appFiber,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appConfig.setupRoutes()
	appConfig.fiber.Listen(port)
}

func (app *AppConfig) createCloudFunction(createCFRequest *CreateCFRequest) (string, error) {
	start := time.Now()

	cfName := fmt.Sprintf("%s", strings.ToLower(createCFRequest.Name))
	name := app.getLocation() + "/functions/" + cfName

	fmt.Println("🤖 Deployment initialized.")
	fmt.Println("Packaging sources...")

	sourceUploadUrl := app.generateUploadURL()
	bodyBuffer := app.packFiles(createCFRequest.Files)

	fmt.Println("Uploading source...")

	_, err := app.uploadFileRequest(sourceUploadUrl, bodyBuffer)
	if err != nil {
		fmt.Printf("uploadFileRequest - err[%#v] \n", err)
		return "", err
	}
	requestBody := app.buildRequestBody(sourceUploadUrl, name)

	fmt.Println("Deploying function...")

	if app.cloudFunctionExists(name) {
		operation, err := app.cloudfunctions.Projects.Locations.Functions.Patch(name, requestBody).Do()
		if err != nil {
			fmt.Printf("Patch - err[%#v] \n", err)
			return "", err
		}

		if operation != nil {
			fmt.Printf("🚀 Update - Function deployed in %d seconds. \n", time.Now().Sub(start))
			return "Updated the Cloud Function", nil
		}

	} else {
		operation, err := app.cloudfunctions.Projects.Locations.Functions.Create(app.getLocation(), requestBody).Do()
		if err != nil {
			fmt.Printf("Create - err[%#v] \n", err)
		}

		if operation != nil {
			fmt.Printf("🚀 Create - Function deployed in %d seconds. \n", time.Now().Sub(start))
			return "Created the Cloud Function", nil
		}
	}
	return "", nil
}

func (app *AppConfig) getLocation() string {
	return fmt.Sprintf("projects/%s/locations/%s", app.projectID, app.region)
}

func (app *AppConfig) generateUploadURL() string {
	response, err := app.cloudfunctions.Projects.Locations.Functions.GenerateUploadUrl(app.getLocation(), &cloudfunctions.GenerateUploadUrlRequest{}).Do()
	if err != nil {
		fmt.Printf("err: %#v \n", err)
		return ""
	}

	return response.UploadUrl
}

func (app *AppConfig) packFiles(filesRequest []FilesRequest) *bytes.Buffer {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)
	defer w.Close()

	for _, file := range filesRequest {
		f, err := w.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(file.Content))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("📦 Source code packaged.")
	return buf
}

func (app *AppConfig) uploadFileRequest(sourceUploadURL string, body *bytes.Buffer) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPut, sourceUploadURL, body)
	req.Header.Set("Content-Type", "application/zip")
	req.Header.Set("X-Goog-Content-Length-Range", "0,104857600") // https://cloud.google.com/storage/docs/json_api/v1/parameters#xgoogcontentlengthrange

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			return nil, err
		}

		err = resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 200 {
			fmt.Println("🛸 Source uploaded to cloud.")
		}
	}

	return resp, err
}

func (app *AppConfig) buildRequestBody(sourceUploadURL string, name string) *cloudfunctions.CloudFunction {
	return &cloudfunctions.CloudFunction{
		Name:              name,                           // Required - Name must start with a letter followed by up to 62 letters, numbers or hyphens, and cannot end with a hyphen
		AvailableMemoryMb: 128,                            // Optional
		EntryPoint:        "helloWorld",                   // Required
		Runtime:           "nodejs8",                      // Required
		HttpsTrigger:      &cloudfunctions.HttpsTrigger{}, // Required
		SourceUploadUrl:   sourceUploadURL,                // Required
		IngressSettings:   "ALLOW_ALL",                    // Optional
	}
}

func (app *AppConfig) cloudFunctionExists(name string) bool {
	_, err := app.cloudfunctions.Projects.Locations.Functions.Get(name).Do()
	if err != nil {
		return false
	}
	return true
}

func (app *AppConfig) getCloudFunction(functionID string) (*cloudfunctions.CloudFunction, error) {
	if functionID == "" {
		return nil, ErrFunctionIDIsEmpty
	}

	name := app.getLocation() + "/functions/" + functionID

	response, err := app.cloudfunctions.Projects.Locations.Functions.Get(name).Do()
	if err != nil {
		return nil, err
	}
	return response, err
}

func (app *AppConfig) setIAMPolicy(functionID string) (*cloudfunctions.Policy, error) {

	if functionID == "" {
		return nil, ErrFunctionIDIsEmpty
	}

	name := app.getLocation() + "/functions/" + functionID

	setIAMPolicyRequest := &cloudfunctions.SetIamPolicyRequest{
		Policy: &cloudfunctions.Policy{
			Bindings: []*cloudfunctions.Binding{
				{
					Role:    "roles/cloudfunctions.invoker",
					Members: []string{"allUsers"},
				},
			},
		},
	}

	response, err := app.cloudfunctions.Projects.Locations.Functions.SetIamPolicy(name, setIAMPolicyRequest).Do()
	if err != nil {
		return nil, err
	}
	return response, err
}

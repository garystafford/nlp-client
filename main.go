// author: Gary A. Stafford
// site: https://programmaticponderings.com
// license: MIT License
// purpose: Simple echo-based microservice: nlp-client

package main

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	serverPort = ":" + getEnv("NLP_CLIENT_PORT", "8080")
	urlRack    = getEnv("RAKE_ENDPOINT", "http://localhost:8080")
	urlProse   = getEnv("PROSE_ENDPOINT", "http://localhost:8080")
	urlLang    = getEnv("LANG_ENDPOINT", "http://localhost:8080")
	urlDynamo  = getEnv("DYNAMO_ENDPOINT", "http://localhost:8080")
	apiKey     = getEnv("API_KEY", "")
	log        = logrus.New()

	// Echo instance
	e = echo.New()
)

func init() {
	log.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-API-Key",
		Skipper: func(c echo.Context) bool {
			if strings.HasPrefix(c.Request().RequestURI, "/health") {
				return true
			}
			return false
		},
		Validator: func(key string, c echo.Context) (bool, error) {
			log.Debugf("API_KEY: %v", apiKey)
			return key == apiKey, nil
		},
	}))

	// Routes
	e.GET("/health", getHealth)
	e.GET("/error", getError)
	e.GET("/routes", getRoutes)
	e.POST("/keywords", getKeywords)
	e.POST("/tokens", getTokens)
	e.POST("/entities", getEntities)
	e.POST("/language", getLanguage)
	e.POST("/record", putDynamo)

	// Start server
	e.Logger.Fatal(e.Start(serverPort))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getRoutes(c echo.Context) error {
	response, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, response)
}

func getHealth(c echo.Context) error {
	var response interface{}
	err := json.Unmarshal([]byte(`{"status":"UP"}`), &response)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, response)
}

func getError(c echo.Context) error {
	var response interface{}
	err := json.Unmarshal([]byte(`{"error":"500 Internal Server Error"}`), &response)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusInternalServerError, response)
}

func getKeywords(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlRack+"/keywords", c.Request().Body)

	return serviceResponse(err, req, c)
}

func getTokens(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlProse+"/tokens", c.Request().Body)

	return serviceResponse(err, req, c)
}

func getEntities(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlProse+"/entities", c.Request().Body)

	return serviceResponse(err, req, c)
}

func getLanguage(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlLang+"/language", c.Request().Body)

	return serviceResponse(err, req, c)
}

func putDynamo(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlDynamo+"/record", c.Request().Body)

	return serviceResponse(err, req, c)
}

func serviceResponse(err error, req *http.Request, c echo.Context) error {
	req.Header.Set("Authorization", c.Request().Header.Get("Authorization"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSONBlob(http.StatusOK, body)
}

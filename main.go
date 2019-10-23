// author: Gary A. Stafford
// site: https://programmaticponderings.com
// license: MIT License
// purpose: Simple echo-based microservice: nlp-client

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/context"
)

var (
	portClient = getEnv("NLP_CLIENT_PORT", "8080")
	urlRack    = getEnv("RACK_PORT", "8080")
	urlProse   = getEnv("PROSE_PORT", "8080")

	// Echo instance
	e = echo.New()
)

func main() {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			if strings.HasPrefix(c.Request().RequestURI, "/health") {
				return true
			}
			return false
		},
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("AUTH_KEY"), nil
		},
	}))

	// Routes
	e.GET("/health", getHealth)
	e.GET("/routes", getRoutes)
	e.POST("/keywords", getKeywords)
	e.POST("/tokens", getTokens)
	e.POST("/entities", getEntities)

	// Start server
	e.Logger.Fatal(e.Start(portClient))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getRoutes(c echo.Context) error {
	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, data)
}

func getHealth(c echo.Context) error {
	var response interface{}
	err := json.Unmarshal([]byte(`{"status":"UP"}`), &response)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, response)
}

func getKeywords(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlRack+"/keywords", c.Request().Body)
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

func getTokens(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlProse+"/tokens", c.Request().Body)
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

func getEntities(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlProse+"/entities", c.Request().Body)
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

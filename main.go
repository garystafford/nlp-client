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
	portClient = os.Getenv("NLP_CLIENT_PORT")
	urlRack = os.Getenv("RACK_ENDPOINT")
	urlProse = os.Getenv("PROSE_ENDPOINT")
)

func main() {
	// Echo instance
	e := echo.New()

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
	e.POST("/keywords", extractKeywords)
	e.POST("/tokens", extractTokens)
	e.POST("/entities", extractEntities)

	// Start server
	e.Logger.Fatal(e.Start(portClient))
}

func getHealth(c echo.Context) error {
	var response interface{}
	err := json.Unmarshal([]byte(`{"status":"UP"}`), &response)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, response)
}

func extractKeywords(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlRack + "/keywords", c.Request().Body)
	req.Header.Set("Authorization", c.Request().Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSONBlob(http.StatusOK, body)
}

func extractTokens(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlProse + "/keywords", c.Request().Body)
	req.Header.Set("Authorization", c.Request().Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSONBlob(http.StatusOK, body)
}

func extractEntities(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlProse + "/entities", c.Request().Body)
	req.Header.Set("Authorization", c.Request().Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSONBlob(http.StatusOK, body)
}

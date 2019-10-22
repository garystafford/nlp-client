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

func main() {
	port := os.Getenv("NLP_CLIENT_PORT")

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
	e.Logger.Fatal(e.Start(port))
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
	url := os.Getenv("RACK_ENDPOINT")
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, c.Request().Body)
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
	jsonMap := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	text := jsonMap["text"] // request body - 'text' json attribute value
	tokens, err := getTokens(text.(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, tokens)
}

func extractEntities(c echo.Context) error {
	jsonMap := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	text := jsonMap["text"] // request body - 'text' json attribute value
	entities, err := getEntities(text.(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, entities)
}

// author: Gary A. Stafford
// site: https://programmaticponderings.com
// license: MIT License
// purpose: NLP microservices: nlp-client
// modified: 2021-06-15

package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/context"
)

var (
	logLevel   = getEnv("LOG_LEVEL", "1") // INFO
	serverPort = getEnv("NLP_CLIENT_PORT", ":8080")
	urlRake    = getEnv("RAKE_ENDPOINT", "http://localhost:8080")
	urlProse   = getEnv("PROSE_ENDPOINT", "http://localhost:8080")
	urlLang    = getEnv("LANG_ENDPOINT", "http://localhost:8080")
	urlDynamo  = getEnv("DYNAMO_ENDPOINT", "http://localhost:8080")
	apiKey     = getEnv("API_KEY", "ChangeMe")
	e          = echo.New()
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getRoutes(c echo.Context) error {
	//response, err := json.MarshalIndent(e.Routes(), "", "  ")
	//if err != nil {
	//	return echo.NewHTTPError(http.StatusInternalServerError, err)
	//}

	return c.JSON(http.StatusOK, e.Routes())
}

func getHealth(c echo.Context) error {
	type HealthStatus struct {
		Status string `json:"status,omitempty"`
	}
	healthStatus := HealthStatus{"Up"}
	//_, err := json.Marshal(healthStatus)
	//if err != nil {
	//	e.Logger.Error(err)
	//}
	return c.JSON(http.StatusOK, healthStatus)
}

func getHealthUpstream(c echo.Context) error {
	var urlHealth = ""
	switch c.Param("app") {
	case "rake":
		urlHealth = urlRake
	case "prose":
		urlHealth = urlProse
	case "lang":
		urlHealth = urlLang
	case "dynamo":
		urlHealth = urlDynamo
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlHealth+"/health", c.Request().Body)

	return serviceResponse(err, req, c)
}

func getError(c echo.Context) error {
	return echo.NewHTTPError(http.StatusInternalServerError)
}

func getKeywords(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlRake+"/keywords", c.Request().Body)

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

func getSentences(c echo.Context) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlProse+"/sentences", c.Request().Body)

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
	req.Header.Set("X-API-Key", c.Request().Header.Get("X-API-Key"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				e.Logger.Error(err)
			}
		}(resp.Body)
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

func run() error {
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
			e.Logger.Debugf("API_KEY: %v", apiKey)
			return key == apiKey, nil
		},
	}))

	// Routes
	e.GET("/health", getHealth)
	e.GET("/health/:app", getHealthUpstream)
	e.GET("/error", getError)
	e.GET("/routes", getRoutes)
	e.POST("/keywords", getKeywords)
	e.POST("/tokens", getTokens)
	e.POST("/entities", getEntities)
	e.POST("/sentences", getSentences)
	e.POST("/language", getLanguage)
	e.POST("/record", putDynamo)

	// Start server
	return e.Start(serverPort)
}

func init() {
	level, _ := strconv.Atoi(logLevel)
	e.Logger.SetLevel(log.Lvl(level))
}

func main() {
	if err := run(); err != nil {
		e.Logger.Fatal(err)
		os.Exit(1)
	}
}

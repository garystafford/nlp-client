package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

func TestGetHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	err := getHealth(c)
	if err != nil {
		t.Error(err)
	}
	res := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(res.Body)
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	type HealthStatus struct {
		Status string `json:"status"`
	}
	expectedStatus := HealthStatus{"Up"}
	var responseBody HealthStatus

	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		t.Error(err)
	}

	if assert.NoError(t, getHealth(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, &expectedStatus, &responseBody)
	}
}

func TestGetError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/routes", nil)
	w := httptest.NewRecorder()
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
	c := e.NewContext(req, w)
	err := getRoutes(c)
	if err != nil {
		t.Error(err)
	}
	res := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(res.Body)
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	type Route struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Name   string `json:"name"`
	}

	prefix := "github.com/garystafford/nlp-client"
	expectedStatus := []Route{
		{"GET", "/health", prefix + ".getHealth"},
		{"GET", "/health/:app", prefix + ".getHealthUpstream"},
		{"GET", "/error", prefix + ".getError"},
		{"GET", "/routes", prefix + ".getRoutes"},
		{"POST", "/keywords", prefix + ".getKeywords"},
		{"POST", "/tokens", prefix + ".getTokens"},
		{"POST", "/entities", prefix + ".getEntities"},
		{"POST", "/sentences", prefix + ".getSentences"},
		{"POST", "/language", prefix + ".getLanguage"},
		{"POST", "/record", prefix + ".putDynamo"},
	}
	var responseBody []Route

	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		t.Error(err)
	}

	// sort both arrays of Route structs so they are in identical order
	sort.Slice(expectedStatus, func(i, j int) bool { return expectedStatus[i].Path < expectedStatus[j].Path })
	sort.Slice(responseBody, func(i, j int) bool { return responseBody[i].Path < responseBody[j].Path })

	if assert.NoError(t, getRoutes(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, &expectedStatus, &responseBody)
	}
}

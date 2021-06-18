package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	API_KEY = "ChangeMe"
)

type Route struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Name   string `json:"name"`
}

type Routes []Route

func TestHealthUsingUnmarshal(t *testing.T) {
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
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	type HealthStatus struct {
		Status string `json:"status"`
	}
	expectedStatus := HealthStatus{"Up"}
	var responseBody HealthStatus

	err = json.Unmarshal([]byte(strings.Trim(string(data), "\n")), &responseBody)
	if err != nil {
		t.Error(err)
	}

	if assert.NoError(t, getHealth(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, &expectedStatus, &responseBody)
	}
}

func TestHealthUsingStrings(t *testing.T) {
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
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	expected := `{"status":"Up"}`
	actual := strings.Trim(string(data), "\n")

	if assert.NoError(t, getHealth(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, &expected, &actual)
	}
}

func TestGetError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/routes", nil)
	req.Header.Set("X-API-Key", API_KEY)
	req.Header.Set("Content-Type", "application/json")

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
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}


	//type Routes []struct {
	//	Method string `json:"method"`
	//	Path   string `json:"path"`
	//	Name   string `json:"name"`
	//}

	expectedStatus := Routes{{"GET", "/health", "main.getHealth"}}
	t.Log(expectedStatus[0].Name)
	var responseBody Routes

	//err = json.Unmarshal(data, &responseBody)
	jsonBody := []byte(strings.Trim(string(data), "\n"))
	t.Log(jsonBody)
	err = json.Unmarshal(jsonBody, &responseBody)
	if err != nil {
		t.Errorf("json.Unmarshal: %v", err)
	}

	if assert.NoError(t, getRoutes(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, &expectedStatus[0].Name, &responseBody)
	}
}
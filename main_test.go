package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	var err error

	database, err = openDatabase()
	if err != nil {
		log.Printf("failed to open database: %v", err)
		os.Exit(1)
	}
	defer database.Close()

	os.Exit(m.Run())
}

func TestSetByIDHandler(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet, "/api/sets/1", nil)
	recorder := httptest.NewRecorder()

	setByIDHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", recorder.Code)
	}

	var actualSet CardSet

	if err := json.NewDecoder(recorder.Body).Decode(&actualSet); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if actualSet.ID != 1 {
		t.Errorf("expected ID 1, got %d", actualSet.ID)
	}

	if actualSet.Title != "Thai Basics" {
		t.Errorf("expected title Thai Basics, got %s", actualSet.Title)
	}
}

func TestSetByIDHandlerNotFound(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/sets/999", nil)
	recorder := httptest.NewRecorder()

	setByIDHandler(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", recorder.Code)
	}
}

func TestSetsHandlerPost(t *testing.T) {

	body := strings.NewReader(`{"title":"Greetings"}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/sets",
		body,
	)

	recorder := httptest.NewRecorder()

	setsHandler(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", recorder.Code)
	}

	var createdSet CardSet

	if err := json.NewDecoder(recorder.Body).Decode(&createdSet); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdSet.Title != "Greetings" {
		t.Errorf("expected title Greetings, got %s", createdSet.Title)
	}

	if createdSet.ID <= 0 {
		t.Errorf("expected positive ID, got %d", createdSet.ID)
	}
}

func TestSetsHandlerPostEmptyTitle(t *testing.T) {
	body := strings.NewReader(`{"title":""}`)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/sets",
		body,
	)

	recorder := httptest.NewRecorder()

	setsHandler(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", recorder.Code)
	}
}

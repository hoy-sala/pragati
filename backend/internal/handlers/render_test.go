package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pragati/backend/internal/models"
)

func TestRenderJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	data := models.APIResponse{
		Data: map[string]string{"key": "value"},
	}
	renderJSON(rec, http.StatusOK, data)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}

	var resp models.APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected Data to be a map, got %T", resp.Data)
	}
	if dataMap["key"] != "value" {
		t.Errorf("expected Data.key='value', got '%v'", dataMap["key"])
	}
}

func TestRenderJSON_Error(t *testing.T) {
	rec := httptest.NewRecorder()
	errResp := models.APIResponse{
		Error: &models.APIError{
			Code:    "NOT_FOUND",
			Message: "resource not found",
		},
	}
	renderJSON(rec, http.StatusNotFound, errResp)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	var resp models.APIResponse
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp.Error == nil {
		t.Fatal("expected Error to be non-nil")
	}
	if resp.Error.Code != "NOT_FOUND" {
		t.Errorf("expected Error.Code 'NOT_FOUND', got '%s'", resp.Error.Code)
	}
	if resp.Error.Message != "resource not found" {
		t.Errorf("expected Error.Message 'resource not found', got '%s'", resp.Error.Message)
	}
}

func TestRenderJSON_WithPagination(t *testing.T) {
	rec := httptest.NewRecorder()
	data := models.APIResponse{
		Data: []string{"a", "b"},
		Meta: models.Pagination{Offset: 0, Limit: 10, Total: 2},
	}
	renderJSON(rec, http.StatusOK, data)

	var resp models.APIResponse
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp.Meta == nil {
		t.Fatal("expected Meta to be non-nil")
	}

	meta, ok := resp.Meta.(map[string]interface{})
	if !ok {
		t.Fatalf("expected Meta to be a map, got %T", resp.Meta)
	}
	if meta["total"] != float64(2) {
		t.Errorf("expected Meta.total=2, got %v", meta["total"])
	}
}

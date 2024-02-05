package request_dto_test

import (
	"bytes"
	"context"
	"eau-de-go/internal/transport/http/request_dto"
	"github.com/google/uuid"
	"net/http"
	"testing"
)

func TestMakeCreateAppUserParamsFromRequest_HappyPath(t *testing.T) {
	requestBody := `{
		"username": "testuser",
		"first_name": "Test",
		"last_name": "User",
		"email": "testuser@example.com",
		"password": "password123",
		"is_active": true
	}`
	req, _ := http.NewRequest("POST", "/appuser", bytes.NewBuffer([]byte(requestBody)))
	params, err := request_dto.MakeCreateAppUserParamsFromRequest(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if params.Username != "testuser" || params.FirstName != "Test" || params.LastName != "User" || params.Email != "testuser@example.com" || params.Password != "password123" {
		t.Fatalf("Expected params to match input, got %v", params)
	}
}

func TestMakeCreateAppUserParamsFromRequest_MissingStringFieldGetEmptyString(t *testing.T) {
	requestBody := `{
		"username": "testuser",
		"email": "testuser@example.com",
		"password": "password123",
		"is_active": true
	}`
	req, _ := http.NewRequest("POST", "/appuser", bytes.NewBuffer([]byte(requestBody)))
	params, err := request_dto.MakeCreateAppUserParamsFromRequest(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if params.Username != "testuser" || params.FirstName != "" || params.LastName != "" || params.Email != "testuser@example.com" || params.Password != "password123" {
		t.Fatalf("Expected params to match input, got %v", params)
	}
}

func TestMakeCreateAppUserParamsFromRequest_BadRequest(t *testing.T) {
	requestBody := `{
		"username": "testuser",
		"first_name": "Test",
		"last_name": "User",
		"email": "testuser@example.com",
		"password": "password123",
		"is_active": "true"
	}`
	req, _ := http.NewRequest("POST", "/appuser", bytes.NewBuffer([]byte(requestBody)))
	_, err := request_dto.MakeCreateAppUserParamsFromRequest(req)

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestMakeUpdateAppUserParamsFromRequest_HappyPath(t *testing.T) {
	requestBody := `{
		"first_name": "Updated",
		"last_name": "User"
	}`
	req, _ := http.NewRequest("PUT", "/appuser", bytes.NewBuffer([]byte(requestBody)))
	req = req.WithContext(context.WithValue(req.Context(), "jwt_claims", map[string]interface{}{"id": uuid.New().String()}))
	params, err := request_dto.MakeUpdateAppUserParamsFromRequest(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if params.FirstName.String != "Updated" || params.LastName.String != "User" {
		t.Fatalf("Expected params to match input, got %v", params)
	}
}

func TestMakeUpdateAppUserParamsFromRequest_MissingStringFieldGetNil(t *testing.T) {
	requestBody := `{
		"first_name": "Updated"
	}`
	req, _ := http.NewRequest("PUT", "/appuser", bytes.NewBuffer([]byte(requestBody)))
	req = req.WithContext(context.WithValue(req.Context(), "jwt_claims", map[string]interface{}{"id": uuid.New().String()}))
	params, err := request_dto.MakeUpdateAppUserParamsFromRequest(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if params.FirstName.String != "Updated" || params.LastName.Valid != false || params.LastName.String != "" {
		t.Fatalf("Expected params to match input, got %v", params)
	}
}

func TestMakeUpdateAppUserParamsFromRequest_BadRequest(t *testing.T) {
	requestBody := `{
		"first_name": 123,
		"last_name": "User"
	}`
	req, _ := http.NewRequest("PUT", "/appuser", bytes.NewBuffer([]byte(requestBody)))
	req = req.WithContext(context.WithValue(req.Context(), "jwt_claims", map[string]interface{}{"id": uuid.New().String()}))
	_, err := request_dto.MakeUpdateAppUserParamsFromRequest(req)

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestMakeUpdateAppUserParamsFromRequest_NoJwtClaims(t *testing.T) {
	requestBody := `{
		"first_name": "Updated",
		"last_name": "User"
	}`
	req, _ := http.NewRequest("PUT", "/appuser", bytes.NewBuffer([]byte(requestBody)))
	_, err := request_dto.MakeUpdateAppUserParamsFromRequest(req)

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestMakeUpdateAppUserParamsFromRequest_InvalidJwtClaims(t *testing.T) {
	requestBody := `{
		"first_name": "Updated",
		"last_name": "User"
	}`
	req, _ := http.NewRequest("PUT", "/appuser", bytes.NewBuffer([]byte(requestBody)))
	req = req.WithContext(context.WithValue(req.Context(), "jwt_claims", map[string]interface{}{"id": "invalid_uuid"}))
	_, err := request_dto.MakeUpdateAppUserParamsFromRequest(req)

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

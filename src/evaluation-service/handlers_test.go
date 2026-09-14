package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandlerReturnsOK(t *testing.T) {
	a := &App{}
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	a.healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebeu %d", rec.Code)
	}
}

func TestEvaluationHandlerMissingParamsReturns400(t *testing.T) {
	a := &App{}
	req := httptest.NewRequest(http.MethodGet, "/evaluate", nil)
	rec := httptest.NewRecorder()

	a.evaluationHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava status 400 sem user_id/flag_name, recebeu %d", rec.Code)
	}
}

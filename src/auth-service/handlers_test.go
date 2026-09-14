package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandlerReturnsOK(t *testing.T) {
	app := &App{}
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	app.healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebeu %d", rec.Code)
	}
	if got := rec.Body.String(); got != `{"status":"ok"}`+"\n" {
		t.Errorf("corpo inesperado: %q", got)
	}
}

func TestMasterKeyAuthMiddlewareRejectsWrongKey(t *testing.T) {
	app := &App{MasterKey: "segredo-correto"}
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/keys", nil)
	req.Header.Set("Authorization", "Bearer chave-errada")
	rec := httptest.NewRecorder()

	app.masterKeyAuthMiddleware(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("esperava status 403, recebeu %d", rec.Code)
	}
	if called {
		t.Error("o handler protegido nao deveria ter sido chamado")
	}
}

func TestMasterKeyAuthMiddlewareAllowsCorrectKey(t *testing.T) {
	app := &App{MasterKey: "segredo-correto"}
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/keys", nil)
	req.Header.Set("Authorization", "Bearer segredo-correto")
	rec := httptest.NewRecorder()

	app.masterKeyAuthMiddleware(next).ServeHTTP(rec, req)

	if !called {
		t.Error("o handler protegido deveria ter sido chamado")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("esperava status 200, recebeu %d", rec.Code)
	}
}

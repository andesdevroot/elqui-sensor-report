package copernicus

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestClient_Authenticate_OK verifica el flujo feliz: el endpoint responde 200
// con access_token y el cliente lo devuelve, enviando el grant correcto.
func TestClient_Authenticate_OK(t *testing.T) {
	type captured struct {
		method      string
		contentType string
		form        url.Values
	}
	reqCh := make(chan captured, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		reqCh <- captured{
			method:      r.Method,
			contentType: r.Header.Get("Content-Type"),
			form:        r.PostForm,
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"token-123","token_type":"Bearer","expires_in":600}`))
	}))
	defer srv.Close()

	c := &Client{TokenURL: srv.URL, ClientID: "mi-id", ClientSecret: "mi-secreto"}

	token, err := c.Authenticate(context.Background())
	if err != nil {
		t.Fatalf("Authenticate devolvió error: %v", err)
	}
	if token != "token-123" {
		t.Errorf("token = %q, se esperaba %q", token, "token-123")
	}

	got := <-reqCh
	if got.method != http.MethodPost {
		t.Errorf("método = %q, se esperaba POST", got.method)
	}
	if got.contentType != "application/x-www-form-urlencoded" {
		t.Errorf("Content-Type = %q, se esperaba application/x-www-form-urlencoded", got.contentType)
	}
	if got.form.Get("grant_type") != "client_credentials" {
		t.Errorf("grant_type = %q, se esperaba client_credentials", got.form.Get("grant_type"))
	}
	if got.form.Get("client_id") != "mi-id" {
		t.Errorf("client_id = %q, se esperaba mi-id", got.form.Get("client_id"))
	}
	if got.form.Get("client_secret") != "mi-secreto" {
		t.Errorf("client_secret = %q, se esperaba mi-secreto", got.form.Get("client_secret"))
	}
}

// TestClient_Authenticate_Unauthorized verifica el rechazo de credenciales (401).
func TestClient_Authenticate_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid_client"}`))
	}))
	defer srv.Close()

	c := &Client{TokenURL: srv.URL, ClientID: "id", ClientSecret: "incorrecto"}

	_, err := c.Authenticate(context.Background())
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("Se esperaba ErrUnauthorized, se obtuvo: %v", err)
	}
}

// TestClient_Authenticate_SinCredenciales verifica que falle sin tocar la red
// cuando no hay credenciales configuradas.
func TestClient_Authenticate_SinCredenciales(t *testing.T) {
	c := &Client{}

	_, err := c.Authenticate(context.Background())
	if err == nil {
		t.Fatal("Se esperaba error por credenciales faltantes, se obtuvo nil")
	}
}

// TestNewClient_DesdeEntorno verifica que NewClient tome las credenciales de las
// variables de entorno y aplique los valores por defecto.
func TestNewClient_DesdeEntorno(t *testing.T) {
	t.Setenv(EnvClientID, "id-env")
	t.Setenv(EnvClientSecret, "secreto-env")

	c := NewClient()
	if c.ClientID != "id-env" {
		t.Errorf("ClientID = %q, se esperaba %q", c.ClientID, "id-env")
	}
	if c.ClientSecret != "secreto-env" {
		t.Errorf("ClientSecret = %q, se esperaba %q", c.ClientSecret, "secreto-env")
	}
	if c.TokenURL != DefaultTokenURL {
		t.Errorf("TokenURL = %q, se esperaba el endpoint por defecto", c.TokenURL)
	}
	if c.HTTPClient == nil || c.HTTPClient.Timeout <= 0 {
		t.Error("HTTPClient debe existir y tener un timeout por defecto")
	}
}

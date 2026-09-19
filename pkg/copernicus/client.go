// Package copernicus implementa el cliente HTTP de Copernicus Data Space
// (Sentinel-2 L2A): autenticación OAuth2 y, más adelante, búsqueda por AOI y
// descarga de bandas.
package copernicus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// DefaultTokenURL es el endpoint OAuth2 de Copernicus Data Space (realm CDSE).
const DefaultTokenURL = "https://identity.dataspace.copernicus.eu/auth/realms/CDSE/protocol/openid-connect/token"

// Variables de entorno con las credenciales de Copernicus Data Space.
const (
	EnvClientID     = "CDSE_CLIENT_ID"
	EnvClientSecret = "CDSE_CLIENT_SECRET"
)

// DefaultTimeout es el timeout por defecto de las llamadas HTTP.
const DefaultTimeout = 30 * time.Second

// ErrUnauthorized se devuelve cuando el token endpoint rechaza las credenciales.
var ErrUnauthorized = errors.New("copernicus: credenciales rechazadas (401)")

// Client es el cliente HTTP de Copernicus Data Space.
type Client struct {
	// TokenURL es el endpoint OAuth2; si está vacío se usa DefaultTokenURL.
	TokenURL string
	// ClientID y ClientSecret son las credenciales OAuth2. NewClient las lee de
	// las variables de entorno CDSE_CLIENT_ID y CDSE_CLIENT_SECRET.
	ClientID     string
	ClientSecret string
	// HTTPClient permite configurar timeout y transporte; si es nil se usa uno
	// con DefaultTimeout. Las credenciales nunca se escriben en el código.
	HTTPClient *http.Client
}

// NewClient construye un Client con el endpoint y el timeout por defecto,
// tomando las credenciales de CDSE_CLIENT_ID y CDSE_CLIENT_SECRET.
func NewClient() *Client {
	return &Client{
		TokenURL:     DefaultTokenURL,
		ClientID:     os.Getenv(EnvClientID),
		ClientSecret: os.Getenv(EnvClientSecret),
		HTTPClient:   &http.Client{Timeout: DefaultTimeout},
	}
}

// Authenticate obtiene un access token con el grant OAuth2 client_credentials.
// Las credenciales viajan en el cuerpo (x-www-form-urlencoded), nunca en la URL.
func (c *Client) Authenticate(ctx context.Context) (string, error) {
	if c.ClientID == "" || c.ClientSecret == "" {
		return "", fmt.Errorf("autenticar: faltan credenciales (%s y %s)", EnvClientID, EnvClientSecret)
	}

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenEndpoint(), strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("construir request de token: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return "", fmt.Errorf("solicitar token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("%w: %s", ErrUnauthorized, responseExcerpt(resp))
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint devolvió %s: %s", resp.Status, responseExcerpt(resp))
	}

	var payload struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decodificar respuesta de token: %w", err)
	}
	if payload.AccessToken == "" {
		return "", errors.New("respuesta de token sin access_token")
	}

	return payload.AccessToken, nil
}

// tokenEndpoint devuelve el endpoint configurado o, si está vacío, el default.
func (c *Client) tokenEndpoint() string {
	if c.TokenURL != "" {
		return c.TokenURL
	}
	return DefaultTokenURL
}

// httpClient devuelve el cliente HTTP configurado o uno con el timeout por defecto.
func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: DefaultTimeout}
}

// responseExcerpt devuelve un extracto del cuerpo para incluirlo en el error.
func responseExcerpt(resp *http.Response) string {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 256))
	if err != nil {
		return "(cuerpo ilegible)"
	}
	return strings.TrimSpace(string(body))
}

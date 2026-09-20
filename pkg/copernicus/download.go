package copernicus

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DownloadTimeout es el timeout de la descarga de una banda (archivos ~100 MB).
const DownloadTimeout = 300 * time.Second

// Errores de la descarga de bandas.
var (
	// ErrAssetNotFound se devuelve cuando el item no contiene el asset pedido.
	ErrAssetNotFound = errors.New("asset no encontrado")
	// ErrNoHTTPSURL se devuelve cuando el asset solo ofrece esquemas no HTTP (p. ej. s3://).
	ErrNoHTTPSURL = errors.New("sin URL HTTPS para descargar")
)

// DownloadBand descarga la banda bandKey del item en destDir y devuelve la ruta
// del archivo escrito. Prefiere la alternativa HTTPS del asset y, si no existe,
// usa el href principal cuando ya es http(s). Reintenta con backoff exponencial
// ante respuestas 5xx y escribe por streaming (archivo .part + rename).
func (c *Client) DownloadBand(ctx context.Context, token string, item STACItem, bandKey, destDir string) (string, error) {
	asset, ok := item.Assets[bandKey]
	if !ok {
		return "", fmt.Errorf("%w: el item %s no tiene %q (disponibles: %s)", ErrAssetNotFound, item.ID, bandKey, strings.Join(assetKeys(item.Assets), ", "))
	}

	url, err := asset.downloadURL()
	if err != nil {
		return "", fmt.Errorf("banda %s del item %s: %w", bandKey, item.ID, err)
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", fmt.Errorf("crear directorio %s: %w", destDir, err)
	}
	dest := filepath.Join(destDir, item.ID+"_"+bandKey+path.Ext(asset.Href))

	for attempt := 0; ; attempt++ {
		status, statusText, err := c.downloadToFile(ctx, token, url, dest)
		switch {
		case err != nil:
			return "", err
		case status == http.StatusOK:
			return dest, nil
		case status == http.StatusUnauthorized:
			return "", fmt.Errorf("%w: descarga de %s rechazada", ErrUnauthorized, bandKey)
		case status >= 500 && attempt < c.MaxRetries:
			if err := sleep(ctx, c.RetryBaseDelay<<uint(attempt)); err != nil {
				return "", fmt.Errorf("esperar reintento de descarga: %w", err)
			}
		default:
			return "", fmt.Errorf("descargar %s: el servidor devolvió %s", bandKey, statusText)
		}
	}
}

// downloadURL elige la URL de descarga: la alternativa HTTPS si existe, o el
// href principal cuando ya es http(s). Devuelve ErrNoHTTPSURL si el asset solo
// ofrece esquemas no HTTP (p. ej. s3://).
func (a Asset) downloadURL() (string, error) {
	if alt, ok := a.Alternate["https"]; ok && alt.Href != "" {
		return alt.Href, nil
	}
	if strings.HasPrefix(a.Href, "http://") || strings.HasPrefix(a.Href, "https://") {
		return a.Href, nil
	}
	return "", fmt.Errorf("%w: el asset solo ofrece %q", ErrNoHTTPSURL, a.Href)
}

// downloadToFile hace un GET y escribe el cuerpo en dest (vía archivo .part).
// Devuelve el status HTTP para que el llamador decida los reintentos.
func (c *Client) downloadToFile(ctx context.Context, token, url, dest string) (int, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, "", fmt.Errorf("construir request de descarga: %w", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.downloadClient().Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("descargar %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 512))
		return resp.StatusCode, resp.Status, nil
	}

	tmp := dest + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return 0, "", fmt.Errorf("crear archivo temporal %s: %w", tmp, err)
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(tmp)
		return 0, "", fmt.Errorf("escribir %s: %w", tmp, err)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return 0, "", fmt.Errorf("cerrar %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, dest); err != nil {
		os.Remove(tmp)
		return 0, "", fmt.Errorf("renombrar %s → %s: %w", tmp, dest, err)
	}

	return resp.StatusCode, resp.Status, nil
}

// downloadClient devuelve un cliente HTTP con el timeout de descarga. Si el
// cliente del catálogo trae un Transport propio, se reutiliza.
func (c *Client) downloadClient() *http.Client {
	hc := &http.Client{Timeout: DownloadTimeout}
	if c.HTTPClient != nil {
		hc.Transport = c.HTTPClient.Transport
	}
	return hc
}

// assetKeys devuelve las claves de assets ordenadas, para mensajes de error.
func assetKeys(assets map[string]Asset) []string {
	keys := make([]string, 0, len(assets))
	for k := range assets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

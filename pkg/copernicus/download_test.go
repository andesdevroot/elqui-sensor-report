package copernicus

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// itemConAlternate arma un STACItem con un asset que ofrece alternativa HTTPS.
func itemConAlternate(id, bandKey, hrefS3, hrefHTTPS string) STACItem {
	return STACItem{
		ID: id,
		Assets: map[string]Asset{
			bandKey: {
				Href: hrefS3,
				Alternate: map[string]struct{ Href string }{
					"https": {Href: hrefHTTPS},
				},
			},
		},
	}
}

func TestClient_DownloadBand_OK(t *testing.T) {
	content := []byte("jp2-falso")
	authCh := make(chan string, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCh <- r.Header.Get("Authorization")
		_, _ = w.Write(content)
	}))
	defer srv.Close()

	item := itemConAlternate("S2B_TEST_001", "B04_10m", "s3://eodata/x/B04_10m.jp2", srv.URL+"/B04_10m.jp2")
	c := &Client{}
	dir := t.TempDir()

	dest, err := c.DownloadBand(context.Background(), "token-abc", item, "B04_10m", dir)
	if err != nil {
		t.Fatalf("DownloadBand devolvió error: %v", err)
	}
	if want := filepath.Join(dir, "S2B_TEST_001_B04_10m.jp2"); dest != want {
		t.Errorf("destino = %q, se esperaba %q", dest, want)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("leer archivo: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("contenido = %q, se esperaba %q", got, content)
	}
	if _, err := os.Stat(dest + ".part"); !os.IsNotExist(err) {
		t.Error("no debe quedar archivo .part tras una descarga exitosa")
	}
	if auth := <-authCh; auth != "Bearer token-abc" {
		t.Errorf("Authorization = %q, se esperaba Bearer token-abc", auth)
	}
}

func TestClient_DownloadBand_FallbackHref(t *testing.T) {
	// Sin alternativas: usa el href principal cuando ya es http(s).
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	item := STACItem{
		ID: "S2B_TEST_002",
		Assets: map[string]Asset{
			"B08_10m": {Href: srv.URL + "/B08_10m.jp2"},
		},
	}

	c := &Client{}
	dest, err := c.DownloadBand(context.Background(), "token-abc", item, "B08_10m", t.TempDir())
	if err != nil {
		t.Fatalf("DownloadBand devolvió error: %v", err)
	}
	if got, _ := os.ReadFile(dest); string(got) != "ok" {
		t.Errorf("contenido = %q, se esperaba ok", got)
	}
}

func TestClient_DownloadBand_SoloS3(t *testing.T) {
	item := STACItem{
		ID: "S2B_TEST_003",
		Assets: map[string]Asset{
			"B04_10m": {Href: "s3://eodata/Sentinel-2/B04_10m.jp2"},
		},
	}

	c := &Client{}
	_, err := c.DownloadBand(context.Background(), "token-abc", item, "B04_10m", t.TempDir())
	if !errors.Is(err, ErrNoHTTPSURL) {
		t.Fatalf("Se esperaba ErrNoHTTPSURL, se obtuvo: %v", err)
	}
	if !strings.Contains(err.Error(), "s3://") {
		t.Errorf("el error debe mencionar s3://, se obtuvo: %v", err)
	}
}

func TestClient_DownloadBand_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	item := itemConAlternate("S2B_TEST_004", "B04_10m", "s3://eodata/x/B04_10m.jp2", srv.URL+"/B04_10m.jp2")
	c := &Client{}

	_, err := c.DownloadBand(context.Background(), "token-malo", item, "B04_10m", t.TempDir())
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("Se esperaba ErrUnauthorized, se obtuvo: %v", err)
	}
}

func TestClient_DownloadBand_AssetInexistente(t *testing.T) {
	item := STACItem{ID: "S2B_TEST_005", Assets: map[string]Asset{"B04_10m": {Href: "https://x/B04_10m.jp2"}}}
	c := &Client{}

	_, err := c.DownloadBand(context.Background(), "token-abc", item, "SCL_20m", t.TempDir())
	if !errors.Is(err, ErrAssetNotFound) {
		t.Fatalf("Se esperaba ErrAssetNotFound, se obtuvo: %v", err)
	}
	if !strings.Contains(err.Error(), "B04_10m") {
		t.Errorf("el error debe listar los assets disponibles, se obtuvo: %v", err)
	}
}

func TestClient_DownloadBand_ReintentaEn5xx(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) < 3 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		_, _ = w.Write([]byte("contenido-final"))
	}))
	defer srv.Close()

	item := itemConAlternate("S2B_TEST_006", "B04_10m", "s3://eodata/x/B04_10m.jp2", srv.URL+"/B04_10m.jp2")
	c := &Client{MaxRetries: 3, RetryBaseDelay: time.Millisecond}

	dest, err := c.DownloadBand(context.Background(), "token-abc", item, "B04_10m", t.TempDir())
	if err != nil {
		t.Fatalf("DownloadBand devolvió error tras reintentar: %v", err)
	}
	if got, _ := os.ReadFile(dest); string(got) != "contenido-final" {
		t.Errorf("contenido = %q, se esperaba contenido-final", got)
	}
	if got := attempts.Load(); got != 3 {
		t.Errorf("intentos = %d, se esperaban 3 (2 fallos 5xx + 1 éxito)", got)
	}
}

func TestClient_DownloadBand_ReintentosAgotados(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	item := itemConAlternate("S2B_TEST_007", "B04_10m", "s3://eodata/x/B04_10m.jp2", srv.URL+"/B04_10m.jp2")
	c := &Client{MaxRetries: 1, RetryBaseDelay: time.Millisecond}

	_, err := c.DownloadBand(context.Background(), "token-abc", item, "B04_10m", t.TempDir())
	if err == nil {
		t.Fatal("Se esperaba error al agotar los reintentos, se obtuvo nil")
	}
	if got := attempts.Load(); got != 2 {
		t.Errorf("intentos = %d, se esperaban 2 (1 + 1 reintento)", got)
	}
}

package copernicus

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// featureCollectionFixture simula la respuesta del catálogo STAC con 3 escenas:
// dos bajo el umbral de nubosidad y una sobre él (85%), que debe filtrarse en Go.
const featureCollectionFixture = `{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "id": "S2A_MSIL2A_20260115_T19JCG_001",
      "bbox": [-71.30, -29.95, -71.20, -29.85],
      "properties": {"datetime": "2026-01-15T14:30:00Z", "eo:cloud_cover": 12.5},
      "assets": {
        "B04": {"href": "https://example.com/B04_001.tif"},
        "B08": {"href": "https://example.com/B08_001.tif"}
      }
    },
    {
      "type": "Feature",
      "id": "S2B_MSIL2A_20260120_T19JCG_002",
      "bbox": [-71.31, -29.96, -71.21, -29.86],
      "properties": {"datetime": "2026-01-20T14:30:00Z", "eo:cloud_cover": 44.0},
      "assets": {"B04": {"href": "https://example.com/B04_002.tif"}}
    },
    {
      "type": "Feature",
      "id": "S2A_MSIL2A_20260125_T19JCG_003",
      "bbox": [-71.32, -29.97, -71.22, -29.87],
      "properties": {"datetime": "2026-01-25T14:30:00Z", "eo:cloud_cover": 85.0},
      "assets": {"B04": {"href": "https://example.com/B04_003.tif"}}
    }
  ]
}`

func TestClient_SearchByAOI_OK(t *testing.T) {
	type captured struct {
		method string
		path   string
		auth   string
		ct     string
		raw    string
		body   map[string]any
	}
	reqCh := make(chan captured, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		reqCh <- captured{
			method: r.Method,
			path:   r.URL.Path,
			auth:   r.Header.Get("Authorization"),
			ct:     r.Header.Get("Content-Type"),
			raw:    string(raw),
			body:   body,
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(featureCollectionFixture))
	}))
	defer srv.Close()

	c := &Client{SearchURL: srv.URL + "/v1/search"}
	start := time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC)

	items, err := c.SearchByAOI(context.Background(), "token-abc", "POINT (-71.25 -29.90)", start, end, 50)
	if err != nil {
		t.Fatalf("SearchByAOI devolvió error: %v", err)
	}

	// Post-filtro en Go: la escena con 85% de nubosidad debe quedar fuera.
	if len(items) != 2 {
		t.Fatalf("Se esperaban 2 items tras el filtro, se obtuvieron %d", len(items))
	}
	if items[0].ID != "S2A_MSIL2A_20260115_T19JCG_001" {
		t.Errorf("ID = %q", items[0].ID)
	}
	if items[0].CloudCover != 12.5 {
		t.Errorf("CloudCover = %v, se esperaba 12.5", items[0].CloudCover)
	}
	if want := time.Date(2026, time.January, 15, 14, 30, 0, 0, time.UTC); !items[0].Datetime.Equal(want) {
		t.Errorf("Datetime = %v, se esperaba %v", items[0].Datetime, want)
	}
	if items[0].BBox != [4]float64{-71.30, -29.95, -71.20, -29.85} {
		t.Errorf("BBox = %v", items[0].BBox)
	}
	if got := items[0].Assets["B04"].Href; got != "https://example.com/B04_001.tif" {
		t.Errorf("Assets[B04].Href = %q", got)
	}

	got := <-reqCh
	if got.method != http.MethodPost {
		t.Errorf("método = %q, se esperaba POST", got.method)
	}
	if got.path != "/v1/search" {
		t.Errorf("path = %q, se esperaba /v1/search", got.path)
	}
	if got.auth != "Bearer token-abc" {
		t.Errorf("Authorization = %q, se esperaba Bearer token-abc", got.auth)
	}
	if got.ct != "application/json" {
		t.Errorf("Content-Type = %q, se esperaba application/json", got.ct)
	}
	if strings.Contains(got.raw, "eo:cloud_cover") {
		t.Error("el filtro de nubosidad no debe viajar en el body STAC")
	}
	cols, _ := got.body["collections"].([]any)
	if len(cols) != 1 || cols[0] != "sentinel-2-l2a" {
		t.Errorf("collections = %v, se esperaba [sentinel-2-l2a]", got.body["collections"])
	}
	if got.body["limit"] != float64(100) {
		t.Errorf("limit = %v, se esperaba 100", got.body["limit"])
	}
	if got.body["datetime"] != "2026-01-15T00:00:00Z/2026-01-31T00:00:00Z" {
		t.Errorf("datetime = %v", got.body["datetime"])
	}
	intersects, _ := got.body["intersects"].(map[string]any)
	if intersects["type"] != "Point" {
		t.Errorf("intersects.type = %v, se esperaba Point", intersects["type"])
	}
	coords, _ := intersects["coordinates"].([]any)
	if len(coords) != 2 || coords[0] != -71.25 || coords[1] != -29.9 {
		t.Errorf("intersects.coordinates = %v", intersects["coordinates"])
	}
}

func TestClient_SearchByAOI_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"detail":"token inválido"}`))
	}))
	defer srv.Close()

	c := &Client{SearchURL: srv.URL}
	start := time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)

	_, err := c.SearchByAOI(context.Background(), "token-malo", "POINT (-71.25 -29.90)", start, start.AddDate(0, 1, 0), 50)
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("Se esperaba ErrUnauthorized, se obtuvo: %v", err)
	}
}

func TestClient_SearchByAOI_ReintentaEn5xx(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) < 3 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"type":"FeatureCollection","features":[]}`))
	}))
	defer srv.Close()

	c := &Client{SearchURL: srv.URL, MaxRetries: 3, RetryBaseDelay: time.Millisecond}
	start := time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)

	items, err := c.SearchByAOI(context.Background(), "token-abc", "POINT (-71.25 -29.90)", start, start.AddDate(0, 1, 0), 50)
	if err != nil {
		t.Fatalf("SearchByAOI devolvió error tras reintentar: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("Se esperaban 0 items, se obtuvieron %d", len(items))
	}
	if got := attempts.Load(); got != 3 {
		t.Errorf("intentos = %d, se esperaban 3 (2 fallos 5xx + 1 éxito)", got)
	}
}

func TestClient_SearchByAOI_ReintentosAgotados(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := &Client{SearchURL: srv.URL, MaxRetries: 1, RetryBaseDelay: time.Millisecond}
	start := time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)

	_, err := c.SearchByAOI(context.Background(), "token-abc", "POINT (-71.25 -29.90)", start, start.AddDate(0, 1, 0), 50)
	if err == nil {
		t.Fatal("Se esperaba error al agotar los reintentos, se obtuvo nil")
	}
	if got := attempts.Load(); got != 2 {
		t.Errorf("intentos = %d, se esperaban 2 (1 + 1 reintento)", got)
	}
}

func TestClient_SearchByAOI_TokenVacio(t *testing.T) {
	c := &Client{}
	start := time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)

	_, err := c.SearchByAOI(context.Background(), "", "POINT (-71.25 -29.90)", start, start.AddDate(0, 1, 0), 50)
	if err == nil {
		t.Fatal("Se esperaba error con token vacío, se obtuvo nil")
	}
}

func TestParseWKT(t *testing.T) {
	t.Run("POINT", func(t *testing.T) {
		geom, err := parseWKT("POINT (-71.25 -29.90)")
		if err != nil {
			t.Fatalf("parseWKT devolvió error: %v", err)
		}
		if geom.Type != "Point" {
			t.Errorf("Type = %q, se esperaba Point", geom.Type)
		}
		coords, ok := geom.Coordinates.([2]float64)
		if !ok || coords != [2]float64{-71.25, -29.9} {
			t.Errorf("Coordinates = %v", geom.Coordinates)
		}
	})

	t.Run("POLYGON", func(t *testing.T) {
		wkt := "POLYGON ((-71.30 -29.95, -71.20 -29.95, -71.20 -29.85, -71.30 -29.85, -71.30 -29.95))"
		geom, err := parseWKT(wkt)
		if err != nil {
			t.Fatalf("parseWKT devolvió error: %v", err)
		}
		if geom.Type != "Polygon" {
			t.Errorf("Type = %q, se esperaba Polygon", geom.Type)
		}
		rings, ok := geom.Coordinates.([][][2]float64)
		if !ok || len(rings) != 1 || len(rings[0]) != 5 {
			t.Errorf("Coordinates = %#v", geom.Coordinates)
		}
	})

	t.Run("geometría no soportada", func(t *testing.T) {
		if _, err := parseWKT("CIRCLE (1 2)"); err == nil {
			t.Fatal("Se esperaba error para una geometría no soportada, se obtuvo nil")
		}
	})
}

package copernicus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DefaultSearchURL es el endpoint del catálogo STAC de Copernicus Data Space.
const DefaultSearchURL = "https://stac.dataspace.copernicus.eu/v1/search"

// DefaultMaxRetries y DefaultRetryBaseDelay configuran el backoff exponencial
// ante respuestas 5xx del catálogo STAC.
const (
	DefaultMaxRetries     = 3
	DefaultRetryBaseDelay = 500 * time.Millisecond
)

// ErrInvalidAOI se devuelve cuando el AOI no es un WKT soportado (POINT o POLYGON).
var ErrInvalidAOI = errors.New("AOI inválido")

// STACItem es una escena del catálogo STAC de Copernicus Data Space.
type STACItem struct {
	ID         string           `json:"id"`
	Datetime   time.Time        `json:"datetime"`
	CloudCover float64          `json:"cloud_cover"`
	BBox       [4]float64       `json:"bbox"`
	Assets     map[string]Asset `json:"assets"`
}

// Asset es un recurso descargable de una escena (p. ej. la banda B04).
type Asset struct {
	Href string `json:"href"`
	// Alternate contiene ubicaciones alternativas del asset (p. ej. "https").
	Alternate map[string]struct{ Href string } `json:"alternate"`
	// FileSize es el tamaño del recurso en bytes, cuando el catálogo lo informa.
	FileSize int64 `json:"file:size"`
	// Type es el tipo MIME del recurso (p. ej. image/jp2).
	Type string `json:"type"`
}

// searchRequest es el cuerpo de la consulta STAC.
type searchRequest struct {
	Collections []string  `json:"collections"`
	Intersects  *geometry `json:"intersects"`
	Datetime    string    `json:"datetime"`
	Limit       int       `json:"limit"`
}

// geometry es la geometría GeoJSON mínima que acepta el campo intersects.
type geometry struct {
	Type        string `json:"type"`
	Coordinates any    `json:"coordinates"`
}

// featureCollection es la respuesta del catálogo STAC.
type featureCollection struct {
	Features []rawItem `json:"features"`
}

// rawItem refleja la forma anidada de un Item STAC antes de aplanarlo a STACItem.
type rawItem struct {
	ID         string     `json:"id"`
	BBox       [4]float64 `json:"bbox"`
	Properties struct {
		Datetime   time.Time `json:"datetime"`
		CloudCover float64   `json:"eo:cloud_cover"`
	} `json:"properties"`
	Assets map[string]Asset `json:"assets"`
}

// SearchByAOI busca escenas Sentinel-2 L2A que intersectan el AOI (WKT: POINT o
// POLYGON) en el rango [startDate, endDate]. El filtro por nubosidad se aplica
// después, en Go, para mantener la consulta STAC simple y debuggeable.
func (c *Client) SearchByAOI(ctx context.Context, token, aoiWKT string, startDate, endDate time.Time, maxCloudCover float64) ([]STACItem, error) {
	if token == "" {
		return nil, errors.New("búsqueda STAC: token vacío")
	}
	geom, err := parseWKT(aoiWKT)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(searchRequest{
		Collections: []string{"sentinel-2-l2a"},
		Intersects:  geom,
		Datetime:    startDate.UTC().Format(time.RFC3339) + "/" + endDate.UTC().Format(time.RFC3339),
		Limit:       100,
	})
	if err != nil {
		return nil, fmt.Errorf("codificar consulta STAC: %w", err)
	}

	resp, err := c.doSearch(ctx, token, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("%w: %s", ErrUnauthorized, responseExcerpt(resp))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("búsqueda STAC devolvió %s: %s", resp.Status, responseExcerpt(resp))
	}

	var fc featureCollection
	if err := json.NewDecoder(resp.Body).Decode(&fc); err != nil {
		return nil, fmt.Errorf("decodificar FeatureCollection: %w", err)
	}

	items := make([]STACItem, 0, len(fc.Features))
	for _, f := range fc.Features {
		if f.Properties.CloudCover > maxCloudCover {
			continue
		}
		items = append(items, STACItem{
			ID:         f.ID,
			Datetime:   f.Properties.Datetime,
			CloudCover: f.Properties.CloudCover,
			BBox:       f.BBox,
			Assets:     f.Assets,
		})
	}

	return items, nil
}

// doSearch ejecuta la consulta con reintentos y backoff exponencial ante 5xx.
func (c *Client) doSearch(ctx context.Context, token string, payload []byte) (*http.Response, error) {
	var lastStatus string
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		if attempt > 0 {
			if err := sleep(ctx, c.RetryBaseDelay<<uint(attempt-1)); err != nil {
				return nil, fmt.Errorf("esperar reintento de búsqueda STAC: %w", err)
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.searchEndpoint(), bytes.NewReader(payload))
		if err != nil {
			return nil, fmt.Errorf("construir request STAC: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := c.httpClient().Do(req)
		if err != nil {
			return nil, fmt.Errorf("solicitar búsqueda STAC: %w", err)
		}
		if resp.StatusCode >= 500 {
			lastStatus = resp.Status
			_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 512))
			resp.Body.Close()
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("búsqueda STAC: reintentos agotados tras error %s", lastStatus)
}

// searchEndpoint devuelve el endpoint configurado o, si está vacío, el default.
func (c *Client) searchEndpoint() string {
	if c.SearchURL != "" {
		return c.SearchURL
	}
	return DefaultSearchURL
}

// sleep espera d o hasta que el contexto se cancele.
func sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// parseWKT convierte un WKT POINT o POLYGON en una geometría GeoJSON.
func parseWKT(wkt string) (*geometry, error) {
	s := strings.TrimSpace(wkt)
	upper := strings.ToUpper(s)

	switch {
	case strings.HasPrefix(upper, "POINT"):
		x, y, err := parsePoint(strings.TrimSpace(s[len("POINT"):]))
		if err != nil {
			return nil, err
		}
		return &geometry{Type: "Point", Coordinates: [2]float64{x, y}}, nil
	case strings.HasPrefix(upper, "POLYGON"):
		rings, err := parsePolygon(strings.TrimSpace(s[len("POLYGON"):]))
		if err != nil {
			return nil, err
		}
		return &geometry{Type: "Polygon", Coordinates: rings}, nil
	default:
		return nil, fmt.Errorf("%w: solo se soportan POINT y POLYGON", ErrInvalidAOI)
	}
}

// parsePoint parsea "x y" (con o sin paréntesis) y devuelve sus coordenadas.
func parsePoint(s string) (float64, float64, error) {
	s = strings.Trim(strings.TrimSpace(s), "()")
	fields := strings.FieldsFunc(s, func(r rune) bool { return r == ' ' || r == ',' || r == '\t' })
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("%w: se esperaban 2 coordenadas y se obtuvieron %d", ErrInvalidAOI, len(fields))
	}
	x, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parsear longitud %q: %w", fields[0], err)
	}
	y, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parsear latitud %q: %w", fields[1], err)
	}
	return x, y, nil
}

// parsePolygon parsea "((x y, ...), ...)" y devuelve sus anillos en formato GeoJSON.
func parsePolygon(s string) ([][][2]float64, error) {
	var rings [][][2]float64
	depth := 0
	var cur strings.Builder
	for _, r := range s {
		switch r {
		case '(':
			depth++
			if depth == 2 {
				cur.Reset()
			}
		case ')':
			if depth == 2 {
				ring, err := parseRing(cur.String())
				if err != nil {
					return nil, err
				}
				rings = append(rings, ring)
			}
			depth--
		default:
			if depth == 2 {
				cur.WriteRune(r)
			}
		}
	}
	if len(rings) == 0 {
		return nil, fmt.Errorf("%w: POLYGON sin anillos", ErrInvalidAOI)
	}
	return rings, nil
}

// parseRing parsea "x y, x y, ..." en una lista de pares de coordenadas.
func parseRing(s string) ([][2]float64, error) {
	var pts [][2]float64
	for _, raw := range strings.Split(s, ",") {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		x, y, err := parsePoint(raw)
		if err != nil {
			return nil, err
		}
		pts = append(pts, [2]float64{x, y})
	}
	if len(pts) < 4 {
		return nil, fmt.Errorf("%w: un anillo requiere al menos 4 puntos, se obtuvieron %d", ErrInvalidAOI, len(pts))
	}
	return pts, nil
}

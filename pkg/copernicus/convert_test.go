package copernicus

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestConvertToGeoTIFF_BinarioAusente(t *testing.T) {
	orig := gdalBinary
	gdalBinary = "gdal_translate-que-no-existe-xyz"
	defer func() { gdalBinary = orig }()

	err := ConvertToGeoTIFF(context.Background(), "entrada.jp2", "salida.tif")
	if !errors.Is(err, ErrGDALUnavailable) {
		t.Fatalf("Se esperaba ErrGDALUnavailable, se obtuvo: %v", err)
	}
}

// TestConvertToGeoTIFF_OK sustituye gdal_translate por un script falso que copia
// la entrada a la salida, para verificar la invocación sin depender de GDAL real.
func TestConvertToGeoTIFF_OK(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "falso-gdal")
	if err := os.WriteFile(script, []byte("#!/bin/sh\ncat \"$2\" > \"$3\"\n"), 0o755); err != nil {
		t.Fatalf("escribir script: %v", err)
	}

	orig := gdalBinary
	gdalBinary = script
	defer func() { gdalBinary = orig }()

	src := filepath.Join(dir, "entrada.jp2")
	if err := os.WriteFile(src, []byte("jp2-falso"), 0o644); err != nil {
		t.Fatalf("escribir entrada: %v", err)
	}
	dst := filepath.Join(dir, "salida.tif")

	if err := ConvertToGeoTIFF(context.Background(), src, dst); err != nil {
		t.Fatalf("ConvertToGeoTIFF devolvió error: %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("leer salida: %v", err)
	}
	if string(got) != "jp2-falso" {
		t.Errorf("salida = %q, se esperaba jp2-falso", got)
	}
}

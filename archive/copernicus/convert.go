package copernicus

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ErrGDALUnavailable se devuelve cuando gdal_translate no está disponible.
var ErrGDALUnavailable = errors.New("gdal_translate no disponible")

// gdalBinary es el ejecutable de conversión; variable para sustituirlo en tests.
var gdalBinary = "gdal_translate"

// ConvertToGeoTIFF convierte srcJP2 a GeoTIFF (dstTIF) invocando gdal_translate
// como subproceso. Requiere GDAL instalado en el sistema (no es una dependencia
// de Go del proyecto).
func ConvertToGeoTIFF(ctx context.Context, srcJP2, dstTIF string) error {
	if _, err := exec.LookPath(gdalBinary); err != nil {
		return fmt.Errorf("%w: %v", ErrGDALUnavailable, err)
	}

	out, err := exec.CommandContext(ctx, gdalBinary, "-q", srcJP2, dstTIF).CombinedOutput()
	if err != nil {
		return fmt.Errorf("gdal_translate falló: %w (salida: %s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

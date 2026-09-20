# 05 — Fuentes de datos

Fuente satelital principal: **Element 84 Earth Search**. Al final se listan las alternativas evaluadas y la validación del motor ET0.

## Earth Search — STAC público (fuente principal)

- **Endpoint**: https://earth-search.aws.element84.com/v1
- **Qué es**: catálogo **STAC** público de Element 84 sobre AWS Open Data.
- **Acceso**: **anónimo, sin tokens ni registro**.
- **Colección**: `sentinel-2-l2a` — Sentinel-2 L2A servido como **COG** (Cloud Optimized GeoTIFF).
- **Lectura**: `pystac-client` para buscar por **AOI**, rango de fechas y nubosidad; **`rasterio` con lectura por ventana** (`windowed read`) para leer solo el recorte del AOI — nunca se descarga la escena completa.
- **Bandas**: `B04` (rojo) y `B08` (NIR) para el NDVI; opcionalmente `SCL` para enmascarar nubes y píxeles inválidos. En L2A la reflectancia viene escalada, así que hay que aplicar el `scale`/`offset` que declare cada asset.

### Cálculo

```text
NDVI = (B08 - B04) / (B08 + B04)
Kcb  = 1.51 × NDVI − 0.23
```

- La relación `Kcb = 1.51 × NDVI − 0.23` corresponde al **método INIA** (INIA Intihuasi), validado en Coquimbo.
- TODO: completar la cita exacta del paper (autoría, año, DOI o URL) antes de publicar recomendaciones basadas en esta fórmula.

## Alternativas evaluadas

- **Microsoft Planetary Computer** — catálogo STAC con Sentinel-2 L2A, búsqueda también anónima; buena segunda fuente si Earth Search falla.
- **AWS Open Data** — el dataset crudo detrás de Earth Search: el mismo dato sin depender del catálogo.
- **Mundi Web Services** — catálogo alternativo (DIAS), si más adelante se necesita otra vía de acceso.

## Fuente de validación ET0 (motor Go, se mantiene)

El motor `internal/analysis` sigue calculando ET0 con Hargreaves-Samani; su validación no cambia.

Fuente primaria: **FAO-56** — Allen, R.G., Pereira, L.S., Raes, D. & Smith, M. (1998). *Crop evapotranspiration: Guidelines for computing crop water requirements*. FAO Irrigation and Drainage Paper 56. Roma: FAO.

- Ecuación 52 (Hargreaves-Samani): `ET0 = 0.0023 (Tmedia + 17.8) (Tmax - Tmin)^0.5 Ra`
- Radiación extraterrestre Ra: ecuaciones 21–25 (Cap. 3), con `Gsc = 0.0820 MJ m-2 min-1` y `λ = 2.45 MJ kg-1` (Ra en mm/día = Ra[MJ m-2 d-1] / 2.45).

### Casos de referencia

**Caso 1 — FAO-56, Ejemplo 20 (Lyon, Francia).** Ancla externa publicada.
- Latitud 45.7167°N · 15 de julio (día del año J = 196)
- Tmax = 26.6 °C · Tmin = 14.8 °C · Tmedia = 20.7 °C
- **ET0 esperado = 5.0 mm/día** (FAO-56, ec. 52; el mismo ejemplo da 4.56 mm/día con Penman-Monteith)

**Caso 2 — Valle del Elqui, verano.** Valor calculado aplicando la ec. 52.
- Latitud 30.0°S · 15 de enero (J = 15)
- Tmax = 30.0 °C · Tmin = 15.0 °C · Tmedia = 22.5 °C
- **ET0 esperado = 6.32 mm/día**

**Caso 3 — Valle del Elqui, invierno.** Valor calculado aplicando la ec. 52.
- Latitud 30.0°S · 15 de julio (J = 196)
- Tmax = 18.0 °C · Tmin = 6.0 °C · Tmedia = 12.0 °C
- **ET0 esperado = 1.90 mm/día**

### Validación del término radiativo (Ra)

**FAO-56, Ejemplo 8**: 3 de septiembre (J = 246) a 20°S → Ra = 32.2 MJ m⁻² día⁻¹ = **13.1 mm/día** en evaporación equivalente.

### Pendiente

- TODO: validación cruzada con **INIA** (estaciones del Valle del Elqui) para v2. Los casos 2 y 3 son valores calculados con la ec. 52, no mediciones publicadas; el único ancla externa publicada es el Ejemplo 20 de FAO-56.

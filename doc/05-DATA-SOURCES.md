# 05 — Fuentes de datos

v1 prioriza dos fuentes: **DGA + Sentinel-2**. El resto queda documentado para no perder el hilo.

## DGA — Sistema Hidrométrico

- URL: https://dga.mop.gob.cl/
- Qué aporta: información hidrométrica oficial de la Dirección General de Aguas (caudales, niveles y registros asociados).
- TODO: documentar endpoint o formato de descarga (CSV, Excel, API) y condiciones de uso.
- Rol en v1: **fuente prioritaria**.

## Copernicus Sentinel-2

- Acceso: [Copernicus Data Space](https://dataspace.copernicus.eu/) (registro gratis).
- Qué aporta: imágenes multiespectrales; índices NDVI/NDWI para estimar estrés hídrico de la vegetación.
- TODO: definir el mecanismo de obtención (descarga manual vs. API) y el preprocesamiento mínimo.
- Rol en v1: **fuente prioritaria**.

## CIREN — Coquimbo

- Plataforma de monitoreo hídrico en desarrollo.
- TODO: investigar si publica API o datos descargables.
- Rol: v2+.

## SMAP (NASA)

- Humedad de suelo con resolución de 9 km.
- Acceso: [NASA Earthdata](https://www.earthdata.nasa.gov/) (registro gratis).
- Rol: v2+.

## SMOS (ESA)

- Humedad de suelo con más de 15 años de registros.
- Acceso: Copernicus Data Space (gratis).
- Rol: v2+.

## Estrategia v1

Priorizar **DGA + Sentinel-2**: son las fuentes más directas para el balance hídrico y el estrés hídrico. CIREN, SMAP y SMOS quedan para v2+: su resolución y formatos agregan complejidad que v1 no necesita.

## Fuente de validación ET0

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

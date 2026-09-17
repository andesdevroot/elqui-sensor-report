# 02 — Modelo de datos

Los structs viven en `internal/models`. Los tags JSON usan `snake_case`.

> Contrato de datos de v1. La implementación llega con TDD (Fases 1–2 del roadmap); este documento fija campos, tags e invariantes.

## SensorReading

Lectura individual de un sensor de humedad de suelo.

```go
// SensorReading representa una lectura puntual de un sensor de humedad de suelo.
type SensorReading struct {
    Timestamp  time.Time `json:"timestamp"`
    HumedadPct float64   `json:"humedad_pct"`
    TempC      float64   `json:"temp_c"`
    ParcelaID  string    `json:"parcela_id"`
}
```

## ET0Result

Evapotranspiración de referencia calculada para una fecha.

```go
// ET0Result contiene la ET0 (mm) calculada para una fecha y el método usado.
type ET0Result struct {
    Fecha     time.Time `json:"fecha"`
    ET0mm     float64   `json:"et0_mm"`
    Metodo    string    `json:"metodo"`
    Confianza float64   `json:"confianza"`
}
```

## EfficiencyReport

Resumen de eficiencia hídrica de una parcela en un período.

```go
// EfficiencyReport resume la eficiencia hídrica de una parcela en un período.
type EfficiencyReport struct {
    Periodo       string  `json:"periodo"`
    ConsumoRealMm float64 `json:"consumo_real_mm"`
    OptimoMm      float64 `json:"optimo_mm"`
    EficienciaPct float64 `json:"eficiencia_pct"`
    Semaforo      string  `json:"semaforo"`
}
```

## Invariantes

- `HumedadPct` ∈ [0, 100]
- `TempC` ∈ [-10, 50]
- `EficienciaPct` ∈ [0, ∞)

## TODOs

- TODO: definir dónde se validan los invariantes (¿`ingest` o `models`?) y qué pasa con una lectura fuera de rango: ¿error o advertencia?
- TODO: definir los valores y umbrales de `Semaforo` (p. ej. verde/amarillo/rojo).
- TODO: definir la semántica de `Metodo` y el rango de `Confianza` (¿0.0–1.0?).
- TODO: definir el formato de `Periodo` (p. ej. `2026-03` o `2026-03-01..2026-03-31`).

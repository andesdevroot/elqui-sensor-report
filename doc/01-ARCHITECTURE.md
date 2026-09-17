# 01 — Arquitectura

## Capas

```text
┌──────────────────────────────────────────────────────────────┐
│                          cmd/elqui                           │
│                    main.go — entry point                     │
└───────────────────────────────┬──────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────┐
│                          internal/                           │
│                                                              │
│       ingest ───────────► analysis ───────────► report       │
│       CSV → lecturas      ET0 + balance         Markdown     │
└───────────────────────────────┬──────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────┐
│                             pkg/                             │
│                                                              │
│        dga (datos DGA)          sentinel (NDVI/NDWI)         │
└───────────────────────────────┬──────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────┐
│                       internal/models                        │
│         SensorReading · ET0Result · EfficiencyReport         │
└──────────────────────────────────────────────────────────────┘
```

**Regla de dependencias**: `cmd` orquesta e importa `internal/` y `pkg/`; `internal/` importa `pkg/` y `models`; `pkg/` no importa `internal/`; `internal/models` no importa nada del proyecto.

## Decisión: Go

- **stdlib robusta**: `encoding/csv`, `encoding/json`, `net/http` y `time` cubren todo lo necesario en v1.
- **Binario estático**: `go build -ldflags="-s -w"` produce un único binario, sin runtime ni instalación.
- **Consistencia con promptc**: mismo stack, estructura y metodología que el proyecto hermano del mismo autor (`github.com/andesdevroot/promptc`).

## Principio: sin dependencias externas en v1

`go.mod` no declara dependencias. Si una tarea parece requerir una librería externa, primero se busca una alternativa con stdlib; si no existe, la decisión se documenta en este archivo antes de agregarla.

## Flujo de datos

```text
CSV del sensor → []SensorReading → ET0Result → balance hídrico → EfficiencyReport → Markdown
```

1. `ingest` parsea el CSV del sensor y valida invariantes → `[]SensorReading`.
2. `analysis` calcula la ET0 (Hargreaves-Samani) → `ET0Result`.
3. `analysis` compara el consumo real contra el óptimo del período → `EfficiencyReport` (con semáforo).
4. `report` renderiza el Markdown final que recibe el agricultor.

Las fuentes externas (DGA, Sentinel-2) se consultan desde `pkg/`; su estado y TODOs viven en `doc/05-DATA-SOURCES.md`.

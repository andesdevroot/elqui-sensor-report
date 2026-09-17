# 03 — Metodología de trabajo

## Spec-Driven Development

- `doc/` es la **fuente de verdad** del proyecto: qué se construye, por qué y con qué criterios.
- `Task.md` es la **cola de trabajo**: las tareas pendientes y su estado.
- El código se escribe para cumplir la especificación. Si la spec cambia, se actualiza `doc/` en el mismo commit (o en un commit previo).

## TDD estricto: RED → GREEN → REFACTOR

1. **RED**: escribir un test que falle y demuestre la carencia.
2. **GREEN**: implementar el código mínimo que hace pasar el test.
3. **REFACTOR**: limpiar el código manteniendo la suite en verde.

## Flujo por tarea

```text
leer doc/ → escribir test → go test (RED) → implementar → go test (GREEN) → refactor → commit
```

## Commits atómicos

- Un commit por tarea completada. Si una tarea no cabe en un commit, se divide en tareas más pequeñas.
- Mensaje en formato **Conventional Commits** (formato exacto en `doc/04-CONVENTIONS.md`).

## Regla de oro

Nunca hacer commit con tests rojos: `go test -v ./...` en verde es requisito previo de cada commit.

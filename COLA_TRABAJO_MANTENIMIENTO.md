# Documentación: Cola de Trabajo (Pendientes de Mantenimiento)

Este endpoint especializado está diseñado para el personal técnico o de supervisión que necesita ver **qué tareas tienen pendientes de atender** hoy mismo.

A diferencia del listado general, este endpoint filtra automáticamente los tickets terminados y los **ordena por urgencia**.

## 1. Ver Cola de Trabajo
Muestra los tickets con estado `abierto` o `en_progreso`.

- **Ruta:** `GET /api/user/tickets/cola-trabajo`
- **Autenticación:** Requerida (Bearer Token)

### Parámetros de Consulta (Query Params)
- `pag`: (Opcional) Número de página (default: 1).
- `buscar`: (Opcional) Búsqueda por asunto o descripción.
- `propiedad_id`: (Opcional) Filtrar tareas de un solo inmueble.

### Lógica de Ordenamiento
El sistema ordena los resultados automáticamente bajo este criterio:
1.  **Prioridad (Descendente)**: Primero verás los de prioridad `alta`, luego `media` y al final `baja`.
2.  **Antigüedad (Ascendente)**: Dentro de la misma prioridad, verás primero los tickets más antiguos (los que llevan más tiempo esperando).

---

### Ejemplo de Respuesta (JSON)
```json
{
  "datos": [
    {
      "id": 12,
      "unidad_nombre": "Departamento 401",
      "asunto": "Corte de luz total",
      "descripcion": "No hay energía en ninguna toma...",
      "prioridad": "alta",
      "estado": "abierto",
      "fecha_apertura": "2026-05-08T10:00:00Z"
    },
    {
      "id": 5,
      "unidad_nombre": "Oficina 2",
      "asunto": "Fuga leve",
      "prioridad": "media",
      "estado": "en_progreso",
      "fecha_apertura": "2026-05-07T09:00:00Z"
    }
  ],
  "paginacion": {
    "total": 2,
    "pagina_actual": 1,
    "por_pagina": 10,
    "paginas": 1
  }
}
```

---

## 2. Recomendación de Uso
Para el personal de mantenimiento, se recomienda usar este endpoint como la **pantalla principal de tareas**. 
*   **Acción 1**: Ver la lista y elegir el ticket de arriba (el más urgente).
*   **Acción 2**: Cambiar el estado a `en_progreso` usando el endpoint de `PATCH /estado`.
*   **Acción 3**: Al terminar, cambiar el estado a `resuelto`. El ticket desaparecerá automáticamente de esta lista.

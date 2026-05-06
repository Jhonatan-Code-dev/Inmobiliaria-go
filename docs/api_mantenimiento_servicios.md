# Documentación de API: Mantenimiento y Servicios (1.7)

Este módulo gestiona el ciclo de vida de los incidentes técnicos (Tickets) y el control de mediciones de servicios variables (Agua, Luz, etc.) vinculados a los contratos.

---

## 🛠️ 1. Mantenimiento (Tickets)

Los tickets permiten reportar averías o requerimientos de servicio en una unidad específica.

### 1.1 Listar Tickets
`GET /api/user/tickets`

**Filtros (Query Params):**
- `pag`: (int) Página actual.
- `por_pagina`: (int) Registros por página.
- `unidad_id`: (int) Filtrar por una habitación/local específico.
- `estado`: (string) `abierto`, `en_progreso`, `resuelto`, `cancelado`.

**Respuesta Exitosa (200 OK):**
```json
{
  "datos": [
    {
      "id": 10,
      "unidad_id": 3,
      "unidad": "A-101",
      "titulo": "Fuga de agua en baño",
      "descripcion": "Goteo constante en la ducha principal",
      "prioridad": "alta",
      "estado": "abierto",
      "creado_en": "2026-04-12T10:00:00Z"
    }
  ],
  "paginacion": { ... }
}
```

---

### 1.2 Crear Ticket
`POST /api/user/tickets`

**Request Body:**
```json
{
  "unidad_id": 3,
  "titulo": "Falla eléctrica",
  "descripcion": "No hay luz en la cocina, parece ser un breaker.",
  "prioridad": "media"
}
```

---

### 1.3 Actualizar Ticket
`PUT /api/user/tickets/:id`

Permite cambiar datos del reporte o avanzar su estado.

**Request Body:**
```json
{
  "titulo": "Falla eléctrica (Actualizado)",
  "descripcion": "Se requiere cambio de cableado.",
  "estado": "en_progreso"
}
```

---

## 💧 2. Servicios (Mediciones)

Módulo para el registro de lecturas de medidores (agua, luz) para posterior facturación basada en consumo.

### 2.1 Listar Mediciones
`GET /api/user/servicios`

**Filtros (Query Params):**
- `contrato_id`: (int) Ver lecturas asociadas a un inquilino específico.
- `tipo`: (string) `agua`, `luz`.

---

### 2.2 Registrar Lectura (Crear)
`POST /api/user/servicios`

**Request Body:**
```json
{
  "unidad_id": 3,
  "contrato_id": 1,
  "tipo": "luz",
  "mes": 4,
  "anio": 2026,
  "lectura_anterior": 1200.50,
  "lectura_actual": 1280.75
}
```
*Nota: El sistema calculará el consumo automáticamente (80.25 unidades).*

---

### 2.3 Corregir Lectura (Actualizar)
`PUT /api/user/servicios/:id`

Si se cometió un error en la digitación de la lectura actual.

**Request Body:**
```json
{
  "lectura_actual": 1285.00
}
```

---

## 💡 Instrucciones para el Consumo (Frontend)

1. **Prioridad de Tickets:** El campo `prioridad` es opcional pero recomendado. Los valores estándar son `baja`, `media`, `alta`.
2. **Ciclo de Estados:** Se recomienda que el frontend muestre botones de acción (Ej: "Marcar en Progreso", "Cerrar Ticket") que envíen el estado correspondiente al endpoint `PUT`.
3. **Lecturas:** Antes de registrar una lectura actual, se recomienda obtener la última lectura registrada para mostrarla como "Lectura Anterior" en el formulario.

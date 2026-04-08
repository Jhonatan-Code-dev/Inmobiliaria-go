# API de Gastos (Módulo de Egresos Simplificado)

Este módulo permite gestionar los gastos de la empresa, visualizar registros históricos con filtros avanzados y paginación.

**Base URL:** `http://localhost:5000/api/user/gastos`

## Listar Gastos
Obtiene una lista paginada de hasta 10 gastos por página.

- **URL:** `GET /`
- **Auth:** Requiere Bearer Token o Cookie (`token_usuario`).
- **Query Params:**
    - `pag` (int): Número de página (ej. `1`).
    - `anio` (int): Filtrar por año (ej. `2024`).
    - `mes` (int): Filtrar por mes (1-12).
    - `desde` (string): Fecha inicio `YYYY-MM-DD`.
    - `hasta` (string): Fecha fin `YYYY-MM-DD`.
    - `fecha` (string): Fecha exacta `YYYY-MM-DD`.

### Respuesta (200 OK)
```json
{
  "datos": [
    {
      "id": 1,
      "empresa_id": 1,
      "monto": 50.00,
      "fecha": "2024-03-10T00:00:00Z",
      "tipo_pago_id": 1,
      "descripcion": "Reparación de tubería"
    }
  ],
  "paginacion": {
    "total": 1,
    "paginas": 1,
    "pagina": 1,
    "por_pagina": 10
  }
}
```

---

## Registrar Gasto
Crea un nuevo registro de gasto y genera directamente el registro en caja.

- **URL:** `POST /`
- **Body:**
```json
{
  "empresa_id": 1,
  "monto": 50.50,
  "fecha": "2024-03-01",
  "tipo_pago_id": 3,
  "descripcion": "Pago mensual Starlink"
}
```
*Nota: El monto se envía y se recibe como número decimal (ej. 50.50).*

### Respuesta (201 Created)
Retorna el objeto creado con su ID.

---

## Listar Tipos de Pago
Obtiene la lista de los métodos de pago disponibles en la base de datos.

- **URL:** `GET /tipos-pago`
- **Respuesta:**
```json
[
  {"id": 1, "nombre": "efectivo"},
  {"id": 2, "nombre": "transferencia"},
  {"id": 3, "nombre": "yape"}
]
```

---

## Actualizar Gasto
Actualiza los datos de un gasto existente.

- **URL:** `PUT /{id}`
- **Body:** Similar al de creación.

---

## Eliminar Gasto
Elimina físicamente el registro de gasto.

- **URL:** `DELETE /{id}`
- **Respuesta:** `{"message": "gasto eliminado"}`

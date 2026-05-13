# API de Mediciones de Servicios (Luz y Agua)

Esta documentación contiene únicamente los endpoints necesarios para gestionar las lecturas de medidores y la cobranza automática de servicios en las habitaciones.

## 1. Obtener Lectura Anterior
Utilice este endpoint para mostrar al usuario cuál fue la última lectura registrada y así calcular el punto de partida del nuevo mes.

- **URL**: `GET /api/user/servicios/ultimo/:contrato_id`
- **Método**: `GET`
- **Parámetros de URL**: 
    - `contrato_id` (Integer): ID del contrato de alquiler.
- **Parámetros de Consulta (Query)**:
    - `tipo` (String): `luz` o `agua` (Opcional, por defecto `luz`).
- **Respuesta Exitosa (200 OK)**:
```json
{
  "id": 10,
  "contrato_id": 1,
  "tipo_servicio": "luz",
  "lectura_actual": 150.00,
  "fecha_lectura": "2024-04-10T00:00:00Z"
}
```
*Si no existen lecturas previas, el campo `lectura_actual` devolverá `0`.*

---

## 2. Registrar Nueva Lectura y Generar Cobro
Este endpoint realiza tres acciones en un solo paso: registra la lectura, calcula el consumo y genera el cargo (deuda) para el inquilino.

- **URL**: `POST /api/user/servicios/registrar-y-cobrar`
- **Método**: `POST`
- **Cuerpo de la Petición (JSON)**:
```json
{
  "contrato_id": 1,
  "tipo_servicio": "luz",
  "lectura_actual": 200.00,
  "fecha_lectura": "2024-05-12",
  "precio_unitario": 1.80
}
```
- **Campos**:
    - `contrato_id`: ID del contrato.
    - `tipo_servicio`: "luz" o "agua".
    - `lectura_actual`: El valor actual que marca el medidor.
    - `fecha_lectura`: Fecha de la toma de lectura (YYYY-MM-DD).
    - `precio_unitario`: El costo por cada unidad consumida.
- **Respuesta Exitosa (201 Created)**:
```json
{
  "id": 11,
  "contrato_id": 1,
  "tipo_servicio": "luz",
  "lectura_anterior": 150.00,
  "lectura_actual": 200.00,
  "consumo": 50.00,
  "monto": 90.00,
  "procesado": true,
  "cargo_id": 505
}
```

---

## 3. Listar Historial de Mediciones
Para mostrar una tabla con el histórico de consumos de una habitación.

- **URL**: `GET /api/user/servicios`
- **Método**: `GET`
- **Filtros (Query)**:
    - `contrato_id` (Integer): ID del contrato.
    - `pag` (Integer): Número de página.
    - `por_pagina` (Integer): Cantidad de registros por página.
- **Respuesta Exitosa (200 OK)**: Devuelve un objeto paginado con el array de mediciones.

---

## 4. Corregir Lectura Errónea
En caso de error al digitar la lectura actual.

- **URL**: `api/user/servicios/:id`
- **Método**: `PUT`
- **Cuerpo de la Petición (JSON)**:
```json
{
  "lectura_actual": 195.00
}
```
- **Restricción**: El sistema recalculará el consumo y el monto del cargo automáticamente, siempre y cuando el cargo no haya sido procesado o pagado.

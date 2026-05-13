# Documentación Completa: Módulo de Mediciones y Servicios

Esta guía contiene la lista definitiva de endpoints para la gestión integral de servicios (Luz, Agua, etc.) en el sistema de alquileres.

## 1. Consulta de Lectura Anterior
Obtiene el valor del medidor del mes previo para un contrato específico.

- **URL**: `GET /api/user/servicios/ultimo/:contrato_id`
- **Query Params**: `tipo` (luz|agua)
- **Respuesta**: El objeto de la última medición o `{"lectura_actual": 0}` si es nuevo.

---

## 2. Registro Individual con Cobro Automático
Registra una lectura y genera la deuda en el sistema financiero instantáneamente.

- **URL**: `POST /api/user/servicios/registrar-y-cobrar`
- **Body**:
```json
{
  "contrato_id": 1,
  "tipo_servicio": "luz",
  "lectura_actual": 200.00,
  "fecha_lectura": "2024-05-12",
  "precio_unitario": 1.50
}
```

---

## 3. Registro Masivo (Carga por Lote)
Permite registrar las lecturas de múltiples habitaciones en una sola petición. Ideal para fin de mes.

- **URL**: `POST /api/user/servicios/masivo`
- **Body**: Un array de objetos como el del punto anterior.
```json
[
  { "contrato_id": 1, "tipo_servicio": "luz", "lectura_actual": 200, ... },
  { "contrato_id": 2, "tipo_servicio": "luz", "lectura_actual": 185, ... }
]
```

---

## 4. Listar Unidades Pendientes de Lectura
Muestra qué habitaciones aún no tienen registrada su lectura para el mes en curso.

- **URL**: `GET /api/user/servicios/pendientes`
- **Query Params**: `tipo` (luz|agua)

---

## 5. Selector de Contratos Activos
Endpoint ligero diseñado específicamente para llenar desplegables (dropdowns) en el frontend. Retorna solo contratos con estado 'activo' o 'vencido'.

- **URL**: `GET /api/user/alquileres/activos/selector`
- **Respuesta**: Array de objetos con ID, Nombre de Cliente y Código de Unidad.

---

## 6. Historial y Mantenimiento
- **Listar Todo**: `GET /api/user/servicios?contrato_id=X` (Paginado)
- **Ver Detalle**: `GET /api/user/servicios/:id`
- **Actualizar/Corregir**: `PUT /api/user/servicios/:id` (Solo lectura actual)
- **Eliminar**: `DELETE /api/user/servicios/:id` (Solo si no está pagado)

---

## Pruebas de Funcionamiento (Test Report)
El módulo ha sido validado con pruebas unitarias (`TestRegistrarYCobrar`):
- [x] Validación de cálculo de consumo (Actual - Anterior).
- [x] Validación de monto financiero (Consumo * Precio).
- [x] Generación automática de Cargo en el estado de cuenta.
- [x] Integridad de datos en base de datos.

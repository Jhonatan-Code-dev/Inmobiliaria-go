# Documentación Profesional: Módulo de Mediciones y Servicios (V2)

Este documento centraliza todos los endpoints necesarios para el registro, cálculo y facturación de servicios variables (Luz, Agua, etc.) en el sistema de Rentals.

## 1. Consultar Estado Actual (Pre-llenado)
Antes de registrar una lectura, el frontend debe obtener los datos históricos para calcular el consumo esperado.

- **URL**: `GET /api/user/servicios/ultimo/:contrato_id`
- **Query Params**: `tipo` (opcional, default `luz`)
- **Respuesta**:
```json
{
  "contrato_id": 12,
  "tipo_servicio": "luz",
  "lectura_actual": 1250.50,
  "fecha_lectura": "2024-04-12T00:00:00Z"
}
```

---

## 2. Registrar Lectura y Generar Cobro (Directo)
Este endpoint registra la medición y **crea automáticamente una deuda** en el estado de cuenta del inquilino.

- **URL**: `POST /api/user/servicios/registrar-y-cobrar`
- **Cuerpo (JSON)**:
```json
{
  "contrato_id": 12,        // Acepta número o string "12"
  "tipo_servicio": "luz",   // "luz" o "agua"
  "lectura_actual": 1300.0, // Valor final del medidor
  "lectura_anterior": 1250, // Opcional. Si no se envía, usa la última lectura del sistema.
  "precio_unitario": 1.50,  // Costo por kWh o m3
  "factor": 1.0,            // Opcional. Multiplicador del medidor (ej. x10). Default 1.0
  "cargo_fijo": 5.00,       // Opcional. Monto extra (Mantenimiento/Cargo Fijo). Default 0.0
  "fecha_lectura": "2024-05-13" 
}
```
### Fórmula de Cálculo Aplicada:
`Consumo = (Lectura Actual - Lectura Anterior) * Factor`
`Monto Final = (Consumo * Precio Unitario) + Cargo Fijo`

### 2.1. Inicialización de Medidor (Primer Registro)
Si es la primera vez que registras el medidor y este no empieza en cero (ej. ya marca `1250`), usa el campo `lectura_anterior` para establecer el punto de partida.

- **Cuerpo (JSON)**:
```json
{
  "contrato_id": 12,
  "tipo_servicio": "luz",
  "lectura_actual": 1300,
  "lectura_anterior": 1250, // "Cuanto tenía" al iniciar
  "precio_unitario": 1.50
}
```
*El sistema calculará el consumo sobre la diferencia (50 unidades) y guardará `1300` como la nueva base para el próximo mes.*

---

## 3. Registro Masivo
Ideal para cargar todas las habitaciones de un edificio en un solo paso.

- **URL**: `POST /api/user/servicios/masivo`
- **Cuerpo**: Array de objetos con la misma estructura que el registro individual.

---

## 4. Selector de Contratos Activos (Dropdown Helper)
Para llenar los selectores del frontend con contratos válidos.

- **URL**: `GET /api/user/alquileres/activos/selector`
- **Respuesta**:
```json
[
  { "id": 12, "cliente_nombre": "Juan Perez", "unidad_codigo": "HAB-201" },
  { "id": 13, "cliente_nombre": "Maria Lopez", "unidad_codigo": "HAB-202" }
]
```

---

## 5. Gestión y Correcciones
- **Listar Pendientes**: `GET /api/user/servicios/pendientes` (Muestra lecturas aún no cobradas).
- **Eliminar**: `DELETE /api/user/servicios/:id` (Solo si el cargo asociado no ha sido pagado).
- **Actualizar**: `PUT /api/user/servicios/:id` (Permite corregir la lectura actual si hubo error de dedo).

---

## Notas Técnicas para Frontend:
1. **Robusto**: El campo `contrato_id` es flexible; puedes enviarlo como número o como texto.
2. **Validación**: El servidor rechazará lecturas donde la `actual` sea menor a la `anterior` (después de aplicar el factor).
3. **Cargos**: Cada registro genera un registro en la tabla `cargos` con concepto "Consumo de [Tipo]".

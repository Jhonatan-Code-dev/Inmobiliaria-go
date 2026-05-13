# Documentación de Servicios para Frontend (Luz y Agua)

Esta guía detalla cómo integrar los nuevos endpoints para el registro de lecturas de medidores y la generación automática de cobros.

## 1. Consultar Última Lectura (Lectura Anterior)

Antes de registrar una nueva lectura, usa este endpoint para obtener el valor del medidor del mes pasado. Esto permite que el usuario solo tenga que ingresar el valor actual.

- **URL**: `GET /api/user/servicios/ultimo/:contrato_id`
- **Query Params**: 
    - `tipo`: `luz` o `agua` (por defecto es `luz`)
- **Respuesta Exitosa (200 OK)**:
```json
{
  "id": 45,
  "contrato_id": 12,
  "tipo_servicio": "luz",
  "lectura_actual": 1250.50,
  "fecha_lectura": "2024-04-12T00:00:00Z"
}
```
*Nota: Si no hay lecturas previas, devolverá `{"lectura_actual": 0}`.*

---

## 2. Registrar Lectura y Generar Cobro Automático

Este es el endpoint principal para el proceso "Directo". Registra la lectura, calcula el consumo y crea la deuda en el sistema financiero.

- **URL**: `POST /api/user/servicios/registrar-y-cobrar`
- **Cuerpo de la Petición (JSON)**:
```json
{
  "contrato_id": 12,
  "tipo_servicio": "luz",
  "lectura_actual": 1300.00,
  "fecha_lectura": "2024-05-12",
  "precio_unitario": 1.50
}
```
- **Lógica Interna (Informativo)**:
    - `Consumo = 1300.00 - 1250.50 = 49.50`
    - `Monto = 49.50 * 1.50 = 74.25`
    - Se crea un Cargo con concepto "Consumo de Luz" por S/. 74.25.

- **Respuesta Exitosa (201 Created)**:
```json
{
  "id": 46,
  "contrato_id": 12,
  "tipo_servicio": "luz",
  "lectura_anterior": 1250.50,
  "lectura_actual": 1300.00,
  "consumo": 49.50,
  "monto": 74.25,
  "procesado": true,
  "cargo_id": 1025
}
```

---

## 3. Flujo Sugerido en la Interfaz (UI)

1. **Paso 1**: El usuario selecciona el inquilino/habitación (Contrato).
2. **Paso 2**: El frontend llama a `GET /ultimo/:id?tipo=luz`.
3. **Paso 3**: Se muestra en pantalla: *"Lectura Anterior: 1250.50"*.
4. **Paso 4**: El usuario ingresa la *"Lectura Actual"* (ej. 1300.00).
5. **Paso 5**: El usuario ingresa el *"Precio por Unidad"* (ej. 1.50).
6. **Paso 6**: El usuario presiona "Guardar y Cobrar".
7. **Paso 7**: Se llama a `POST /registrar-y-cobrar`.
8. **Paso 8**: Éxito. El sistema muestra un mensaje: *"Lectura registrada. Se ha generado un cargo de S/. 74.25"*.

---

## 4. Endpoints de Soporte (Mantenimiento)

### Listar Historial de Mediciones
- **URL**: `GET /api/user/servicios?contrato_id=12&pag=1&por_pagina=10`
- **Uso**: Para mostrar una tabla con todos los meses de luz/agua de un inquilino.

### Corregir Error
- **URL**: `PUT /api/user/servicios/:id`
- **Body**: `{"lectura_actual": 1295.00}`
- **Restricción**: Solo se puede editar si no ha sido procesado o si el cargo asociado sigue pendiente.

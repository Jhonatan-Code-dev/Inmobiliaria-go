# Documentación: Módulo de Contratos Digitales

Este módulo permite la gestión de plantillas de contratos y la generación automática de documentos de alquiler basados en datos reales de clientes y unidades.

## 1. Gestión de Plantillas

### Listar Plantillas
`GET /api/user/alquileres/plantillas`

**Descripción:** Retorna todas las plantillas guardadas por la empresa.

---

### Guardar/Actualizar Plantilla
`POST /api/user/alquileres/plantillas`

**Payload:**
```json
{
  "id": 0,
  "nombre": "Contrato Estándar Residencial",
  "contenido": "# CONTRATO DE ALQUILER\n\nYo, {{cliente_nombre}}..."
}
```
*   Si `id` es `0`, crea una nueva.
*   Si `id` es mayor a `0`, actualiza la existente.

---

### Eliminar Plantilla
`DELETE /api/user/alquileres/plantillas/:id`

---

## 2. Generación de Contratos

### Generar Documento de Contrato
`GET /api/user/alquileres/:id/generar-documento`

**Query Parameters:**
*   `plantilla_id` (opcional): ID de la plantilla a usar. Si no se envía, usa una plantilla básica por defecto.

**Respuesta:**
```json
{
  "alquiler_id": 15,
  "contenido": "# CONTRATO DE ALQUILER\n\nPor el presente documento..."
}
```

---

## 3. Variables Soportadas (Placeholders)

Puedes usar las siguientes variables en tus plantillas:

| Variable | Descripción |
| :--- | :--- |
| `{{cliente_nombre}}` | Nombre completo del cliente |
| `{{cliente_documento}}` | Número de documento (DNI/RUC) |
| `{{unidad_codigo}}` | Código de la unidad alquilada |
| `{{monto_renta}}` | Monto de la renta mensual |
| `{{moneda}}` | Moneda del contrato (PEN/USD) |
| `{{fecha_inicio}}` | Fecha de inicio (YYYY-MM-DD) |
| `{{dia_vencimiento}}` | Día del mes para el pago |

---

## Guía de Consumo Frontend

1.  **Vista de Alquiler:** Agregar un botón "Generar Contrato".
2.  **Selector de Plantilla:** Al hacer clic, mostrar un modal para elegir una plantilla (opcional).
3.  **Visualización:** Llamar al endpoint de generación y mostrar el `contenido` (Markdown o Texto Plano) en un componente que permita imprimir o copiar.

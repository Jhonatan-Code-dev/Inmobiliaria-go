# Documentación: Configuración Global de Asistencia

Esta documentación detalla los nuevos endpoints para gestionar la política de asistencia de la empresa.

## 1. Concepto General
El sistema ahora permite establecer un **horario base** para toda la empresa. 
- Si un trabajador tiene un horario específico asignado, se usará ese.
- Si un trabajador **no** tiene un horario asignado, el sistema usará esta configuración global para determinar si su marca de entrada es "puntual" o "tarde".

---

## 2. Obtener Configuración Actual
Obtiene los parámetros actuales de la política de asistencia de la empresa.

- **URL:** `/api/user/asistencia/configuracion`
- **Método:** `GET`
- **Query Params:**
    - `empresa_id` (obligatorio): ID de la empresa.
- **Headers:**
    - `Authorization`: Bearer `<JWT_TOKEN>`

### Respuesta Exitosa (200 OK)
```json
{
  "hora_entrada": "08:00",
  "hora_salida": "17:00",
  "tolerancia_minutos": 15,
  "dias_laborables": "1,2,3,4,5"
}
```

---

## 3. Actualizar Configuración Global
Establece o actualiza los parámetros de asistencia para todos los trabajadores sin horario específico.

- **URL:** `/api/user/asistencia/configuracion`
- **Método:** `POST`
- **Query Params:**
    - `empresa_id` (obligatorio): ID de la empresa.
- **Headers:**
    - `Authorization`: Bearer `<JWT_TOKEN>`
    - `Content-Type`: `application/json`

### Cuerpo de la Petición (JSON)
| Campo | Tipo | Descripción | Ejemplo |
| :--- | :--- | :--- | :--- |
| `hora_entrada` | string | Hora de entrada esperada (HH:mm) | `"08:30"` |
| `hora_salida` | string | Hora de salida esperada (HH:mm) | `"18:00"` |
| `tolerancia_minutos` | int | Minutos de gracia antes de marcar "tarde" | `10` |
| `dias_laborables` | string | Días de la semana (1=Lunes, 7=Domingo) | `"1,2,3,4,5,6"` |

### Ejemplo de Petición
```json
{
  "hora_entrada": "09:00",
  "hora_salida": "18:00",
  "tolerancia_minutos": 10,
  "dias_laborables": "1,2,3,4,5"
}
```

### Respuesta Exitosa (200 OK)
Retorna el objeto de configuración actualizado.

---

## 4. Notas para el Frontend
- **Validación:** Asegurarse de enviar las horas en formato de 24h (`HH:mm`).
- **Días Laborables:** Es una cadena separada por comas. El frontend puede presentar un selector de días (Checkboxes) y unir los valores seleccionados.
---

## 5. Reporte de Asistencia Enriquecido
El endpoint de reporte ahora incluye información del horario esperado para facilitar la visualización en el frontend.

- **URL:** `/api/user/asistencia/reporte`
- **Método:** `GET`

### Nuevos Campos en cada registro:
| Campo | Tipo | Descripción |
| :--- | :--- | :--- |
| `hora_entrada_esperada` | string | Hora de entrada según la config. global o individual. |
| `hora_salida_esperada` | string | Hora de salida según la config. global o individual. |

**Nota:** El campo `estado` se re-evalúa dinámicamente en el reporte basándose en la configuración actual para garantizar coherencia visual.

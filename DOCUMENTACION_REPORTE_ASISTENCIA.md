# Documentación: Reporte Detallado de Asistencia

Este endpoint permite obtener un listado exhaustivo de las marcas de asistencia de los trabajadores, con capacidades avanzadas de filtrado y búsqueda.

## 1. Endpoint
- **URL:** `/api/user/asistencia/reporte`
- **Método:** `GET`
- **Autenticación:** Requerida (Bearer Token)

---

## 2. Parámetros de Consulta (Query Params)

| Parámetro | Tipo | Descripción | Ejemplo |
| :--- | :--- | :--- | :--- |
| `empresa_id` | int | **Obligatorio**. ID de la empresa. | `1` |
| `buscar` | string | Busca por el nombre del trabajador (coincidencia parcial). | `"yona"` |
| `fecha` | string | Filtra por una fecha específica (formato YYYY-MM-DD). | `"2026-05-09"` |
| `desde` | string | Fecha inicial para rango (YYYY-MM-DD). | `"2026-05-01"` |
| `hasta` | string | Fecha final para rango (YYYY-MM-DD). | `"2026-05-31"` |
| `estado` | string | Filtra por estado: `puntual`, `tarde`, `falta`, `justificado`. | `"tarde"` |
| `pag` | int | Número de página (default: 1). | `1` |
| `limite` | int | Registros por página (default: 50). | `15` |

---

## 3. Ejemplo de Respuesta (200 OK)

```json
{
    "success": true,
    "data": [
        {
            "id": 5,
            "empresa_id": 1,
            "usuario_id": 1,
            "usuario_nombre": "yona",
            "fecha": "2026-05-09T00:00:00Z",
            "hora_entrada": "2026-05-09T16:40:19Z",
            "hora_salida": "2026-05-09T17:31:00Z",
            "estado": "tarde",
            "notas": null,
            "horas_trabajadas": 0.84,
            "hora_entrada_esperada": "08:00",
            "hora_salida_esperada": "17:00"
        }
    ],
    "total": 3,
    "pagina": 1,
    "limite": 15
}
```

---

## 4. Lógica de Filtrado de Fechas
1. **Fecha Específica:** Si se envía el parámetro `fecha`, el sistema ignorará `desde` y `hasta`, y devolverá solo los registros de ese día.
2. **Rango de Fechas:** Si se envían `desde` y `hasta`, el sistema devolverá todos los registros contenidos en ese intervalo (inclusive).
3. **Búsqueda por Nombre:** El parámetro `buscar` filtra la lista en tiempo real para mostrar solo los trabajadores cuyo nombre coincida con el texto enviado.

---

## 5. Integración con el Frontend
- Se recomienda usar el parámetro `buscar` para implementar una barra de búsqueda en tiempo real.
- Para el filtrado de fechas, se puede usar un componente de calendario que envíe `desde` y `hasta` al seleccionar un rango, o simplemente `fecha` para una vista diaria.
- El campo `estado` se calcula dinámicamente basándose en la configuración global de la empresa si el trabajador no tiene un horario específico.

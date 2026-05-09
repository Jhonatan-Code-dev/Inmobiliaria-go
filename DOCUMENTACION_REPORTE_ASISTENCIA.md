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
            "horas_trabajadas_formato": "0H 50M 24S",
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

---

## 6. Endpoints de Exportación

El sistema provee dos endpoints adicionales para descargar el mismo reporte en formatos binarios. Estos endpoints **soportan exactamente los mismos parámetros de filtro** (`empresa_id`, `buscar`, `fecha`, `desde`, `hasta`, `estado`) que el reporte principal, asegurando que lo exportado coincida con lo visualizado. Se ignora la paginación para exportar el reporte completo.

### 6.1. Exportar a Excel
- **URL:** `/api/user/asistencia/reporte/excel`
- **Método:** `GET`
- **Content-Type Retornado:** `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`

### 6.2. Exportar a PDF
- **URL:** `/api/user/asistencia/reporte/pdf`
- **Método:** `GET`
- **Content-Type Retornado:** `application/pdf`

### Ejemplo de Consumo en Frontend (Axios / Blob)
Al consumir estos endpoints de exportación desde el frontend, **es crítico configurar `responseType: 'blob'`** para no corromper el archivo descargado.

```javascript
import axios from 'axios';

async function descargarReporte(tipo = 'excel', filtros) {
    const queryParams = new URLSearchParams(filtros).toString();
    const url = `/api/user/asistencia/reporte/${tipo}?${queryParams}`;

    const response = await axios.get(url, {
        headers: { Authorization: `Bearer TU_TOKEN` },
        responseType: 'blob', // OBLIGATORIO para archivos binarios
    });

    const urlBlob = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = urlBlob;
    link.setAttribute('download', `reporte_asistencia.${tipo === 'excel' ? 'xlsx' : 'pdf'}`);
    document.body.appendChild(link);
    link.click();
    link.remove();
}
```

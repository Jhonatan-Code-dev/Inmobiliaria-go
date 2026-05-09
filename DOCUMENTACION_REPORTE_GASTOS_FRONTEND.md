# Guía de Consumo: Reportes de Gastos (Excel & PDF)

Esta documentación detalla cómo integrar los nuevos endpoints de reportes de gastos en el frontend.

## 1. Endpoints Disponibles

| Formato | Método | Endpoint | Content-Type |
| :--- | :--- | :--- | :--- |
| **Excel** | `GET` | `/api/user/gastos/reporte/excel` | `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet` |
| **PDF** | `GET` | `/api/user/gastos/reporte/pdf` | `application/pdf` |

## 2. Parámetros de Consulta (Query Params)

Los endpoints aceptan los mismos filtros que el listado normal de gastos. Todos son opcionales.

| Parámetro | Tipo | Descripción | Ejemplo |
| :--- | :--- | :--- | :--- |
| `empresa_id` | `int` | **Obligatorio**. ID de la empresa. | `?empresa_id=1` |
| `anio` | `int` | Filtra gastos por un año específico. | `&anio=2024` |
| `mes` | `int` | Filtra gastos por un mes (1-12). | `&mes=4` |
| `desde` | `string` | Fecha de inicio del rango (YYYY-MM-DD). | `&desde=2024-01-01` |
| `hasta` | `string` | Fecha de fin del rango (YYYY-MM-DD). | `&hasta=2024-03-31` |

---

## 3. Implementación en el Frontend

### Caso A: Descarga Directa (window.open)
Si tu backend maneja la autenticación por Cookies o no requiere headers personalizados para la descarga, esta es la forma más sencilla:

```javascript
/**
 * @param {'excel' | 'pdf'} formato 
 * @param {Object} filtros 
 */
const descargarReporte = (formato, filtros) => {
  const url = new URL(`${window.location.origin}/api/user/gastos/reporte/${formato}`);
  
  // Agregar filtros a la URL
  Object.keys(filtros).forEach(key => {
    if (filtros[key]) url.searchParams.append(key, filtros[key]);
  });

  // Abrir en nueva pestaña para iniciar descarga
  window.open(url.toString(), '_blank');
};

// Ejemplo de uso:
descargarReporte('pdf', { empresa_id: 1, anio: 2024, mes: 5 });
```

### Caso B: Usando Axios (Con Token Bearer)
Si usas JWT en el Header `Authorization`, debes pedir el archivo como un `blob` para que la librería no intente parsearlo como JSON.

```javascript
import axios from 'axios';

/**
 * @param {'excel' | 'pdf'} formato 
 * @param {Object} filtros 
 */
const exportarGastos = async (formato, filtros) => {
  try {
    const token = localStorage.getItem('token'); // O donde guardes tu JWT

    const response = await axios.get(`/api/user/gastos/reporte/${formato}`, {
      params: filtros,
      responseType: 'blob', // REQUERIDO
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    // Crear un link temporal en el DOM
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;

    // Nombre del archivo
    const extension = formato === 'excel' ? 'xlsx' : 'pdf';
    link.setAttribute('download', `reporte_gastos_${Date.now()}.${extension}`);
    
    document.body.appendChild(link);
    link.click();
    
    // Limpieza
    link.remove();
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error("Error descargando el reporte:", error);
    alert("No se pudo generar el reporte.");
  }
};
```

---

## 4. Notas Técnicas

1.  **Seguridad**: Los endpoints verifican que el `empresa_id` solicitado coincida con el `empresa_id` de la sesión del usuario (vía middleware `TenantAuth`).
2.  **Límites**: El reporte no está paginado. Generará un listado de hasta 10,000 registros por defecto para asegurar que el historial esté completo.
    *   El **PDF** incluye numeración de páginas y un resumen del total al final del documento.

---

## 5. Solución de Problemas: El archivo se descarga sin extensión

Si al descargar el archivo recibes un nombre como `a5f1f611-3dda...` sin extensión, revisa lo siguiente:

### 1. En el Servidor (Backend)
Asegúrate de que el header `Content-Disposition` use comillas dobles para el nombre del archivo. El backend ya ha sido actualizado para enviar:
`Content-Disposition: attachment; filename="reporte_gastos.xlsx"`

### 2. En el Cliente (Frontend - Caso Axios)
Cuando usas `axios` o `fetch` con `responseType: 'blob'`, el navegador pierde el nombre del archivo original enviado por el servidor a menos que lo extraigas de los headers. Para evitar complicaciones, **fuerza el nombre y la extensión** al crear el enlace de descarga:

```javascript
// FORZAR NOMBRE Y EXTENSIÓN CORRECTA
const extension = formato === 'excel' ? 'xlsx' : 'pdf';
link.setAttribute('download', `reporte_gastos.${extension}`);
```

### 3. Atributo `download`
Si usas un enlace `<a>` simple, siempre incluye el atributo `download` con el nombre deseado:
`<a href="/api/..." download="reporte.pdf">Descargar</a>`

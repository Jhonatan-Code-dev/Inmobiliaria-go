# API: Sistema de Control de Asistencia (Documentación Completa)

Esta guía detalla cómo integrar el módulo de asistencia en el Frontend. El sistema maneja marcaciones inteligentes, cálculos automáticos de horas y supervisión administrativa con filtros avanzados.

---

## 1. Configuración General

*   **Base URL:** `{{server}}/api/user/asistencia`
*   **Autenticación:** Requiere `Authorization: Bearer <JWT_TOKEN>`
*   **Contexto Multi-Tenant:** La mayoría de los endpoints requieren el query param `empresa_id`.

---

## 2. Flujo del Trabajador (App Móvil / Panel Usuario)

### 2.1 Marcación Inteligente (Entrada/Salida)
Este es el endpoint principal para el trabajador. El sistema determina automáticamente si es entrada o salida basándose en si ya existe un registro hoy.

*   **Endpoint:** `/marcar`
*   **Método:** `POST`
*   **Body:** `{}` (Vacio)
*   **Lógica:**
    *   **Primer clic del día:** Registra **ENTRADA**. Si pasa de la hora configurada (más tolerancia), el estado será `tarde`.
    *   **Segundo clic del día:** Registra **SALIDA**. Calcula automáticamente `horas_trabajadas`.

### 2.2 Ver Mi Historial
*   **Endpoint:** `/mi-historial`
*   **Método:** `GET`
*   **Descripción:** Retorna todos los registros de asistencia del usuario autenticado.

---

## 3. Flujo del Administrador (Panel de Supervisión)

### 3.1 Supervisión Global de Personal
Visualiza a todos los empleados y filtra por comportamiento.

*   **Endpoint:** `/registros?empresa_id=1`
*   **Método:** `GET`
*   **Parámetros de Filtro:**
    *   `estado`: `tarde`, `puntual`, `falta`. (Ej: `?estado=tarde` para ver infractores).
    *   `desde` / `hasta`: Rango de fechas (`YYYY-MM-DD`).
    *   `pag`: Número de página.
*   **Campo Clave:** `usuario_nombre` viene incluido para mostrarlo en las tablas.

### 3.2 Gestión de Horarios
*   **Asignar/Editar:** `POST /horarios?empresa_id=1`
*   **Ver Detalle:** `GET /horarios/detalle?empresa_id=1&usuario_id=4`
*   **Estructura del Horario:**
    ```json
    {
      "usuario_id": 4,
      "hora_entrada": "08:00",
      "hora_salida": "17:00",
      "tolerancia_minutos": 15,
      "dias_laborables": "1,2,3,4,5" // Lunes a Viernes
    }
    ```

### 3.3 Corrección de Errores
Si un usuario marcó mal, el administrador puede eliminar la marca para permitir una nueva.
*   **Endpoint:** `/registros/:id?empresa_id=1`
*   **Método:** `DELETE`

---

## 4. Gestión de Permisos y Justificaciones

1.  **Trabajador solicita:** `POST /permisos` enviando `fecha` y `motivo`.
2.  **Admin visualiza:** `GET /permisos?empresa_id=1&estado=pendiente`.
3.  **Admin decide:** `PUT /permisos/:id/estado?empresa_id=1` enviando `{"estado": "aprobado", "respuesta": "..."}`.

---

## 5. Ejemplo de Consumo (React + Axios)

### Registrar Asistencia (Botón Único)
```javascript
const handleMarcado = async () => {
  try {
    const res = await axios.post('/api/user/asistencia/marcar');
    const data = res.data;
    
    if (data.hora_salida) {
      alert(`Salida registrada. Trabajaste ${data.horas_trabajadas.toFixed(2)} horas.`);
    } else {
      const estadoMsg = data.estado === 'tarde' ? ' (Llegaste Tarde)' : ' (A tiempo)';
      alert(`Entrada registrada: ${new Date(data.hora_entrada).toLocaleTimeString()}${estadoMsg}`);
    }
  } catch (err) {
    alert("Error: " + err.response.data.message);
  }
};
```

### Tabla de Supervisión (Admin)
```javascript
// Obtener todos los que llegaron tarde hoy
const fetchTardanzas = async () => {
  const hoy = new Date().toISOString().split('T')[0];
  const res = await axios.get(`/api/user/asistencia/registros`, {
    params: {
      empresa_id: 1,
      estado: 'tarde',
      desde: hoy,
      hasta: hoy
    }
  });
  setTardanzas(res.data);
};
```

---

## Notas Técnicas para el Frontend:
1.  **Zonas Horarias:** El backend ya maneja la conversión a la hora local de la empresa según su país. El Frontend solo debe mostrar las fechas que recibe.
2.  **Estados Visuales:** Se recomienda usar colores: Verde (`puntual`), Naranja/Rojo (`tarde`), Gris (`falta`).
3.  **Decimales:** El campo `horas_trabajadas` es un decimal (float64). Puedes usar `.toFixed(2)` para mostrarlo como `8.50`.

# Guía de Integración Frontend: Módulo de Asistencia

Esta documentación detalla los endpoints disponibles para el módulo de **Control de Asistencia** y cómo deben ser consumidos desde el Frontend (Aplicación Móvil o Web). 

> **Nota:** Todos los endpoints requieren que se envíe el JWT del usuario autenticado en el Header `Authorization: Bearer <token>`.

---

## 1. Operaciones del Empleado (App Móvil / Portal del Trabajador)

Estos endpoints están diseñados para ser usados por el propio trabajador.

### 1.1 Marcar Asistencia (Entrada / Salida)
El backend es inteligente y determina si es entrada o salida automáticamente. Con solo presionar el botón "Marcar", el backend registrará la entrada (si es el primer clic del día) o la salida (si es el segundo clic).

*   **Endpoint:** `/api/user/asistencia/marcar`
*   **Método:** `POST`
*   **Body:** `Vacio` (No requiere enviar datos en el cuerpo)

**Ejemplo de respuesta exitosa (200 OK):**
```json
{
  "id": 15,
  "empresa_id": 1,
  "usuario_id": 4,
  "fecha": "2024-05-08T00:00:00Z",
  "hora_entrada": "2024-05-08T08:15:00Z",
  "hora_salida": null,
  "estado": "tarde",
  "notas": null,
  "horas_trabajadas": null
}
```

### 1.2 Ver Mi Historial
Permite al trabajador ver sus marcas de asistencia recientes.

*   **Endpoint:** `/api/user/asistencia/mi-historial`
*   **Método:** `GET`

**Ejemplo de respuesta:** (Retorna un Array de objetos similares al anterior).

### 1.3 Solicitar un Permiso / Justificación
Si el trabajador falta o necesita justificar una tardanza.

*   **Endpoint:** `/api/user/asistencia/permisos`
*   **Método:** `POST`
*   **Body Request:**
```json
{
  "fecha": "2024-05-10", 
  "motivo": "Cita médica en Essalud a las 10:00am"
}
```

---

## 2. Operaciones del Administrador (Panel de Control)

Estos endpoints son para RRHH o los administradores de la inmobiliaria para gestionar al personal.

### 2.1 Listar Asistencia Global
Muestra todos los registros de la empresa con opciones de filtrado.

*   **Endpoint:** `/api/user/asistencia/registros`
*   **Método:** `GET`
*   **Query Params Disponibles:**
    *   `empresa_id` (Obligatorio)
    *   `usuario_id` (Opcional - Filtra por empleado específico)
    *   `desde` (Opcional - Formato `YYYY-MM-DD`)
    *   `hasta` (Opcional - Formato `YYYY-MM-DD`)
    *   `pag` y `limite` (Opcional para paginación)

**Ejemplo de petición URL:**
`/api/user/asistencia/registros?empresa_id=1&desde=2024-05-01&hasta=2024-05-31`

### 2.2 Asignar Horario a un Empleado
Crea o actualiza el horario y nivel de tolerancia de un trabajador específico.

*   **Endpoint:** `/api/user/asistencia/horarios?empresa_id=1`
*   **Método:** `POST`
*   **Body Request:**
```json
{
  "usuario_id": 4,
  "hora_entrada": "08:00",
  "hora_salida": "17:00",
  "tolerancia_minutos": 15,
  "dias_laborables": "1,2,3,4,5" // Lunes a Viernes
}
```

### 2.3 Aprobar o Rechazar Permiso
El administrador revisa una solicitud de permiso y decide su estado.

*   **Endpoint:** `/api/user/asistencia/permisos/:id/estado?empresa_id=1`
*   **Método:** `PUT`
*   **Reemplazar `:id`** en la URL por el ID real del permiso.
*   **Body Request:**
```json
{
  "estado": "aprobado", // "aprobado" o "rechazado"
  "respuesta": "Documentación médica verificada, justificación aceptada."
}
```

---

## 3. Ejemplo de Consumo en JS / React (Axios)

**Botón Inteligente de "Marcar Asistencia" en la App:**

```javascript
import axios from 'axios';

const marcarAsistencia = async () => {
  try {
    const token = localStorage.getItem('token');
    const response = await axios.post('/api/user/asistencia/marcar', {}, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    const registro = response.data;
    if (registro.hora_salida === null) {
      alert(`Entrada registrada: Estado ${registro.estado}`);
    } else {
      alert(`Salida registrada. Trabajaste ${registro.horas_trabajadas.toFixed(2)} horas.`);
    }

  } catch (error) {
    // Si ya marcó entrada y salida el mismo día, retorna un error
    alert("Error: " + (error.response?.data?.message || "No se pudo marcar asistencia"));
  }
};
```

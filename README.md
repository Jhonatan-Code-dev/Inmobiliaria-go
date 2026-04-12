# Inmobiliaria: Sistema de Gestión Inmobiliaria

Bienvenido al repositorio principal de Inmobiliaria. Este sistema es una plataforma completa (SaaS) orientada a la **gestión eficiente de propiedades, contratos de alquiler, finanzas (ingresos y egresos) y control de mantenimiento mentalizado para corredores y empresas de bienes raíces**.

El backend está desarrollado en Go con una arquitectura limpia centrada en el dominio, asegurando una gran capacidad de escalabilidad, seguridad en las transacciones e historial auditable.

---

## Roles del Sistema

El sistema implementa un estricto control de acceso basado en roles (RBAC). Cuenta **únicamente** con los siguientes 4 roles definidos:

1. **Administrador (Admin):** Acceso total al sistema y la configuración de la empresa.
2. **Supervisor:** Supervisión general, monitoreo de métricas y acceso completo a los reportes.
3. **Vendedor:** Gestión operativa, control de clientes, inmuebles y creación de contratos.
4. **Inventario:** Destinado al personal de campo para el control estricto del estado físico de las unidades y sus equipamientos.

---

## Módulos Funcionales y Requerimientos

A continuación, se detalla la tabla de los requerimientos y módulos que componen el ecosistema de este software inmobiliario:

| # | Módulo | Descripción / Funcionalidades Principales |
| :--- | :--- | :--- |
| **1.1** | **Usuarios, Roles y Accesos** | • Autenticación segura (JWT + Cookies HttpOnly).<br>• Gestión de usuarios (staff).<br>• Control y bloqueo de acceso basado en los 4 roles (RBAC). |
| **1.2** | **Gestión de Propiedades** | • Registro de inmuebles padre (Edificios, Complejos).<br>• Gestión de unidades (Departamentos, Oficinas) y cambios de estado (*Disponible, Ocupado, Mantenimiento*). |
| **1.3** | **Gestión de Clientes** | • Registro de inquilinos / clientes.<br>• Mapeo de datos personales, contacto directo y contactos de emergencia. |
| **1.4** | **Contratos de Alquiler** | • Creación de nuevos contratos.<br>• Asociación `Cliente` ↔ `Unidad`.<br>• Control de fechas, montos y depósitos de garantía.<br>• Estados del contrato automáticos (*Activo, Finalizado, Cancelado*).<br>• Cálculo dinámico de moras diarias/mensuales. |
| **1.5** | **Cobranzas y Pagos** *(Ingresos)* | • Emisión de cargos automáticos en fechas de vencimiento y cargos manuales.<br>• Registro auditable de pagos (completos o parciales).<br>• Panel de control de saldos pendientes y morosidad. |
| **1.6** | **Gastos** *(Egresos)* | • Registro de egresos para contabilidad cruzada.<br>• Organización de desembolsos por categorías y métodos de transferencia/pago. |
| **1.7** | **Mantenimiento y Servicios** | • Emisión y gestión de tickets de mantenimiento (*Abierto, En Progreso, Resuelto*).<br>• Control de métricas de servicios variables puntuales (*agua, luz, etc.*). |
| **1.8** | **Inventario de Unidades** | *(Versión simplificada)*<br>• Registro interno de la unidad *(Muebles, Equipamiento, Llaves, Estado general)*.<br>• Control de estado de las cosas *(Bueno, Regular, Dañado)*.<br>• Historial auditable de los cambios/actualizaciones de equipamiento. |
| **1.9** | **Reportes** | • Generación de gráficas y KPIs consolidados.<br>• Reportes financieros (Caja, Ingresos vs Egresos).<br>• Reportes de ocupación y morosidad (Riesgo).<br>• Funcionalidad de Exportación en formatos amigables *(PDF / Excel)*. |

---

*Desarrollado para el manejo corporativo estructurado de rentas y operaciones.*

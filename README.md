# INMOBILIARIA - Sistema de Control de Rentas y Operaciones Inmobiliarias

Este sistema consiste en una plataforma de software backend premium de nivel empresarial, diseñada en Go y estructurada mediante Clean Architecture. Permite gestionar la operación de propiedades en renta, emisión automatizada de contratos de alquiler, control de cobros y morosidad, facturación de consumos variables y mantenimiento preventivo de las unidades, todo bajo un estricto aislamiento de datos corporativos por empresa.

---

## Arquitectura de Seguridad y Cifrado

El backend de la plataforma se ha construido con un enfoque estricto en la seguridad de los datos financieros e inmobiliarios de los inquilinos y propietarios:

1. **Cifrado de Contraseñas**:
   - El almacenamiento de credenciales de usuario se realiza aplicando Bcrypt con un factor de costo seguro (cost = 10 o superior). Ninguna contraseña se almacena en texto plano.
2. **Autenticación REST y Sesiones**:
   - Autenticación stateless de alta velocidad mediante JWT (JSON Web Tokens) firmados con el algoritmo criptográfico de hash simétrico HS256.
   - Los tokens emitidos se configuran de forma segura con tiempos de expiración parametrizables (por defecto 15d en producción) y se distribuyen preferencialmente en las peticiones a través de Cookies con atributos HttpOnly, Secure y SameSite=Lax para mitigar de forma definitiva ataques de tipo Cross-Site Scripting (XSS) y Cross-Site Request Forgery (CSRF).
3. **Seguridad y Aislamiento por Empresa**:
   - El sistema implementa un control de acceso seguro y unificado. El middleware TenantAuth valida el token de sesión y extrae de forma automática el empresa_id asociado.
   - Todas las consultas de lectura, creación, actualización o eliminación en la base de datos a través de Ent ORM incluyen obligatoriamente el filtro contextual de empresa_id, bloqueando cualquier intento de fuga o cruce de datos no autorizado entre diferentes cuentas de empresas.

---

## Roles del Sistema (RBAC)

El sistema implementa un control de acceso basado en roles (Role-Based Access Control) con únicamente los siguientes 3 roles definidos, asegurando que cada perfil interactúe estrictamente con lo que necesita:

1. **Administrador (Admin):** Acceso total al sistema, configuración de pasarelas, contabilidad corporativa, y control y creación de la nómina/personal (staff) de la empresa.
2. **Supervisor:** Monitoreo general de la operación. Acceso exclusivo a auditorías, estadísticas de rendimiento y generación del set de reportes gerenciales para la toma de decisiones.
3. **Vendedor:** Gestión operativa diaria. Creación y edición de inmuebles, registro de clientes/inquilinos, generación de nuevos contratos de alquiler y recepción de pagos.

---

## Módulos Funcionales y Requerimientos del Sistema

Todos los módulos del backend y sus funcionalidades se presentan a continuación en una tabla consolidada y ordenada:

| Código | Módulo | Requerimiento Funcional | Descripción / Funcionalidad Detallada |
| :--- | :--- | :--- | :--- |
| **1.1** | Gestión de Propiedades | Registro de Inmuebles Padre | Permite registrar edificios, condominios o galerías comerciales con sus metadatos de dirección y servicios. |
| **1.2** | Gestión de Propiedades | Control de Unidades de Arriendo | Permite dar de alta departamentos, oficinas o locales comerciales asociados a un inmueble, con transiciones de estado (*Disponible, Ocupado, Mantenimiento*). |
| **2.1** | Clientes e Inquilinos | Registro Centralizado de Inquilinos | Alta y edición de arrendatarios con validaciones de tipo de documento (DNI, RUC, carnet de extranjería). |
| **2.2** | Clientes e Inquilinos | Historial del Cliente | Almacenamiento de referencias laborales, cuentas bancarias, contactos de emergencia y estado del perfil. |
| **3.1** | Contratos de Alquiler | Emisión de Contratos | Vinculación formal de un Cliente con una Unidad de Arriendo, especificando fecha de inicio, término, montos y depósitos en garantía. |
| **3.2** | Contratos de Alquiler | Gestión del Ciclo de Vida | Control de estados del contrato (*Borrador, Activo, Finalizado, Cancelado*). |
| **3.3** | Contratos de Alquiler | Motor de Cálculo de Moras | Cálculo en tiempo real de días transcurridos tras vencimiento de fecha de pago y aplicación automática de penalizaciones. |
| **4.1** | Cobranzas y Pagos | Emisión Automática de Cargos | Generación automática mensual de cobros programados de arriendo. |
| **4.2** | Cobranzas y Pagos | Emisión Manual de Cargos | Posibilidad de aplicar penalizaciones extraordinarias o cargos manuales. |
| **4.3** | Cobranzas y Pagos | Registro de Pagos e Ingresos | Registro de pagos completos o abonos parciales, con fecha de abono y método de pago. |
| **4.4** | Cobranzas y Pagos | Estado de Cuenta | Consulta y generación en tiempo real del saldo pendiente acumulado e historial de transacciones del inquilino. |
| **5.1** | Contabilidad de Gastos | Registro de Egresos | Registro de gastos de administración, mantenimiento general e impuestos corporativos para contabilidad cruzada. |
| **5.2** | Contabilidad de Gastos | Clasificación de Gastos | Vinculación del gasto con categorías específicas e inyección del comprobante físico de pago. |
| **6.1** | Soporte y Mantenimiento | Emisión de Tickets de Soporte | Registro de solicitudes de soporte técnico imputados a unidades específicas. |
| **6.2** | Soporte y Mantenimiento | Flujo del Ticket | Gestión de prioridades (*Baja, Media, Alta*) y estados del ticket (*Abierto, En Progreso, Resuelto, Cerrado*). |
| **6.3** | Soporte y Mantenimiento | Lectura de Medidores | Captura periódica de consumo variable de agua y energía eléctrica para cobro indexado. |
| **7.1** | Reportes Gerenciales | Tendencia Mensual (Balance) | Serie temporal histórica de ingresos versus gastos y balance neto mensual (`GET /api/user/reportes/ingresos-gastos`). |
| **7.2** | Reportes Gerenciales | Distribución por Método de Pago | Porcentaje y volumen de ingresos según método de cobro utilizado (`GET /api/user/reportes/metodos-pago`). |
| **7.3** | Reportes Gerenciales | Distribución por Categoría de Gastos| Clasificación de egresos agrupados por tipo de pago de gasto (`GET /api/user/reportes/categorias-gastos`). |
| **7.4** | Reportes Gerenciales | Rentabilidad por Propiedad | Tasa de ocupación, ingresos y asignación pro-rata de gastos corporativos por inmueble (`GET /api/user/reportes/rentabilidad-propiedades`). |
| **7.5** | Reportes Gerenciales | Carga de Trabajo de Mantenimiento | Desglose cuantitativo de tickets de soporte agrupados por estado y urgencia (`GET /api/user/reportes/tickets-mantenimiento`). |

---

## Instrucciones de Ejecución Local (Desarrollo)

### Requisitos Previos
* Go (Versión 1.20 o superior instalada).
* Base de Datos: PostgreSQL o MySQL disponible.
* Google Wire: Generador de inyección de dependencias.
  ```bash
  go install github.com/google/wire/cmd/wire@latest
  ```

### Pasos para Arrancar
1. Clonar el proyecto y situarse en el directorio raíz.
2. Configurar el archivo `.env`:
   Crea o edita el archivo `.env` en la raíz del proyecto agregando tus variables de entorno correspondientes:
   ```env
   APP_ENV=local
   PORT=4000
   BASE_DATOS_1=postgres://tu_usuario:tu_password@localhost:5432/alquilamax?sslmode=disable&TimeZone=UTC
   JWT_SECRET=TU_SUPER_CLAVE_SECRETA_JWT_DESARROLLO
   JWT_ACCESS_DURATION=15d
   COOKIE_MAX_AGE=15d
   COOKIE_SECURE=false
   ```
3. Generar Contenedor de Inyección de Dependencias (Wire):
   Cada vez que agregues nuevos controladores, servicios o repositorios, compila la DI corriendo:
   ```bash
   go generate ./di
   ```
4. Ejecutar el Servidor:
   El sistema aplicará automáticamente todas las migraciones del esquema de base de datos gracias a Ent ORM al iniciar.
   ```bash
   go run cmd/server/main.go
   ```
   El servidor backend estará listo y escuchando peticiones en: `http://localhost:4000`.

---

## Instrucciones de Despliegue en Producción (Ubuntu / Debian Linux)

El sistema está preparado para un despliegue seguro, automatizado y estable usando Systemd como gestor de servicios del sistema y Caddy como servidor web de proxy inverso con generación automática de certificados SSL (HTTPS).

### Paso 1: Compilar el Binario
Realiza una compilación limpia y optimizada en tu entorno local para la arquitectura de tu servidor de producción (normalmente Linux de 64 bits):

```bash
env GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o alquilamax cmd/server/main.go
```
* `-ldflags="-w -s"` elimina la información de depuración para reducir al máximo el tamaño del binario final.
* Transfiere el archivo compilado al servidor utilizando `scp` o `sftp`.

### Paso 2: Configurar el Servicio del Sistema (Systemd)
Para garantizar que la aplicación se mantenga corriendo en segundo plano de manera persistente y se reinicie automáticamente ante caídas del sistema, utilizaremos un archivo de servicio Systemd.

1. Crea el archivo de servicio correspondiente:
   ```bash
   sudo nano /etc/systemd/system/alquilamax.service
   ```
2. Pega la siguiente estructura (puedes guiarte de la plantilla de `/deploy/service-rentals.sh`):
   ```ini
   [Unit]
   Description=Inmobiliaria Backend Service
   After=network.target

   [Service]
   Type=simple
   User=yona
   WorkingDirectory=/home/yona/rentas/rentals-go
   ExecStart=/home/yona/rentas/rentals-go/alquilamax
   Restart=always
   RestartSec=5s

   # Variables de Entorno para Producción
   Environment="APP_ENV=production"
   Environment="PORT=7000"
   Environment="BASE_DATOS_1=postgres://tu_usuario_prod:tu_password_prod@localhost:5432/alquilamax_prod?sslmode=disable&TimeZone=UTC"
   Environment="JWT_SECRET=ALQUILAMAX_SUPER_SECRET_KEY_2026"
   Environment="JWT_ACCESS_DURATION=15d"
   Environment="COOKIE_MAX_AGE=15d"
   Environment="COOKIE_SECURE=true"

   # Gestión de recursos
   LimitNOFILE=65535

   [Install]
   WantedBy=multi-user.target
   ```
3. Registra el servicio, inicialízalo y configúralo para que arranque automáticamente con el sistema operativo:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable alquilamax
   sudo systemctl start alquilamax
   ```
4. Monitorea los logs de ejecución en tiempo real para verificar que no haya problemas con la base de datos o puertos:
   ```bash
   sudo journalctl -u alquilamax.service -f
   ```

### Paso 3: Configurar el Proxy Inverso (Caddy Server)
Utilizaremos Caddy para direccionar las peticiones entrantes del dominio de manera pública y segura mediante certificados SSL automatizados.

1. Abre el archivo de configuración de tu servidor Caddy:
   ```bash
   sudo nano /etc/caddy/Caddyfile
   ```
2. Agrega tu bloque de configuración asociando tu dominio público al puerto local configurado en producción (siguiendo `/deploy/CADDY.sh`):
   ```caddy
   alquilamax.duckdns.org {
       # Redirección en proxy inverso hacia el backend en el puerto 7000
       reverse_proxy localhost:7000 {
           header_up Host {host}
           header_up X-Real-IP {remote_host}
           header_up X-Forwarded-For {remote_host}
           header_up X-Forwarded-Proto {scheme}
       }
       
       # Habilitar compresión de respuesta para mejorar velocidad
       encode gzip zstd
   }
   ```
3. Valida la sintaxis de tu configuración y recarga el servidor Caddy sin caídas en caliente:
   ```bash
   caddy validate --config /etc/caddy/Caddyfile
   sudo systemctl reload caddy
   ```

Con estos pasos, tu backend corporativo estará corriendo de manera robusta y segura en producción bajo una conexión HTTPS encriptada.

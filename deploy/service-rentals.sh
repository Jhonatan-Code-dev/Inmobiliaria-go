# ==============================================================================
# PASOS PARA DESPLEGAR EL SERVICIO ALQUILAMAX EN LINUX (ALMALINUX)
# ==============================================================================

# 1. Crear la carpeta de la aplicación
sudo mkdir -p /opt/alquilamax

# 2. Asegurar permisos de la carpeta (root como dueño)
sudo chown root:root /opt/alquilamax

# 3. Mover el binario compilado a la carpeta destino
# Nota: Asegúrate de haber compilado el binario antes (ver comando al final)
sudo mv /home/opc/alquilamax-bin /opt/alquilamax/

# 4. Dar permisos de ejecución al binario
sudo chmod +x /opt/alquilamax/alquilamax-bin

# 5. Corregir etiquetas de seguridad SELinux (Crucial para AlmaLinux/RHEL)
sudo restorecon -Rv /opt/alquilamax/

# 6. Crear/Editar el archivo de servicio de systemd
sudo vi /etc/systemd/system/alquilamax.service

/*
[Unit]
Description=Servicio Backend Alquilamax
After=network.target postgresql-17.service

[Service]
User=opc
Group=opc
WorkingDirectory=/opt/alquilamax
ExecStart=/opt/alquilamax/alquilamax-bin
Restart=always
RestartSec=5

# Variables de Entorno para Producción
Environment="APP_ENV=production"
Environment="PORT=7000"
Environment="BASE_DATOS_1="
Environment="JWT_SECRET=ALQUILAMAX_SUPER_SECRET_KEY_2026"
Environment="JWT_ACCESS_DURATION=15d"
Environment="COOKIE_MAX_AGE=15d"
Environment="COOKIE_SECURE=true"
Environment="ALLOWED_ORIGINS=*"

[Install]
WantedBy=multi-user.target
*/

# ==============================================================================
# COMANDOS DE GESTIÓN DEL SERVICIO
# ==============================================================================

# Recargar configuración de systemd
sudo systemctl daemon-reload

# Habilitar y arrancar el servicio
sudo systemctl enable alquilamax
sudo systemctl start alquilamax
sudo systemctl stop alquilamax
# Verificar estado
sudo systemctl status alquilamax

# Ver logs en tiempo real (Para depuración)
sudo journalctl -u alquilamax.service -f

# ==============================================================================
# PROCEDIMIENTO DE ACTUALIZACIÓN (RE-DEPLOY)
# ==============================================================================
# Solo ejecuta esto cuando subas un nuevo binario:

sudo systemctl stop alquilamax && \
sudo mv /home/opc/alquilamax-bin /opt/alquilamax/alquilamax-bin && \
sudo chmod +x /opt/alquilamax/alquilamax-bin && \
sudo restorecon -v /opt/alquilamax/alquilamax-bin && \
sudo systemctl start alquilamax && \
sudo systemctl status alquilamax

# ==============================================================================
# COMANDO PARA COMPILAR (EJECUTAR EN TU MÁQUINA LOCAL)
# ==============================================================================
GOOS=linux GOARCH=amd64 go build -o alquilamax-bin ./cmd/api/main.go



#-------------------------------------------------------
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags="-w -s" -o alquilamax-bin ./cmd/api/main.go

#------------------------------------------------------
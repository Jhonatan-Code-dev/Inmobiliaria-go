
GOOS=linux GOARCH=amd64 go build -o binario-linux ./cmd/api/main.go

🚀 1. COMANDOS (SIN CREAR SERVICE AUTOMÁTICO)
sudo mkdir -p /opt/alquilamax && \
sudo chown -R opc:opc /opt/alquilamax && \
sudo mv /home/opc/binario-linux /opt/alquilamax/binario-linux && \
sudo chmod +x /opt/alquilamax/binario-linux && \
sudo restorecon -Rv /opt/alquilamax/

sudo vi /etc/systemd/system/alquilamax.service

[Unit]
Description=AlquilaMax API
After=network.target postgresql-17.service

[Service]
User=opc
Group=opc
WorkingDirectory=/opt/alquilamax
ExecStart=/opt/alquilamax/binario-linux
Restart=always
RestartSec=5

# ===========================
# VARIABLES DE ENTORNO
# ===========================

# Base de datos
Environment="BASE_DATOS_1=postgres://yona:912059555Perez@localhost:5432/alquilamax?sslmode=disable&TimeZone=UTC"

# JWT
Environment="JWT_SECRET=PAGA_CAUSA_SUPER_SECRET_KEY_2026"
Environment="JWT_ACCESS_DURATION=15d"

# Cookies
Environment="COOKIE_MAX_AGE=15d"
Environment="COOKIE_SECURE=false"

# CORS
Environment="ALLOWED_ORIGINS=*"

# App
Environment="APP_ENV=production"
Environment="PORT=5000"

[Install]
WantedBy=multi-user.target

sudo systemctl daemon-reload && \
sudo systemctl enable alquilamax && \
sudo systemctl start alquilamax && \
sudo systemctl status alquilamax



🔁 6. ACTUALIZAR BINARIO



sudo systemctl stop alquilamax && \
sudo mv /home/opc/binario-linux /opt/alquilamax/binario-linux && \
sudo chmod +x /opt/alquilamax/binario-linux && \
sudo restorecon -v /opt/alquilamax/binario-linux && \
sudo systemctl start alquilamax && \
sudo systemctl status alquilamax
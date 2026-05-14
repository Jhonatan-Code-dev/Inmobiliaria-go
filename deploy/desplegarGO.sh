sudo vi /etc/systemd/system/paga-causa.service

[Unit]
Description=Paga Causa Go App
After=network.target

[Service]
User=opc
Group=opc
WorkingDirectory=/opt/paga-causa
ExecStart=/opt/paga-causa/binario-linux
Restart=always
RestartSec=5

# ===========================
# Variables de entorno
# ===========================



Environment="BASE_DATOS_1=postgres://neondb_owner:npg_gYn0Ktyia3jf@ep-noisy-recipe-ah7cyhd8-pooler.c-3.us-east-1.aws.neon.tech/cuotamax?sslmode=require&channel_binding=require"
Environment="JWT_SECRET=PAGA_CAUSA_SUPER_SECRET_KEY_2026"
Environment="JWT_ACCESS_DURATION=15d"
Environment="COOKIE_MAX_AGE=15d"
Environment="COOKIE_SECURE=true"
Environment="VALKEY_ADDR=192.168.18.223:6379"
Environment="VALKEY_AUTH_TTL=15d"

Environment="OCI_USER_OCID=ocid1.user.oc1..aaaaaaaab7obdj34n7gcgpwm74upxnbkyft6luppyf6hemdmrfvrfirouktq"
Environment="OCI_TENANCY_OCID=ocid1.tenancy.oc1..aaaaaaaayekj5vleum2lov5a2uiewlvucfndlfoqyvqelvn27dnzeirn6gpa"
Environment="OCI_FINGERPRINT=52:be:e9:d5:d4:2e:99:37:1b:8c:ac:2b:84:6e:54:cb"
Environment="OCI_REGION=sa-saopaulo-1"
Environment="OCI_PRIVATE_KEY=-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCW3s39vCYAgxVd\nrHObGdhL7lYWs59SFlYgazLwMjGuI0i+cQM3ECQLWjtjuzUjvMOc0TGva0yQoAUa\nkQyrJbH/J+FA8Wuewb3elWk8DrLWQK5SfUsxszkqpF6WpRnY4A9tbcAuV8LGOBGv\nQ2L5MGg58nOULq5TZSmJUtwjGfmMB5UHCXByHjKRyAjFe18MFQaSDkkzlbvokTZI\nC6GwHfUeEGclzU74S7T4uJsRV6TFoFmjJM81NR32yegEuUw8RS4Tdq2ZiE6rLEmA\nDamXK07rxarJ46GK2dI9U8qMQazYzxlqGJGwi9jVByErBRGXLxR/028olo6S84N0\nqR/xdGrJAgMBAAECggEABTcyYJhTqX6em1iJ/hWKhVR52CrRQg+Y3mHKVuSrVEFu\nwzupsGp/zn8p8LCQ/Zlp20bdnHp9gP3xMzoKubNxx/f/WOUHswswVsysQMniy5ac\nmCgx9GfFXQil5vgR8M4NJnut8kQxLPRVAy21kxcF8WPk9D1ZZDNp+E/7oYCFQJ6Y\nXy5rNWWk1A4vxrWvJRGuqE7hcZDJzi/I9T+pNleF7lIYAAxNWjA/P/XJ8etKPsHA\nAf9mmL1998vTTBJFBGtMd3rGLweluKu5BJOMIcBa4+fOQW4lbaovIggRtlm8Ph8G\nqlEQABRcz5m+wWWkhPUSK295a6jExYpGEGrulgF8wQKBgQDQWt8gxZs8Lk5i2QX4\nrhlehdErp/QH/w8anCxlKUVI5VmBBi24aqm7zQNvZoRm9mSRQOrbHrwA6PXRO6JD\nG9dLHe36ZWPD5Az1DJHwmKma9uLyNHgizG5jIes/rMDom/7UD1jemW91ZKEaPHV1\n38WINxxUX6U+jWCXyLJGnpiXwQKBgQC5XsZISyup9/XnMFuY7mEnt4vNAGXfljjb\nq81U2SKfyB0xKcXRvATvGuVGJAMPs/l34FnkQLj1ljwIOJeusl3TgQ2hz447nlL1\n5glrs40kQaRt9AxxBbvdlwc6VQkuWrK4Fqirbt2TnuOM6Vi3ma9VytpYtEWSJPPe\n+xLxALRVCQKBgAEy1J/Coz74YTkOWItyrPCvQmHG6I93NyYHCfZXA3AE6bvlRjQO\nYQWUi1WDuHVDK5buUauLBLfYnzlh53ANY/KprGnJVYaV9EEgnmJM7oTWsL2F8b75\ngBUP5+OI4d80roWXxQIazdpWBts4x9AyxlcfQgl2N2QhuhGFdQkU9nnBAoGAV1Il\n5XhDoVWFKNrGy0u/yI3V9UPyuVhygEh1+Tov7US/O6GJ5jrDuD7bMidUqdkF80pk\nDfnPJyEWNmkySsELIc7xNQAo6Dy0p9EtLubt3d5uLr1//t3MmZ3Dcd8M7CEdf0pt\nOhSHnDqExqRFfneO+MMOCCsjeqydlLSBP7YtPMECgYEAqCrJBu5FwfO8lRKXPbON\nORPNZWMzWQwpZ2akf/mcWmQCqr2kmCJ1zLRWpgOMz2A42/jdYXDQui4cKanI/vMl\ndL8RRuxetgXO4zYj7l+cCNP5LPYGM0HApP6IGwVgFCljLPJyEiTPgPSop7JCjN16\nIQeWVGaLNRZNKyKyDX9dXDc=\n-----END PRIVATE KEY-----\nOCI_API_KEY"
Environment="OCI_NAMESPACE=grny97yseatg"
Environment="OCI_BUCKET_NAME=imagenes-al-dia"

[Install]
WantedBy=multi-user.target




root:912059555Perez@tcp(localhost:3306)/pagacausa?charset=utf8mb4&parseTime=True&loc=UTC
# importante:
permisos de ejecucion:

sudo chmod +x /home/opc/binario-linux


1. Limpiar el contexto de seguridad (Paso clave)
Ejecuta esto para decirle a SELinux que este archivo es un ejecutable legítimo del sistema:

Bash
sudo restorecon -Rv /opt/paga-causa/binario-linux




# 1. Asegurarse de que existe la carpeta destino
sudo mkdir -p /opt/paga-causa
sudo chown opc:opc /opt/paga-causa

# 2. Mover el binario nuevo (sobrescribe el anterior), dar permisos y reiniciar
sudo mv /home/opc/binario-linux /opt/paga-causa/ && \
sudo chmod +x /opt/paga-causa/binario-linux && \
sudo systemctl daemon-reload && \
sudo systemctl restart paga-causa && \
sudo systemctl status paga-causa

para actualizar el archivo:

# Sube el archivo y corre esto:
sudo mv /home/opc/binario-linux /opt/paga-causa/
sudo restorecon -v /opt/paga-causa/binario-linux
sudo chmod +x /opt/paga-causa/binario-linux
sudo systemctl restart paga-causa







sudo systemctl daemon-reload
sudo systemctl enable paga-causa
sudo systemctl start paga-causa
sudo systemctl status paga-causa


sudo systemctl stop paga-causa



ver log:

journalctl -u paga-causa -f


sudo vi /etc/caddy/Caddyfile

# Redirige HTTP a HTTPS
eventregistry.duckdns.org {
    redir https://{host}{uri} 301
}

# HTTPS con proxy inverso
https://eventregistry.duckdns.org {
    reverse_proxy localhost:4000 {
        header_up Host {host}
        header_up X-Forwarded-Proto {scheme}
    }
}

sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload


sudo systemctl restart caddy
sudo systemctl enable caddy
sudo systemctl status caddy



-----------------------------

GOOS=linux GOARCH=amd64 go build -o binario-linux ./cmd/api/main.go

GOOS=linux GOARCH=amd64 go build -o binario-linux ./cmd/api/main.go
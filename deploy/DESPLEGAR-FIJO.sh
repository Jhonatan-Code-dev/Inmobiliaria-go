sudo systemctl stop paga-causa && chmod +x binario-linux && sudo mv binario-linux /opt/paga-causa/binario-linux && sudo restorecon -v /opt/paga-causa/binario-linux && sudo systemctl start paga-causa && sudo systemctl status paga-causa


EDITAR SERVICE:
sudo vi /etc/systemd/system/paga-causa.service
IMPORTATE VER SERVICIO Y ACTUALIZAR EL SERVICIO:

sudo systemctl daemon-reload
sudo systemctl restart paga-causa.service
sudo systemctl status paga-causa.service


# EL PUERTO DEL API SE DEFINE DESDE EL ARCHIVO .env
# Ejemplo:
# PORT=4000



# COMANDO PARA COPIALAR CODIGO EN LOCAL PARA PRODUCCION:

GOOS=linux GOARCH=amd64 go build -o binario-linux ./cmd/api/main.go

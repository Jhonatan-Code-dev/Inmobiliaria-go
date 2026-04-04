# RESUMEN GENERAL DE RAM
free -h
# VER ESPACIO EN DISCO
df -h
# VER SERVICIOS ACTIVOS
systemctl list-unit-files --type=service
# VER SERVICIOS ACTIVOS AHORA
systemctl list-units --type=service
# INFORMACION COMPLETA DEL PROCESADOR
lscpu
#--------------------------------------------------

COMANDOS PARA VER SERVICIOS DETENERLO , DESABILITARLOS O ELIMINARLOS

#-------------------------------

VER SERVICIOS EN EJECUCION
systemctl list-units --type=service --state=running

✅ PASO 2 — Detener un servicio (solo ahora)
sudo systemctl stop NOMBRE


Ejemplo:

sudo systemctl stop cockpit

✅ PASO 3 — Deshabilitar un servicio (no se levanta al reiniciar)
sudo systemctl disable NOMBRE


Ejemplo:

sudo systemctl disable cockpit

☠️ PASO 4 — Detener + Deshabilitar juntos
sudo systemctl stop NOMBRE
sudo systemctl disable NOMBRE

🧹 PASO 5 — Eliminar el servicio si no lo necesitas

Esto solo aplica si el servicio viene de un paquete instalable.

Primero identifica el paquete al que pertenece el servicio:

rpm -qf $(which NOMBRE)


Si eso no funciona:

rpm -qa | grep -i NOMBRE


Luego eliminas el paquete:

sudo dnf remove PAQUETE


Ejemplo:

sudo dnf remove cockpit


# ZRAM CONFIGURANDO EN VPS ME ENTIENDES
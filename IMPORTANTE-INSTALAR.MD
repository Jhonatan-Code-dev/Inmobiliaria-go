sudo apt update
sudo apt install -y \
  libvips-dev \
  pkg-config \
  build-essential

# INSTALACION PARA BASE DE DATOS ORACLE PARA EL SERVIDOR:

https://www.oracle.com/database/technologies/instant-client/linux-x86-64-downloads.html

https://odpi-c.readthedocs.io/en/latest/user_guide/installation.html
comandos de instalacion:

sudo apt-get update
sudo apt-get install -y libaio1
sudo mkdir -p /opt/oracle

# Crear el archivo de configuración con la ruta exacta
sudo sh -c "echo /opt/oracle/instantclient_19_25 > /etc/ld.so.conf.d/oracle-instantclient.conf"

# Aplicar los cambios en el sistema
sudo ldconfig
ls -l /opt/oracle/instantclient_19_25/libclntsh.so*

#-------INSTALACION DE CADDY----------------
sudo dnf install 'dnf-command(copr)'
sudo dnf copr enable @caddy/caddy -y
sudo dnf install caddy -y
sudo setcap cap_net_bind_service=+ep $(which caddy)

sudo systemctl enable caddy
sudo systemctl start caddy
sudo systemctl status caddy

sudo vi /etc/caddy/Caddyfile

# IMPORTANTE ELIMINAR LA PARTE DEL PUERTO QUE ES :80 MUY IMPORTANTE 
#----REINICIAR CADDY

sudo systemctl reload caddy
sudo systemctl restart caddy


-------------PARA DETENER CADDY----------------------
sudo systemctl stop caddy
sudo systemctl disable caddy
-----------------------------------------------------

-------------PARA ACTIVAR CADDY----------------------
sudo systemctl enable caddy
sudo systemctl start caddy
sudo systemctl restart caddy
----------------------------------------------------



alquilamax.duckdns.org {
    # Reverse proxy hacia tu aplicación en el puerto 4000
    reverse_proxy localhost:4000 {
        header_up Host {host}
    }
}

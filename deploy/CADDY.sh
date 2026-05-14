
#-------INSTALACION DE CADDY----------------
sudo dnf install 'dnf-command(copr)'
sudo dnf copr enable @caddy/caddy -y
sudo dnf install caddy -y
sudo setcap cap_net_bind_service=+ep $(which caddy)

sudo systemctl enable caddy
sudo systemctl start caddy
sudo systemctl status caddy

sudo vi /etc/caddy/Caddyfile


eventregistry.online {
    reverse_proxy localhost:8080
}


eventregistry.online {
    # Reescribe la raíz para que apunte internamente a /apk1/
    rewrite * /demo-1.0-SNAPSHOT{path}

    # Proxy inverso hacia el backend
    reverse_proxy localhost:8080 {
        header_up Host {host}
        header_up X-Forwarded-Prefix /demo-1.0-SNAPSHOT
    }
}


eventregistry.online {
    reverse_proxy /academia/* localhost:8080 {
        header_up Host {host}
    }
}




eventregistry.online {
    reverse_proxy localhost:8080 {
        header_up Host {host}
        header_up X-Forwarded-Prefix /apk1
    }
}


#estes es el de sirve los demas no


eventregistry.online {
    # Redirige la raíz (/) a /apk1/ (código 302 temporal o 301 permanente)
    redir / /academia/ 302

    # Proxy inverso para la ruta /apk1/
    reverse_proxy /academia/* localhost:8080 {
        header_up Host {host}
        header_up X-Forwarded-Prefix /academia
    }
}

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



cuota.duckdns.org {
    # Reverse proxy hacia tu aplicación en el puerto 4000
    reverse_proxy localhost:4000 {
        header_up Host {host}
    }
}

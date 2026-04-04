# Se requiere redis

Redis 7.2.11

# buscando alternativas encontre una mejor que redis

# Valkey

Name : valkey
Version : 8.0.6
Release : 2.el9_7
Architecture : x86_64
Size : 1.6 M
Source : valkey-8.0.6-2.el9_7.src.rpm
Repository : appstream
Summary : A persistent key-value database
URL : https://valkey.io
License : BSD-3-Clause AND BSD-2-Clause AND MIT AND BSL-1.0

# para instalar valkey en docker:

https://valkey.io/download/releases/

docker pull valkey/valkey:8.0.6

docker run -d --name valkey -p 6379:6379 --restart unless-stopped valkey/valkey:8.0.6

docker exec -it valkey valkey-server --version

docker exec -it valkey valkey-cli

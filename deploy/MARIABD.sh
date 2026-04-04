sudo dnf module list mariadb

sudo vi /etc/yum.repos.d/MariaDB.repo

#---------------------------------------------------
[mariadb]
name = MariaDB
baseurl = https://rpm.mariadb.org/10.11/rhel9-amd64
gpgkey=https://rpm.mariadb.org/RPM-GPG-KEY-MariaDB
gpgcheck=1
enabled=1
module_hotfixes=1
#---------------------------------------------------

sudo dnf install MariaDB-server -y

sudo systemctl enable mariadb
sudo systemctl start mariadb
sudo systemctl status mariadb


sudo mysql -u root -e "
GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION;
FLUSH PRIVILEGES;
SET PASSWORD FOR 'root'@'localhost' = PASSWORD('912059555Perez');
DROP DATABASE IF EXISTS test;
DELETE FROM mysql.db WHERE Db = 'test' OR Db = 'test_%';
FLUSH PRIVILEGES;
"


sudo systemctl restart mariadb

#reiniciar la maquina
sudo reboot

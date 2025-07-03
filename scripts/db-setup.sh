apt-get update
apt-get install -y mariadb-server

cat <<EOF > /etc/mysql/mariadb.conf.d/99-vagrant.cnf
[mysqld]
bind-address = 0.0.0.0
EOF

service mysql restart

[ -f /vagrant/db/db_init.sql ] && \
  mysql -u root < /vagrant/db/db_init.sql


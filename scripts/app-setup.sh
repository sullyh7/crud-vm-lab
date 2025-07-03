set -e

sudo add-apt-repository ppa:longsleep/golang-backports
apt-get update
apt-get install -y golang-go make

# based on the db ip in the private network (might fetch from .env file instead)
export DB_ADDR="todo:123@tcp(192.168.33.10:3306)/todoapp"

cd /vagrant/app
mkdir -p bin logs
make build

nohup ./bin/todoapp > logs/app.log 2>&1 &
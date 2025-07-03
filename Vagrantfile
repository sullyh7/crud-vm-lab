Vagrant.configure("2") do |config|
  config.vm.define "db" do |db|
    db.vbguest.auto_update = false
    db.vbguest.no_remote = true

    db.vm.hostname = "db.vm"
    db.vm.network "private_network", ip: "192.168.33.10"

    db.vm.synced_folder "./db", "/vagrant/db"
    db.vm.box = "bento/ubuntu-22.04"
    db.vm.box_version = "202502.21.0"
    db.vm.provision "shell", path: "scripts/db-setup.sh"
  end
  config.vm.define "app" do |app|
    app.vbguest.auto_update = false
    
    app.vbguest.no_remote = true

    app.vm.hostname = "app.vm"
    app.vm.network "private_network", ip: "192.168.33.11"

    app.vm.synced_folder "./app", "/vagrant/app"
    app.vm.synced_folder "./scripts", "/vagrant/scripts"

    app.vm.box = "bento/ubuntu-22.04"
    app.vm.box_version = "202502.21.0"
    app.vm.provision "shell", path: "scripts/app-setup.sh"
  end
end



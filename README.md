# VirtualBox & Vagrant Home Lab

This repository demonstrates how to set up a fully reproducible home lab using VirtualBox and Vagrant. Every VM, network, and service is defined in code—no manual steps required beyond installing VirtualBox and Vagrant.

* **Vagrantfile**: Defines two VMs (`db` and `app`) on a private network with fixed IPs.
* **db/db\_init.sql**: SQL script to create database, user, and schema.
* **app/**: Go CRUD service source and module files.
* **scripts/**: Shell provisioners for the DB and App VMs.

---

## 🔧 Prerequisites

* **VirtualBox**
* **Vagrant** 

---

## 🚀 Getting Started

1. **Clone this repo**

   ```bash
   git clone https://github.com/yourusername/vagrant-home-lab.git
   cd vagrant-home-lab
   ```

2. **(Optional) Configure environment variables**

   The project creates any environment variables within the provision scripts, you can also make this external in an env file if needed.

3. **Spin up the lab**

   ```bash
   vagrant up
   ```

   Vagrant will:

   * Launch two Ubuntu 22.04 VMs (`db.vm` & `app.vm`) via VirtualBox
   * Run `scripts/db-setup.sh` on `db.vm` to install MariaDB, configure networking, and apply `db/db_init.sql`
   * Run `scripts/app-setup.sh` on `app.vm` to install Go, build the service, and launch it

4. **Access the VMs**

   ```bash
   vagrant ssh db   # Logs into the database VM
   vagrant ssh app  # Logs into the app VM
   ```

5. **Verify the service**

   On your host machine, in a browser or via curl:

   ```bash
   curl http://(app-ip):8080/api/todos
   ```

   You can also use localhost if you add port forward options from the host machine to the app VM.

---

## 🛠️ Common Commands

* **Halt the lab**

  ```bash
  vagrant halt
  ```

* **Destroy all resources**

  ```bash
  vagrant destroy -f
  ```

* **Reprovision** (re-run scripts)

  ```bash
  vagrant reload --provision
  ```

---

## 🔍 How It Works

* **Infrastructure as Code**: `Vagrantfile` declares two VMs with fixed IPs on a private network.
* **Private Networking**: Host can interact with both VMs; the VMs can communicate with each other directly and with the internet using their NAT network adapter at `192.168.33.10` and `192.168.33.11`.
* **Shell Provisioners**: `scripts/db-setup.sh` and `scripts/app-setup.sh` perform idempotent installs, configuration, and service startup.
* **Go Service**: The `app` VM builds and runs a simple CRUD API (`/api/todos` endpoint) backed by MySQL/MariaDB.

---

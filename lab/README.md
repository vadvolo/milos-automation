# Network Automation Lab

## Installation

### Preparation

You need to install Docker and Docker Compose softaware on yours device:

- [Docker](https://docs.docker.com/engine/install/)
  - [Linux](https://docs.docker.com/desktop/install/linux/)
  - [Mac](https://docs.docker.com/desktop/install/mac-install/)
  - [Windows](https://docs.docker.com/desktop/install/windows-install/)

Pull lab environment from repository:

```
git clone https://github.com/vadvolo/milos-automation.git
cd milos-automation
```

And switch to the `temp` branch:

```
git checkout temp
```

Navigate to the lab folder:

```
cd lab
```

Build and run lab

```
make build
make run
```

To perform import data into netbox please run `make netbox_import`
```
❯ make netbox_import
🧬 loaded config '/etc/netbox/config/configuration.py'
🧬 loaded config '/etc/netbox/config/extra.py'
🧬 loaded config '/etc/netbox/config/logging.py'
🧬 loaded config '/etc/netbox/config/plugins.py'
Installed 742 object(s) from 1 fixture(s)
```

### How to connect to containers

Netbox: `docker exec -u root -t -i netbox-docker-netbox-1 /bin/bash`
Annet: `docker exec -u root -t -i netbox-docker-annet-1 /bin/bash`
Dynamips: `docker exec -u root -t -i netbox-docker-dynamips-lab-1 /bin/bash`

### Create Netbox SuperUser [optional]

1. Go to the netbox container: `docker exec -u root -t -i netbox-docker-netbox-1 /bin/bash`
2. Run command: `/opt/netbox/netbox/manage.py createsuperuser`
3. Input username, email and password:

```
Username (leave blank to use 'root'): milos
Email address: milos@nh.com
Password: milos
Password (again): milos
Superuser created successfully.
```

### NETBOX TOKEN [optional]

1. Go to [LOCAL NETBOX INSTALLATION](http://localhost:8000/users/tokens/)
2. Push `+ Add` button
3. Generate token
4. Copy token
5. exit from netbox container
6. `export NETBOX_TOKEN=a630dcef...`
7. `make annet_restart`

## Lab management

Lab topologies are managed by [Containerlab](https://containerlab.dev/), which is itself deployed in a container. Executable `lab/containerlab` is a simple wrapper for accessing Containerlab container. Available topologies can be found in `lab/topologies`. To spin up the topology:
```bash
./containerlab deploy --topo topologies/lab01_frr-only-test/frr-only-test.clab.yml
```

Connect to a node:
```bash
docker exec -it clab-frr-only-test-frr-r2 bash
```

Destroy the topology:
```bash
./containerlab destroy --topo topologies/lab01_frr-only-test/frr-only-test.clab.yml
```

Full list of commands can be found [here](https://containerlab.dev/cmd/deploy/).

## How to use

Go to annet:

```
docker exec -u root -t -i netbox-docker-annet-1 /bin/bash
```

Run:

- diff

```
python3 -m annet.annet diff lab-r1.nh.com
```

- patch

```
python3 -m annet.annet patch lab-r1.nh.com
```

- deploy

```
python3 -m annet.annet deploy lab-r1.nh.com
```

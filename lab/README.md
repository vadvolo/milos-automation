# Network Automation Lab

## General information

In this repo are located a virtual lab equipment for demonstrating network automation scenarious.

Table of the content:
- [Installation](#installation)
- [Annet description](#annet-description)
- [Useful commands](#useful-commands)
- [Labs](#labs)
  - [lab00. Cisco Base Scenario](./topologies/lab00_cisco_base)
  - [lab01. FRR Base Scenario](./topologies/lab01_frr-only-test)
  - [lab02. Arista Base Scenario](./topologies/lab02_ceos-test)
  - [lab03. Cisco DC Scenario](./topologies/lab03_dc_cisco)

## Installation

### Preparation

First of all, you need to install Docker and Docker Compose softaware on yours device:

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
```

Now you can choose which scenaro you want to run. To start lab you need to run a command `make labXX`, where `XX` is index of the lab. For instance `make lab00`

## Annet description

Annet is an solution for the network configuration management. It provides capabilities for storing templates of configuration, generating and implimentation network configuration. There are four main command for it:
- `gen`
- `diff`
- `patch`
- `deploy`

Annet uses Netbox as Source of Truth about network topology, equipment and resources.

To use annet at presented labs you should prepare generators located at `annet/my_generators`.

### How to use Annet

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

## Useful commands

### How to connect to containers

Netbox: `docker exec -u root -t -i netbox-docker-netbox-1 /bin/bash`
Annet: `docker exec -u root -t -i netbox-docker-annet-1 /bin/bash`
Dynamips: `docker exec -u root -t -i netbox-docker-dynamips-lab-1 /bin/bash`

### How to create Netbox SuperUser [optional]

You don't need to generate your own token if you've made `make netbox_import`.

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

### How to create Netbox Token [optional]

You don't need to generate your own token if you've made `make netbox_import`.

1. Go to [LOCAL NETBOX INSTALLATION](http://localhost:8000/users/tokens/)
2. Push `+ Add` button
3. Generate token
4. Copy token
5. exit from netbox container
6. `export NETBOX_TOKEN=a630dcef...`
7. `make annet_restart`

## Labs

| Lab | Description |
|:---:|:-----------:|
| lab00 | Basic Scenarion with Cisco devices. |


## Lab management

Lab topologies are described by `docker-compose.yml` files in `lab/topologies` folder. Topologies can be deployed by `make lab*_up` and destroyed by `make lab*_down`:
```bash
make lab01_up
make lab01_down
```

Configuration files for nodes in the topology are separated into immutable default configurations `lab/topologies/lab*/default_configs`, used when initializing the lab, and operational configurations `lab/topologies/lab*/operational_configs`, that can be changed by the nodes. Operational configurations are created on lab startup from the defaults, and they are not tracked by Git.

### Lab 01: FRR-only test

To connect to lab nodes, one can use `docker exec` or SSH:
```bash
docker exec -it frr-r1 bash # open shell
# or
docker exec -it frr-r1 vtysh # open FRR CLI
# or
ssh root@172.20.0.110 # password: frr
```

Node list:
- frr-r1: 172.20.0.110
- frr-r2: 172.20.0.111
- frr-r3: 172.20.0.112

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

## Annet description

Annet is an solution for the network configuration management. It provides capabilities for storing templates of configuration, generating and implimentation network configuration. There are four main command for it:
- `gen`
- `diff`
- `patch`
- `deploy`

Annet uses Netbox as Source of Truth about network topology, equipment and resources.

To use annet at presented labs you should prepare generators located at `annet/my_generators`.

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

Lab topologies are managed by [Containerlab](https://containerlab.dev/), which is itself deployed in a container. Executable `lab/containerlab` is a simple wrapper for accessing Containerlab container. Available topologies can be found in `lab/topologies`. To spin up the topology:
```bash
./containerlab deploy --topo topologies/lab01_frr-only-test/frr-only-test.clab.yml
```

Connect to an FRR node:
```bash
docker exec -it clab-frr-only-test-frr-r2 bash
```

Connect to an Arista cEOS node:
```bash
docker exec -it clab-ceos-test-ceos-r3 Cli
  <OR>
ssh admin@clab-ceos-test-ceos-r3
```

Destroy the topology:
```bash
./containerlab destroy --topo topologies/lab01_frr-only-test/frr-only-test.clab.yml
```

Full list of commands can be found [here](https://containerlab.dev/cmd/deploy/).

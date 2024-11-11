# Network Automation Lab

## General information

In this repo are located a virtual lab equipment for demonstrating network automation scenarious.

Table of the content:
- [Installation](#installation)
- [Annet description](#annet-description)
- [Useful commands](#useful-commands)
- [Labs](#labs)
  - [lab00. Cisco Base Scenario](./topologies/lab00_cisco_basic_scenario)
  - [lab01. FRR Base Scenario](./topologies/lab01_frr-only-test)
  - [lab02. Cisco DC Scenario](./topologies/lab02_dc_cisco)
  - [lab03. Multivendor Lab](./topologies/lab03_multivendor)

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

[Annet](https://annetutil.github.io/annet/main/index.html) is a system of network configuration management. It provides capabilities for storing templates of configuration, generating and implimentation network configuration. There are four main command for it:
- `gen`
- `diff`
- `patch`
- `deploy`

To use annet at presented labs you should prepare generators located in the `topologies/lab*/src/lab_generators/` folder.

### How to use Annet

Go to annet:

```
docker exec -u root -t -i netbox-docker-annet-1 /bin/bash
```

Run:

- diff

```
annet diff lab-r1.nh.com
```

- patch

```
annet patch lab-r1.nh.com
```

- deploy

```
annet deploy lab-r1.nh.com
```

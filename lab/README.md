# Network Automation Lab

## General information

This repo contains virtual lab equipment for demonstrating network automation scenarios.

Table of contents:

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

First of all, you need to install Docker and Docker Compose on your device:

- [Docker](https://docs.docker.com/engine/install/)
  - [Linux](https://docs.docker.com/desktop/install/linux/)
  - [Mac](https://docs.docker.com/desktop/install/mac-install/)
  - [Windows](https://docs.docker.com/desktop/install/windows-install/)

Clone this repository:

```bash
git clone https://github.com/vadvolo/milos-automation.git
```

Navigate to the lab folder:

```bash
cd milos-automation/lab
```

Build Annet and Netbox Docker images:

```bash
make build
```

Now you can choose which scenario you want to run. To start a lab you need to run `make labXX`, where `XX` is an index of the lab.
For example, `make lab00` will start `lab00. Cisco Base Scenario`.

## Annet description

[Annet](https://annetutil.github.io/annet/main/index.html) is a network configuration management system. It provides capabilities for storing configuration templates, generating and deploying network configurations.
Annet has four main arguments:

- `gen`
- `diff`
- `patch`
- `deploy`

To use annet in the lab topologies presented in this repository, you should prepare generators located in the `topologies/lab*/src/lab_generators/` folder.

### How to use Annet

Go to Annet container:

```bash
docker exec -u root -t -i annet /bin/bash
```

Run:

- diff

```bash
annet diff lab-r1.nh.com
```

- patch

```bash
annet patch lab-r1.nh.com
```

- deploy

```bash
annet deploy lab-r1.nh.com
```

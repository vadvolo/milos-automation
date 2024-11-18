# Network Automation Lab

## General information

This repo contains virtual lab equipment for demonstrating network automation scenarios with `annet`.

Table of contents:

- [Installation](#installation)
- [Annet description](#annet-description)
- [Labs](#labs)
  - Basic
    - [lab00. Cisco Base Scenario](./topologies/lab00_basic_cisco)
    - [lab01. FRR Base Scenario](./topologies/lab01_basic_frr)
  - DC
    - [lab10. Cisco DC Scenario](./topologies/lab10_dc_cisco)
    - [lab11. FRR DC Scenario](./topologies/lab11_dc_frr)
    - [lab12. Multivendor DC Scenario](./topologies/lab12_dc_arista_cisco_frr)


### Environment

- Netbox url: http://localhost:8000/
- Netbox login/password: `annet/annet`
- Device telnet and ssh login/password: `annet/annet`  

## Lab installation

### Preparation

This steps are the same for all the labs.

1. First of all, you need to install Docker and Docker Compose on your device:
   - [Docker](https://docs.docker.com/engine/install/)
     - [Linux](https://docs.docker.com/desktop/install/linux/)
     - [Mac](https://docs.docker.com/desktop/install/mac-install/)
     - [Windows](https://docs.docker.com/desktop/install/windows-install/)

2. Install `make` utility:
   ```bash
   sudo apt install make  # Linux
   brew install make      # MacOS
   ```

3. Some labs require OS images (i.e. Arista EOS). Please download it according to Lab Guide and put to `../vm_images` directory. 

4. Clone this repository:
   ```bash
   git clone https://github.com/annetutil/annet.git
   ```

   Navigate to the lab folder:
   ```bash
   cd annetutils/contribs/labs
   ```

5. Build Annet and Netbox Docker images:
   ```bash
   make build
   ```

   After some changes you have to run `make rebuild`. It doesn't relate to changes in generators and mesh.

6. Run the Lab
   Now you can choose which scenario you want to run. To start a lab you need to run 
   ```bash
   make labXX
   ```
   where `XX` is an index of the lab.  
   For example, `make lab00` will start `lab00. Basic Cisco Scenario`.

   After this step you will be automatically logged in to annet container as a root. You can also login manually by `docker exec -u root -t -i annet /bin/bash`.

   When Cisco IOS is used in a lab, this step also automatically generates 512 bit RSA keys and enables SSH on devices. It can take a while.

## Annet description

[Annet](https://annetutil.github.io/annet/main/index.html) is a network configuration management system. It provides capabilities for storing configuration templates, generating and deploying network configurations.
Annet has four main arguments:

- `gen` — `annet gen $HOST` to generate desired configuration
- `diff` — `annet diff $HOST` to show text diff between desired and actual configurations
- `patch` — `annet patch $HOST` to prepare configuration patch with related commands
- `deploy` — `annet deploy $HOST` to generate patch and deploy it to devices

`$HOST` can be a single device or list of devices separated by spaces.

There are two main things which you will need to know and change to accomplish this labs:

- [generators](https://annetutil.github.io/annet/main/usage/gen.html) — python classes to yield configuration lines
- [mesh](https://annetutil.github.io/annet/main/mesh/index.html) — python classes to describe BGP design according to connections between devices.

To use annet in the lab topologies presented in this repository, you should prepare generators and mesh located in the `topologies/lab*/src/lab_generators/` folder.

---

Good luck, fellow kids!

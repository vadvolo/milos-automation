## FRR Base Lab

### Introduction

This lab demonstrates basic principles of network automation using FRR devices. The main goal is demonstrating Annet's `Entire` generator type.

Authors:
- [Grigorii Macheev](https://github.com/gregory-mac),
- [Vadim Volovik](https://github.com/vadvolo)
- [Grigorii Solovev](https://github.com/gs1571)

### Objectives

- Understand main principles of writing Annet `Entire` generators

### Topology

![Lab Topology](./images/topology.png)

### Generators

Unlike the `Partial` generators from the previous lab, which create and apply configuration line-by-line, the `Entire` type generates a whole configuration file in one go, which is then copied to the device.
FRR can be managed by `vtysh`, a Cisco-like CLI shell, but it also stores its configuration in a `/etc/frr/frr.conf` file.
We can leverage this fact to manage the routing configuration in a server-like manner, and `Partial` generator will help us to prepare the configuration file.

The generator in this example configures interface descriptions, IP addresses and BGP sessions between FRR routers.
All the parameters are defined by connections in Netbox.

### Lab Guide

| Router | CLI |
|:------:|:----|
| frr-r1 | `docker exec -u root -t -i frr-r1 vtysh` |
| frr-r2 | `docker exec -u root -t -i frr-r2 vtysh` |
| frr-r3 | `docker exec -u root -t -i frr-r3 vtysh` |


**Step 1.**
If it was not done in one of the previous labs, build Netbox and Annet docker images:

```bash
cd annetutils/contribs/labs
make build
```

**Step 2.**

Start the lab:

```bash
make lab01
```

**Step 3.**

Go to the Annet container:

```bash
docker exec -u root -t -i annet /bin/bash
```

Generate configuration for `frr-r1`, `frr-r2`, `frr-r3`:

`annet deploy frr-r1.nh.com frr-r2.nh.com frr-r3.nh.com`

Look at diff

`annet diff frr-r1.nh.com frr-r2.nh.com frr-r3.nh.com`

Deploy it

`annet deploy frr-r1.nh.com frr-r2.nh.com frr-r3.nh.com`

**Step 4.**

Remove connection between `frr-r1` and `frr-r2` in Netbox.

Look at diff

`annet diff frr-r1.nh.com frr-r2.nh.com frr-r3.nh.com`

Deploy it

`annet deploy frr-r1.nh.com frr-r2.nh.com frr-r3.nh.com`

**Step 5.**

Restore connection between `frr-r1` and `frr-r2` in Netbox.

Look at diff

`annet diff frr-r1.nh.com frr-r2.nh.com frr-r3.nh.com`

Deploy it

`annet deploy frr-r1.nh.com frr-r2.nh.com frr-r3.nh.com`

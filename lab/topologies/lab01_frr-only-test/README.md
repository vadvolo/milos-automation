## FRR Base Lab

### Introduction 

This lab demonstrates basic principles network automation using FRR devices. The main goal is demonstrating Entire generator of Annet.

Authors:
- [Grigorii Macheev](https://github.com/gregory-mac),
- [Vadim Volovik](https://github.com/vadvolo)
- [Grigorii Solovev](https://github.com/gs1571)

### Objectives

- Understand main principals of writing annet generators 

### Topology:

![Lab Topology](./images/topology.png)

### Generators

There is only one generator for FRR which is Entire generators. It means that Annet control whole configuration file of the service frr.
We should write that kind like one generator per one configuration file. 

The generator configure ip addresses and descriptions of the interfaces between routers. It configures BGP sessions too.
All the staff depends on connection map in Netbox.

### Lab Guide

| Router | CLI |
|:------:|:----|
| frr-r1 | `docker exec -u root -t -i frr-r1 vtysh` |
| frr-r2 | `docker exec -u root -t -i frr-r1 vtysh` |
| frr-r3 | `docker exec -u root -t -i frr-r1 vtysh` |


**Step 1.**
If it was not done in one of the previous labs, build Netbox and Annet docker images:

```bash
cd annetutils/labs
make build
```

**Step 2.**

Start the lab:

```bash
make lab01
```

**Step 3.**

Go to annet-container

```
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

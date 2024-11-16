## Lab03. Multivendor lab

### Introduction

In this lab we will see how Annet manages software from several different vendors:

- Arista
- FRR
- Cisco

Author:

### Objectives

- To try out Annet in a multivendor environment

### Preparation

#### How to get Arista image

1. Go to https://www.arista.com/en/login
2. Log in or register at Arista.com
3. Go to https://www.arista.com/en/support/software-download
4. Download the `cEOS64-lab-4.33.0F.tar.xz`
5. Prepare docker image: `docker image import cEOS64-lab-4.33.0F.tar.xz arista-ceos:4.33.0F`

Now you're able to run Lab03.

### Topology

![Lab Topology](./images/topology.png)

### Lab Guide

**Step 1.**

If it was not done in one of the previous labs, build Netbox and Annet docker images:

```bash
cd annetutils/contribs/labs
make build
```

**Step 2.**

NOTE: Do not forget to put Cisco IOS image `c7200-jk9s-mz.124-13a.bin` into `../vm_images` directory and Arista image import.

Start the lab:

```bash
make lab12
```

NOTE: On Linux, `make` uses root privileges to execute the following command:

```bash
$(SUDO) find operational_configs -mindepth 1 -not -name '.gitkeep' -delete || true && \
```

which is required to clear operational configs if they exist.

**Step 3.**

Go to the Annet container:

```bash
docker exec -u root -t -i annet /bin/bash
```

Enable SSH on Cisco routers by executing the script:

```bash
for ip in 0 1; do netsshsetup -a 172.20.0.10$ip -b ios -l annet -p annet -P telnet -v cisco --ipdomain nh.com; done
```

**Step 4.**

Go to the Annet container:

```bash
docker exec -u root -t -i annet /bin/bash
```

Generate configuration for spine-1-1, spine-1-2, tor-1-1, tor-1-2, tor-1-3:

`annet gen spine-1-1.nh.com spine-1-2.nh.com tor-1-1.nh.com tor-1-2.nh.com tor-1-3.nh.com`

Look at diff:

`annet diff spine-1-1.nh.com spine-1-2.nh.com tor-1-1.nh.com tor-1-2.nh.com tor-1-3.nh.com`

Deploy it:

`annet deploy spine-1-1.nh.com spine-1-2.nh.com tor-1-1.nh.com tor-1-2.nh.com tor-1-3.nh.com`

**Step 5.**

Assign "Unknown" role to one of the ToRs and deploy configuration on the ToR and every spine.

Go to the [Netbox](http://localhost:8000/), use annet:annet as login:password. Assign tor-1-1.nh.com role "Unknown".

Look at diff:

`annet diff spine-1-1.nh.com spine-1-2.nh.com tor-1-1.nh.com tor-1-2.nh.com tor-1-3.nh.com`

Deploy it:

`annet deploy spine-1-1.nh.com spine-1-2.nh.com tor-1-1.nh.com tor-1-2.nh.com tor-1-3.nh.com`

Restore the role and repeat the actions.

**Step 6.**

Break a connection and check what happens.

Go to [Netbox](http://localhost:8000/), use annet:annet as login:password. Delete the connection between tor-1-1.nh.com and spine-1-1.nh.com.

Look at diff:

`annet diff spine-1-1.nh.com spine-1-2.nh.com tor-1-1.nh.com tor-1-2.nh.com tor-1-3.nh.com`

Deploy it:

`annet deploy spine-1-1.nh.com spine-1-2.nh.com tor-1-1.nh.com tor-1-2.nh.com tor-1-3.nh.com`

Restore the connection and repeat the actions.

**Step 7.**

Drain traffic from one of the spines.

Go to [Netbox](http://localhost:8000/), use annet:annet as login:password. Assign spine-1-1.nh.com tag "maintenance".

Look at diff:

`annet diff spine-1-1.nh.com spine-1-2.nh.com tor-1-1.nh.com tor-1-2.nh.com tor-1-3.nh.com`

Deploy it:

`annet deploy spine-1-1.nh.com spine-1-2.nh.com tor-1-1.nh.com tor-1-2.nh.com tor-1-3.nh.com`

Remove the tag and repeat the actions.

**Step 8.**

After finishing the lab, stop it:

```bash
make services_stop
```

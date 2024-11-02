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
git checkout lab
```

Navigate to the lab folder:

```
cd lab-dc
```

Build and run lab

```
make build
make run
```

### How to connect to containers

Netbox: `docker exec -u root -t -i netbox-docker-netbox-1 /bin/bash`  
Annet: `docker exec -u root -t -i netbox-docker-annet-1 /bin/bash`  
Dynamips: `docker exec -u root -t -i netbox-docker-dynamips-lab-1 /bin/bash`

### Dynamips access

To access to the lab devices use username `cisco` and password `cisco`.
SSH will be enabled if you connect locally to each device and do commands:
```cisco
conf t
crypto key generate rsa general-keys modulus 2048
```

Commands to access to the console:
- spine-1-1.nh.com: `telnet 127.0.0.1 2001`
- spine-1-2.nh.com: `telnet 127.0.0.1 2002`
- tor-1-1.nh.com: `telnet 127.0.0.1 2003`
- tor-1-2.nh.com: `telnet 127.0.0.1 2004`
- tor-1-3.nh.com: `telnet 127.0.0.1 2005`

### Create Netbox SuperUser

**This step has already done in `make run` by coping netbox database.** Please use `root:root`.

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

### NETBOX TOKEN

**This step has already done in `make run` by coping netbox database.** You should use token `export NETBOX_TOKEN=c7fe6f5259a541e6e67ecacf094c439a2450c318`.  

1. Go to [LOCAL NETBOX INSTALLATION](http://localhost:8000/users/tokens/)
2. Push `+ Add` button
3. Generate token
4. Copy token
5. exit from netbox container
6. `export NETBOX_TOKEN=a630dcef...`
7. `make annet_restart`

## How to use

Run:

- diff

```
docker exec -ti netbox-docker-annet-1 python3 -m annet.annet diff spine-1-1.nh.com
```

- patch

```
docker exec -ti netbox-docker-annet-1 python3 -m annet.annet patch spine-1-1.nh.com
```

- deploy

```
docker exec -ti netbox-docker-annet-1 python3 -m annet.annet deploy spine-1-1.nh.com
```

### Topology

![lab-topology](nh2024-lab.png "Title")

Naming:
- Spine - `spine-<pod>-<plnane>`
- ToR - `tor-<pod>-<num>`
- Router ID Spine - `1.2.<pod>.<plane>`
- Router ID ToR - `1.1.<pod>.<num>`
- ASNUM Spine - `6520<pod>`
- ASNUM ToR - `6510<pod><num>`

### Lab actions:
1. Deploy whole configuration.
2. Replace role of one of the ToR's to "Unknown" and deploy configuration on the ToR and every Spine.
3. Restore role of the ToR to "ToR" and deploy configuration on the ToR and every Spine.
4. Assign tag "maintenance" to one of the Spine's and deploy configuration on the Spine.
5. Remove tag "maintenance" from the Spine and deploy configuration again.

### Some comments

#### Required for annet repo
1. FAQ
   - how do implicit
   - how work with rulebook
   - description of ACL language
2. Option to move rule book out of repo

#### Issues
1. ~~swap creds to cisco:cisco~~
2. implicit:
   - ip address
   - no shutdown
3. fix Order
4. cisco deploy error `gnetclisdk.exceptions.GnetcliException: AioRpcError read timeout error. last seen: "Destination filename [startup-config]? "`
5. ann deploy without progress
6. ann doesn't failed if faced with config command error on Cisco IOS
7. after changing rpl do `clear ip bgp * soft`

#### Annet changes

To successfully deploy configuration you need make some changes in annet code:
```diff
diff --git a/annet/implicit.py b/annet/implicit.py
index 850139a..8b7ab33 100644
--- a/annet/implicit.py
+++ b/annet/implicit.py
@@ -133,6 +133,11 @@ def _implicit_tree(device):
                         no shutdown
             """
     elif device.hw.Cisco:
+        text += r"""
+            !interface
+                no shutdown
+                no ip address
+        """
         if device.hw.Cisco.Catalyst:
             # this configuration is not visible in running-config when enabled
             text += r"""
diff --git a/annet/rulebook/texts/cisco.order b/annet/rulebook/texts/cisco.order
index 741770c..3d09f54 100644
--- a/annet/rulebook/texts/cisco.order
+++ b/annet/rulebook/texts/cisco.order
@@ -52,8 +52,6 @@ spanning-tree
 no system jumbomtu %order_reverse
 system jumbomtu

-route-map
-
 service dhcp
 ip dhcp relay
 ipv6 dhcp relay
@@ -74,6 +72,8 @@ interface *

 interface */\S+\.\d+/

+route-map
+
 # удалять eth-trunk можно только после того, как вычистим member interfaces
 undo interface */port-channel\d+/  %order_reverse
 
 router bgp
+    neighbor */[\da-f\.\:]+/ remote-as
+    neighbor */[\da-f\.\:]+/ peer-group
 
 line
\ No newline at end of file
diff --git a/annet/rulebook/texts/cisco.rul b/annet/rulebook/texts/cisco.rul
index 4486a28..5dd2683 100644
--- a/annet/rulebook/texts/cisco.rul
+++ b/annet/rulebook/texts/cisco.rul
@@ -42,11 +42,16 @@ no snmp-server sysobjectid type stack-oid

 !interface */(mgmt|ipmi|Vlan1$)/

+interface Loopback
+    ipv6 address *
+    ip address ~
+
 # SVI/Subifs/Lagg
 interface */(Vlan|Ethernet.*\.|port-channel.*\.?)\d+$/ %diff_logic=cisco.iface.diff
     vrf member
     ipv6 link-local
     ipv6 address *
+    ip address ~
     ipv6 nd ~                      %logic=cisco.misc.no_ipv6_nd_suppress_ra
     mtu

@@ -59,6 +64,7 @@ interface */\w*Ethernet[0-9\/]+$/     %logic=common.permanent %diff_logic=cisco.
     vrf member
     ipv6 link-local
     ipv6 address *
+    ip address ~
     channel-group
     mtu
     storm-control * level
```

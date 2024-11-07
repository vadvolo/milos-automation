# Cisco DC Lab

## Preparation

Before you start please put into `../vm_images` Cisco IOS image `c7200-jk9s-mz.124-13a.bin`

## Topology

![lab-topology](nh2024-lab.png "Title")

Naming:
- Spine - `spine-<pod>-<plnane>`
- ToR - `tor-<pod>-<num>`
- Router ID Spine - `1.2.<pod>.<plane>`
- Router ID ToR - `1.1.<pod>.<num>`
- ASNUM Spine - `6520<pod>`
- ASNUM ToR - `6510<pod><num>`

### Lab Guide

### Lab actions:
1. Deploy whole configuration.
2. Replace role of one of the ToR's to "Unknown" and deploy configuration on the ToR and every Spine.
3. Restore role of the ToR to "ToR" and deploy configuration on the ToR and every Spine.
4. Assign tag "maintenance" to one of the Spine's and deploy configuration on the Spine.
5. Remove tag "maintenance" from the Spine and deploy configuration again.

**Step 1.**  
To start lab please navigate to `annetutils/labs` and run `make lab03`.

**Step 2.**  
Go to annet-container  
```
docker exec -u root -t -i netbox-docker-annet-1 /bin/bash
```

**Step 3.** 

Enable SSH on Cisco routers by script:
```
/home/ubuntu/scripts/netsshsetup/netsshsetup -a 172.20.0.100 -v cisco -b ios -l cisco -p cisco -P telnet --hostname spine-1-1.nh.com
/home/ubuntu/scripts/netsshsetup/netsshsetup -a 172.20.0.101 -v cisco -b ios -l cisco -p cisco -P telnet --hostname spine-1-2.nh.com
/home/ubuntu/scripts/netsshsetup/netsshsetup -a 172.20.0.102 -v cisco -b ios -l cisco -p cisco -P telnet --hostname tor-1-1.nh.com
/home/ubuntu/scripts/netsshsetup/netsshsetup -a 172.20.0.103 -v cisco -b ios -l cisco -p cisco -P telnet --hostname tor-1-2.nh.com
/home/ubuntu/scripts/netsshsetup/netsshsetup -a 172.20.0.104 -v cisco -b ios -l cisco -p cisco -P telnet --hostname tor-1-3.nh.com
```

**Step 4.**

| Router | Command |
|:------:|:------:|
| spine-1-1 |`python3 -m annet.annet gen spine-1-1.nh.com` | 
| spine-1-2 |`python3 -m annet.annet gen spine-1-2.nh.com` |
| tor-1-1 |`python3 -m annet.annet gen tor-1-1.nh.com` | 
| tor-1-2 |`python3 -m annet.annet gen tor-1-2.nh.com` |
| tor-1-2 |`python3 -m annet.annet gen tor-1-3.nh.com` | 

> If you see error below, you need to export NETBOX_TOKEN to the Annet container.
> ```
>   File "/venv/lib/python3.12/site-packages/dataclass_rest/http/requests.py", line 19, in _on_error_default
>     raise ClientError(response.status_code)
> dataclass_rest.exceptions.ClientError: 403
> ```
>
> ```
> export NETBOX_TOKEN="a630dcefcb191982869e7576190e79bfd569d33c"
> ```

<details>
<summary>Output for spine-1-1:</summary>

```
hostname spine-1-1
ip bgp-community new-format
ip community-list standard TOR_NETS permit 65000:1
ip community-list standard GRACEFUL_SHUTDOWN permit 65535:0
interface GigabitEthernet1/0
  no shutdown
  ip address 192.168.11.1 255.255.255.0
  description tor-1-1@Gi1/0
interface GigabitEthernet2/0
  no shutdown
  ip address 192.168.12.1 255.255.255.0
  description tor-1-2@Gi1/0
interface GigabitEthernet3/0
  no shutdown
  ip address 192.168.13.1 255.255.255.0
  description tor-1-3@Gi1/0
interface FastEthernet0/0
  no shutdown
  ip address 172.20.0.100 255.255.255.0
interface FastEthernet0/1
  no shutdown
route-map TOR_IMPORT permit 10
  match community TOR_NETS
route-map TOR_IMPORT deny 9999
route-map TOR_EXPORT permit 10
  match community TOR_NETS
route-map TOR_EXPORT deny 9999
router bgp 65201
  bgp router-id 1.2.1.1
  bgp log-neighbor-changes
  neighbor TOR peer-group
  neighbor TOR route-map TOR_IMPORT in
  neighbor TOR route-map TOR_EXPORT out
  neighbor TOR soft-reconfiguration inbound
  neighbor TOR send-community both
  neighbor 192.168.11.2 remote-as 65111
  neighbor 192.168.12.2 remote-as 65112
  neighbor 192.168.13.2 remote-as 65113
  neighbor 192.168.11.2 peer-group TOR
  neighbor 192.168.12.2 peer-group TOR
  neighbor 192.168.13.2 peer-group TOR
```

</details>

<details>
<summary>Output for spine-1-2:</summary>

```
hostname spine-1-2
ip bgp-community new-format
ip community-list standard TOR_NETS permit 65000:1
ip community-list standard GRACEFUL_SHUTDOWN permit 65535:0
interface GigabitEthernet1/0
  no shutdown
  ip address 192.168.21.1 255.255.255.0
  description tor-1-1@Gi2/0
interface GigabitEthernet2/0
  no shutdown
  ip address 192.168.22.1 255.255.255.0
  description tor-1-2@Gi2/0
interface GigabitEthernet3/0
  no shutdown
  ip address 192.168.23.1 255.255.255.0
  description tor-1-3@Gi2/0
interface FastEthernet0/0
  no shutdown
  ip address 172.20.0.101 255.255.255.0
interface FastEthernet0/1
  no shutdown
route-map TOR_IMPORT permit 10
  match community TOR_NETS
route-map TOR_IMPORT deny 9999
route-map TOR_EXPORT permit 10
  match community TOR_NETS
route-map TOR_EXPORT deny 9999
router bgp 65201
  bgp router-id 1.2.1.2
  bgp log-neighbor-changes
  neighbor TOR peer-group
  neighbor TOR route-map TOR_IMPORT in
  neighbor TOR route-map TOR_EXPORT out
  neighbor TOR soft-reconfiguration inbound
  neighbor TOR send-community both
  neighbor 192.168.21.2 remote-as 65111
  neighbor 192.168.22.2 remote-as 65112
  neighbor 192.168.23.2 remote-as 65113
  neighbor 192.168.21.2 peer-group TOR
  neighbor 192.168.22.2 peer-group TOR
  neighbor 192.168.23.2 peer-group TOR
```

</details>

<details>
<summary>Output for tor-1-1:</summary>

```
hostname tor-1-1
ip bgp-community new-format
ip community-list standard TOR_NETS permit 65000:1
ip community-list standard GRACEFUL_SHUTDOWN permit 65535:0
interface GigabitEthernet1/0
  no shutdown
  ip address 192.168.11.2 255.255.255.0
  description spine-1-1@Gi1/0
interface GigabitEthernet2/0
  no shutdown
  ip address 192.168.21.2 255.255.255.0
  description spine-1-2@Gi1/0
interface FastEthernet0/0
  no shutdown
  ip address 172.20.0.102 255.255.255.0
interface Loopback0
  no shutdown
  ip address 10.0.0.1 255.255.255.255
interface FastEthernet0/1
  no shutdown
interface GigabitEthernet3/0
  no shutdown
route-map SPINE_IMPORT permit 10
  match community TOR_NETS GRACEFUL_SHUTDOWN
  set local-preference 0
route-map SPINE_IMPORT permit 20
  match community TOR_NETS
  set local-preference 100
route-map SPINE_IMPORT deny 9999
route-map SPINE_EXPORT permit 10
  match community TOR_NETS
route-map SPINE_EXPORT deny 9999
route-map CONNECTED permit 10
  match interface Loopback0
  set community 65000:1
route-map CONNECTED deny 9999
router bgp 65111
  bgp router-id 1.1.1.1
  bgp log-neighbor-changes
  redistribute connected route-map CONNECTED
  maximum-paths 16
  neighbor SPINE peer-group
  neighbor SPINE route-map SPINE_IMPORT in
  neighbor SPINE route-map SPINE_EXPORT out
  neighbor SPINE soft-reconfiguration inbound
  neighbor SPINE send-community both
  neighbor 192.168.11.1 remote-as 65201
  neighbor 192.168.21.1 remote-as 65201
  neighbor 192.168.11.1 peer-group SPINE
  neighbor 192.168.21.1 peer-group SPINE
```

</details>

<details>
<summary>Output for tor-1-2:</summary>

```
hostname tor-1-2
ip bgp-community new-format
ip community-list standard TOR_NETS permit 65000:1
ip community-list standard GRACEFUL_SHUTDOWN permit 65535:0
interface GigabitEthernet1/0
  no shutdown
  ip address 192.168.12.2 255.255.255.0
  description spine-1-1@Gi2/0
interface GigabitEthernet2/0
  no shutdown
  ip address 192.168.22.2 255.255.255.0
  description spine-1-2@Gi2/0
interface FastEthernet0/0
  no shutdown
  ip address 172.20.0.103 255.255.255.0
interface Loopback0
  no shutdown
  ip address 10.0.0.2 255.255.255.255
interface FastEthernet0/1
  no shutdown
interface GigabitEthernet3/0
  no shutdown
route-map SPINE_IMPORT permit 10
  match community TOR_NETS GRACEFUL_SHUTDOWN
  set local-preference 0
route-map SPINE_IMPORT permit 20
  match community TOR_NETS
  set local-preference 100
route-map SPINE_IMPORT deny 9999
route-map SPINE_EXPORT permit 10
  match community TOR_NETS
route-map SPINE_EXPORT deny 9999
route-map CONNECTED permit 10
  match interface Loopback0
  set community 65000:1
route-map CONNECTED deny 9999
router bgp 65112
  bgp router-id 1.1.1.2
  bgp log-neighbor-changes
  redistribute connected route-map CONNECTED
  maximum-paths 16
  neighbor SPINE peer-group
  neighbor SPINE route-map SPINE_IMPORT in
  neighbor SPINE route-map SPINE_EXPORT out
  neighbor SPINE soft-reconfiguration inbound
  neighbor SPINE send-community both
  neighbor 192.168.12.1 remote-as 65201
  neighbor 192.168.22.1 remote-as 65201
  neighbor 192.168.12.1 peer-group SPINE
  neighbor 192.168.22.1 peer-group SPINE
```

</details>

<details>
<summary>Output for tor-1-3:</summary>

```
hostname tor-1-3
ip bgp-community new-format
ip community-list standard TOR_NETS permit 65000:1
ip community-list standard GRACEFUL_SHUTDOWN permit 65535:0
interface GigabitEthernet1/0
  no shutdown
  ip address 192.168.13.2 255.255.255.0
  description spine-1-1@Gi3/0
interface GigabitEthernet2/0
  no shutdown
  ip address 192.168.23.2 255.255.255.0
  description spine-1-2@Gi3/0
interface FastEthernet0/0
  no shutdown
  ip address 172.20.0.104 255.255.255.0
interface Loopback0
  no shutdown
  ip address 10.0.0.3 255.255.255.255
interface FastEthernet0/1
  no shutdown
interface GigabitEthernet3/0
  no shutdown
route-map SPINE_IMPORT permit 10
  match community TOR_NETS GRACEFUL_SHUTDOWN
  set local-preference 0
route-map SPINE_IMPORT permit 20
  match community TOR_NETS
  set local-preference 100
route-map SPINE_IMPORT deny 9999
route-map SPINE_EXPORT permit 10
  match community TOR_NETS
route-map SPINE_EXPORT deny 9999
route-map CONNECTED permit 10
  match interface Loopback0
  set community 65000:1
route-map CONNECTED deny 9999
router bgp 65113
  bgp router-id 1.1.1.3
  bgp log-neighbor-changes
  redistribute connected route-map CONNECTED
  maximum-paths 16
  neighbor SPINE peer-group
  neighbor SPINE route-map SPINE_IMPORT in
  neighbor SPINE route-map SPINE_EXPORT out
  neighbor SPINE soft-reconfiguration inbound
  neighbor SPINE send-community both
  neighbor 192.168.13.1 remote-as 65201
  neighbor 192.168.23.1 remote-as 65201
  neighbor 192.168.13.1 peer-group SPINE
  neighbor 192.168.23.1 peer-group SPINE
```

</details>

<details>
<summary>Diff for spine-1-1:</summary>

```diff

```

</details>

<details>
<summary>Diff for spine-1-2:</summary>

```diff

```

</details>

<details>
<summary>Diff for tor-1-1:</summary>

```diff

```

</details>

<details>
<summary>Diff for tor-1-2:</summary>

```diff

```

</details>

<details>
<summary>Diff for tor-1-3:</summary>

```diff

```

</details>

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

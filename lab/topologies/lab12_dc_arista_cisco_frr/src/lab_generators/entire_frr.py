from typing import Optional

from annet.bgp_models import ASN
from annet.generators import Entire
from annet.mesh.executor import MeshExecutionResult
from annet.storage import Device

from .helpers.router import (
    AutonomusSystemIsNotDefined,
    bgp_asnum,
    bgp_groups,
    bgp_mesh,
    router_id,
)


class Frr(Entire):

    TAGS = ["frr"]

    def path(self, device: Device):
        if device.hw.PC:
            return "/etc/frr/frr.conf"

    def reload(self, _) -> str:
        return "sudo /etc/init.d/frr reload"

    def run(self, device: Device):

        mesh_data: MeshExecutionResult = bgp_mesh(device)

        # base configuration
        yield "frr defaults datacenter"
        yield "service integrated-vtysh-config"
        yield ""
        yield "hostname", device.hostname.split(".")[0]
        yield "log file /var/log/frr/frr.log"
        yield ""

        # interface configuration
        for interface in device.interfaces:
            yield "interface", interface.name
            description = ""
            if interface.connected_endpoints:
                remote = interface.connected_endpoints[0]
                description = f"{remote.device.name}@{remote.name}"

            if description:
                yield " description", description
            if interface.ip_addresses:
                for ip in interface.ip_addresses:
                    if ip.family.value == 4:
                        yield " ip address", ip.address
                    if ip.family.value == 6:
                        yield " ipv6 address", ip.address
            yield "exit"
            yield ""

        # bgp configuration
        rid: Optional[str] = router_id(mesh_data)
        try:
            asnum: Optional[ASN] = bgp_asnum(mesh_data)
        except AutonomusSystemIsNotDefined as err:
            raise RuntimeError(f"Device {device.name} has more than one defined autonomus system: {err}")

        if asnum and rid:
            yield "router bgp", asnum
            yield " bgp router-id", rid

            for group in bgp_groups(mesh_data):
                yield " neighbor", group.group_name, "peer-group"

            for peer in mesh_data.peers:
                yield " neighbor", peer.addr, "remote-as", peer.remote_as
                yield " neighbor", peer.addr, "peer-group", peer.group_name

            yield " address-family ipv4 unicast"

            if mesh_data.global_options and mesh_data.global_options.ipv4_unicast and mesh_data.global_options.ipv4_unicast.redistributes:
                for redistribute in mesh_data.global_options.ipv4_unicast.redistributes:
                    yield "  redistribute", redistribute.protocol, "" if not redistribute.policy else f"route-map {redistribute.policy}"

            for group in bgp_groups(mesh_data):
                yield "  neighbor", group.group_name, "route-map", group.import_policy, "in"
                yield "  neighbor", group.group_name, "route-map", group.export_policy, "out"

            if device.device_role.name == "ToR":
                yield "  maximum-paths 16"

            yield " exit-address-family"
            yield "exit"
            yield ""

        # route-map configuration
        yield "bgp community-list standard TOR_NETS seq 5 permit 65000:1"
        if device.device_role.name == "ToR":
            yield """
bgp community-list standard TOR_NETS_WITH_GSHUT seq 5 permit 65000:1
bgp community-list standard TOR_NETS_WITH_GSHUT seq 10 permit graceful-shutdown

route-map SPINE_IMPORT permit 10
 match community TOR_NETS_WITH_GSHUT exact-match
 set local-preference 0
exit

route-map SPINE_IMPORT permit 20
 match community TOR_NETS
 set local-preference 100
exit

route-map SPINE_IMPORT deny 9999
exit

route-map SPINE_EXPORT permit 10
 match community TOR_NETS
exit

route-map SPINE_EXPORT deny 9999
exit

route-map CONNECTED permit 10
 match interface lo
 set community 65000:1
exit

route-map CONNECTED deny 9999
exit
"""
        elif device.device_role.name == "Spine":
            yield """
route-map TOR_IMPORT permit 10
 match community TOR_NETS
exit

route-map TOR_IMPORT deny 9999
exit
"""
        yield "line vty"

from dataclasses import dataclass
from typing import Optional

from annet.bgp_models import ASN
from annet.generators import Entire
from annet.mesh.executor import MeshExecutionResult
from annet.storage import Device

from .helpers.router import AutonomusSystemIsNotDefined, bgp_asnum, bgp_mesh, router_id


class Frr(Entire):
    
    TAGS = ["frr"]
    
    def path(self, device: Device):
        if device.hw.PC:
            return "/etc/frr/frr.conf"

    def reload(self, _) -> str:
        return "sudo /etc/init.d/frr reload"

    def run(self, device: Device):

        mesh_data: MeshExecutionResult = bgp_mesh(device)

        yield "frr defaults datacenter"
        yield "service integrated-vtysh-config"
        yield ""
        yield "hostname", device.hostname.split(".")[0]
        yield "log file /var/log/frr/frr.log"
        yield ""
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

        rid: Optional[str] = router_id(mesh_data)
        try:
            asnum: Optional[ASN] = bgp_asnum(mesh_data)
        except AutonomusSystemIsNotDefined as err:
            RuntimeError(f"Device {device.name} has more than one defined autonomus system: {err}")

        yield "router bgp", asnum
        if router_id:
            yield " bgp router-id", rid
            for peer in mesh_data.peers:
                yield " neighbor", peer.addr, "remote-as", peer.asnum
                yield " neighbor", peer.addr, "interface", peer.source

        yield ""
        yield "line vty"

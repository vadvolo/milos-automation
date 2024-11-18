from dataclasses import dataclass

from annet.generators import Entire
from annet.storage import Device


MTU = 1500


@dataclass
class BgpPeer:
    """Structure of BGP peer"""

    addr: str
    asnum: str
    source: str


class Frr(Entire):
    """Entire generator class for Frrouting"""

    TAGS = ["frr"]

    def path(self, device: Device):
        """Define vendor and path to the configuration file"""

        if device.hw.PC:
            return "/etc/frr/frr.conf"

    def reload(self, _) -> str:
        """define action which should be done in case of configuration file changes"""

        return "sudo /etc/init.d/frr reload"

    def run(self, device: Device):
        """Generate configuration file content"""

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

        router_id = ""
        if device.primary_ip and device.primary_ip.family == 4:
            router_id = device.primary_ip.address.split("/")[0]

        yield "router bgp 65001"
        if router_id:
            yield " bgp router-id", router_id
            for peer in _bgp_peers(device):
                yield " neighbor", peer.addr, "remote-as", peer.asnum
                yield " neighbor", peer.addr, "interface", peer.source

        yield ""
        yield "line vty"
        yield ""


def _bgp_peers(device: Device) -> list[BgpPeer]:
    """Return list of BGP peers"""

    res: list[BgpPeer] = []
    for interface in device.interfaces:
        if interface.connected_endpoints:
            remote = interface.connected_endpoints[0]
            remote_host = remote.device.name
            remote_iface = remote.name
            remote_ip = _get_neighbor_iface_address(device, remote_host, remote_iface)

            res.append(
                BgpPeer(
                    addr=remote_ip,
                    asnum="65001",
                    source=interface.name
                )
            )
    return res


def _get_neighbor_iface_address(device: Device, remote_host: str, remote_iface: str) -> str:
    """Return IP address of remote peer"""

    for neighbour in device.neighbours:
        if neighbour.name == remote_host:
            for interface in neighbour.interfaces:
                if interface.name == remote_iface:
                    return interface.ip_addresses[0].address.split("/")[0]
    raise RuntimeError(f"Remote address not found for neighbor {remote_host}@{remote_iface}")

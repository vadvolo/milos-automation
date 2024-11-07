import re

from annet.generators import PartialGenerator
from annet.storage import Device
from ipaddress import IPv4Address, IPv6Address, IPv4Network, IPv6Network


def parse_device_name(device_name: str) -> dict[str, int]:
    """Parse device name to get """
    re_tor = r"^tor\-(?P<pod>\d+)\-(?P<num>\d+)"
    re_spine = r"^spine\-(?P<pod>\d+)\-(?P<plane>\d+)"

    m = re.match(re_tor, device_name)
    if m:
        tor_attr: dict[str, str] = m.groupdict()
        return {
            "pod": int(tor_attr["pod"]),
            "num": int(tor_attr["num"]),
        }

    m = re.match(re_spine, device_name)
    if m:
        spine_attr: dict[str, str] = m.groupdict()
        return {
            "pod": int(spine_attr["pod"]),
            "plane": int(spine_attr["plane"]),
        }
    raise RuntimeError(f"Could parse name '{device_name}'")


class IpAddress(PartialGenerator):
    
    TAGS = ["l3", "iface"]
    
    def acl_cisco(self, _: Device):
        return """
        interface
            ip address
            ipv6 address
        """
    
    def run_cisco(self, device: Device):
        for interface in device.interfaces:
            if interface.ip_addresses:
                count_ipv4 = 0
                count_ipv6 = 0
                with self.block(f"interface {interface.name}"):
                    for ip_address in interface.ip_addresses:
                        if ip_address.family.value == 4:
                            ip_addr: str = str(IPv4Address(ip_address.address.split("/")[0]))
                            ip_mask: str = str(IPv4Network(ip_address.address, strict=False).netmask)
                            secondary: str = "" if count_ipv4 == 0 else "secondary"
                            yield "ip address", ip_addr, ip_mask, secondary
                            count_ipv4 += 1
                        elif ip_address.family.value == 6:
                            ip_addr: str = str(IPv6Address(ip_address.address.split("/")[0]))
                            ip_mask: int = IPv6Network(ip_address.address, strict=False).prefixlen
                            secondary: str = "" if count_ipv6 == 0 else "secondary"
                            yield f"ipv6 address {ip_addr}/{ip_mask}", secondary
                            count_ipv6 += 1

        if device.device_role.name == "Spine":
            parsed_name_spine: dict[str, int] = parse_device_name(device.name)
            spine_plane: int = parsed_name_spine["plane"]
            for remote_device in device.neighbours:
                if remote_device.device_role.name == "ToR":
                    parsed_name_tor: dict[str, int] = parse_device_name(remote_device.name)
                    tor_num: int = parsed_name_tor["num"]
                    iface_name: str = remote_device.interfaces[0].connected_endpoints[0].name
                    with self.block("interface", iface_name):
                        yield f"ip address 192.168.{spine_plane}{tor_num}.1 255.255.255.0"

        elif device.device_role.name == "ToR":
            parsed_name_tor: dict[str, int] = parse_device_name(device.name)
            tor_num: int = parsed_name_tor["num"]
            for remote_device in device.neighbours:
                if remote_device.device_role.name == "Spine":
                    parsed_name_spine: dict[str, int] = parse_device_name(remote_device.name)
                    spine_plane: int = parsed_name_spine["plane"]
                    iface_name: str = remote_device.interfaces[0].connected_endpoints[0].name
                    with self.block("interface", iface_name):
                        yield f"ip address 192.168.{spine_plane}{tor_num}.2 255.255.255.0"



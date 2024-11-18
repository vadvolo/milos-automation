from ipaddress import IPv4Address, IPv4Network

from annet.generators import PartialGenerator
from annet.storage import Device

from .helpers.router import bgp_mesh


class IpAddresses(PartialGenerator):
    """Partial generator class of IPv4 and IPv6 addresses on interfaces"""

    TAGS = ["l3", "iface"]

    def acl_cisco(self, _: Device):
        """ACL for Cisco devices"""

        return """
        interface
            ip address
            ipv6 address
        """

    def run_cisco(self, device: Device):
        """Generator for Cisco devices"""

        # enrich interfaces by mesh
        bgp_mesh(device)
        for interface in device.interfaces:
            if interface.ip_addresses:
                count_ipv4 = 0
                count_ipv6 = 0
                with self.block(f"interface {interface.name}"):
                    # to deduplicate ip addresses on the interface, it looks like a bug in mesh
                    ip_addr_set: set[tuple[str, int]] = set()
                    for ip_address in interface.ip_addresses:
                        if not (ip_address.address, ip_address.family.value,) in ip_addr_set:
                            ip_addr_set.add((ip_address.address, ip_address.family.value,))
                            if ip_address.family.value == 4:
                                ip_addr: str = str(IPv4Address(ip_address.address.split("/")[0]))
                                ip_mask: str = str(IPv4Network(ip_address.address, strict=False).netmask)
                                secondary: str = "" if count_ipv4 == 0 else "secondary"
                                yield "ip address", ip_addr, ip_mask, secondary
                                count_ipv4 += 1
                            elif ip_address.family.value == 6:
                                secondary = "" if count_ipv6 == 0 else "secondary"
                                yield "ipv6 address", ip_address.address, secondary
                                count_ipv6 += 1

    def acl_arista(self, _: Device):
        """ACL for Arista devices"""
        return """
        interface
            ip address
            ipv6 address
        """

    def run_arista(self, device: Device):
        """Generator for Arista devices"""

        # enrich interfaces by mesh
        bgp_mesh(device)
        for interface in device.interfaces:
            if interface.ip_addresses:
                count_ipv4 = 0
                count_ipv6 = 0
                with self.block(f"interface {interface.name}"):
                    # to deduplicate ip addresses on the interface, it looks like a bug in mesh
                    ip_addr_set: set[tuple[str, int]] = set()
                    for ip_address in interface.ip_addresses:
                        if not (ip_address.address, ip_address.family.value,) in ip_addr_set:
                            ip_addr_set.add((ip_address.address, ip_address.family.value,))
                            if ip_address.family.value == 4:
                                secondary: str = "" if count_ipv4 == 0 else "secondary"
                                yield "ip address", ip_address.address, secondary
                                count_ipv4 += 1
                            elif ip_address.family.value == 6:
                                secondary = "" if count_ipv6 == 0 else "secondary"
                                yield "ipv6 address", ip_address.address, secondary
                                count_ipv6 += 1

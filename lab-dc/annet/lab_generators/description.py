from annet.generators import PartialGenerator
from annet.storage import Device


class Description(PartialGenerator):
    
    TAGS = ["description", "iface"]
    
    def acl_cisco(self, device):
        return """
        interface
            description
        """
    
    def run_cisco(self, device):
        for interface in device.interfaces:
            if interface.connected_endpoints:
                with self.block(f"interface {interface.name}"):
                    remote_device = interface.connected_endpoints[0].device.name.split(".")[0]
                    remote_iface = sorten_port_names(interface.connected_endpoints[0].name, device.device_type)
                    yield f"description {remote_device}@{remote_iface}"
    

# class Ifaces(PartialGenerator):
#     TAGS = ["mgmt", "lldp"]

#     def acl_routeros(self, _):
#         return """
#         ip
#           neighbor              %cant_delete
#             discovery-settings  %cant_delete
#               set               %cant_delete

#         """

#     def run_routeros(self, _: Device):
#         with self.block("ip"):
#             with self.block("neighbor"):
#                 with self.block("discovery-settings"):
#                     yield "set discover-interface-list=LAN lldp-med-net-policy-vlan=2 protocol=cdp"


# class Ntp(PartialGenerator):
#     TAGS = ["mgmt", "ntp"]

#     def acl_routeros(self, _):
#         return """
#         system
#           ntp         %cant_delete
#             client    %cant_delete
#               set     %cant_delete
#               server  %cant_delete
#                 add   %cant_delete
#         """

#     #     /system ntp client
#     # set enabled=yes
#     # /system ntp client servers
#     # add address=ntp0.ntp-servers.net
#     def run_routeros(self, device: Device):
#         with self.block("system"):
#             with self.block("ntp"):
#                 with self.block("client"):
#                     yield "set enabled=yes"

#         with self.block("system"):
#             with self.block("ntp"):
#                 with self.block("client"):
#                     with self.block("server"):
#                         yield "add address=ntp0.ntp-servers.net"


def sorten_port_names(portname: str, device_type) -> str:
    if device_type.manufacturer.name == "Cisco":
        if portname.startswith("GigabitEthernet"):
            return portname.replace("GigabitEthernet", "Gi")
        elif portname.startswith("FastEthernet"):
            return portname.replace("FastEthernet", "Fa")
    return portname

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


def sorten_port_names(portname: str, device_type) -> str:
    if device_type.manufacturer.name == "Cisco":
        if portname.startswith("GigabitEthernet"):
            return portname.replace("GigabitEthernet", "Gi")
        elif portname.startswith("FastEthernet"):
            return portname.replace("FastEthernet", "Fa")
    return portname

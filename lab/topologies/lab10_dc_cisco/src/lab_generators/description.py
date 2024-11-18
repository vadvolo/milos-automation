from annet.adapters.netbox.common.models import DeviceType
from annet.generators import PartialGenerator
from annet.storage import Device


class Description(PartialGenerator):
    """Partial generator class of description on interfaces"""

    TAGS = ["description", "iface"]

    def acl_cisco(self, _: Device):
        """ACL for Cisco devices"""

        return """
        interface
            description
        """

    def run_cisco(self, device: Device):
        """Generator for Cisco devices"""

        for interface in device.interfaces:
            if interface.connected_endpoints:
                with self.block(f"interface {interface.name}"):
                    remote_device = interface.connected_endpoints[0].device.name.split(".")[0]
                    remote_iface = _sorten_port_names(interface.connected_endpoints[0].name, device.device_type)
                    yield f"description {remote_device}@{remote_iface}"

    def acl_arista(self, _: Device):
        """ACL for Arista devices"""

        return """
        interface
            description
        """

    def run_arista(self, device: Device):
        """Generator for Arista devices"""

        for interface in device.interfaces:
            if interface.connected_endpoints:
                with self.block(f"interface {interface.name}"):
                    remote_device = interface.connected_endpoints[0].device.name.split(".")[0]
                    remote_iface = _sorten_port_names(interface.connected_endpoints[0].name, device.device_type)
                    yield f"description {remote_device}@{remote_iface}"


def _sorten_port_names(portname: str, device_type: DeviceType) -> str:
    """Generate short name of interface"""

    if device_type.manufacturer.name == "Cisco":
        if portname.startswith("GigabitEthernet"):
            return portname.replace("GigabitEthernet", "Gi")
        elif portname.startswith("FastEthernet"):
            return portname.replace("FastEthernet", "Fa")
    return portname

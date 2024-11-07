from annet.generators import PartialGenerator
from annet.storage import Device

MTU = 1500

class IfaceDescriptions(PartialGenerator):
    
    TAGS = ["description"]
    
    def acl_cisco(self, device):
        return """
        interface
            description
        """
    
    def run_cisco(self, device):
        for interface in device.interfaces:
            neighbor = ""
            if interface.connected_endpoints:
                for connection in interface.connected_endpoints:
                    neighbor += f"to_{connection.device.name}_{connection.name}"
                with self.block(f"interface {interface.name}"):
                    yield f"description {neighbor}"
            else:
                with self.block(f"interface {interface.name}"):
                    yield f"description disconnected"

class IfaceMtu(PartialGenerator):
    
    TAGS = ["description"]
    
    def acl_cisco(self, device):
        return """
        interface
            mtu
        """
    
    def run_cisco(self, device):
        for interface in device.interfaces:
            if interface.mtu:
                mtu = interface.mtu
            else:
                mtu = MTU
            with self.block(f"interface {interface.name}"):
                yield f"mtu {mtu}"

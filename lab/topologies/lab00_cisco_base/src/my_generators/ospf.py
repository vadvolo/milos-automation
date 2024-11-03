from annet.generators import PartialGenerator
from annet.storage import Device

class Ospf(PartialGenerator):
    
    TAGS = ["ospf"]
    
    def acl_cisco(self, device):
        return """
        router ospf
            network
        """
    
    def run_cisco(self, device):
        for interface in device.interfaces:
            with self.block(f"interface {interface.name}"):
                yield "description MILOS002"

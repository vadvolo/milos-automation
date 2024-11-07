from annet.generators import PartialGenerator
from annet.storage import Device


class Shutdown(PartialGenerator):
    
    TAGS = ["shutdown", "iface"]
    
    def acl_cisco(self, device):
        return """
        interface
            shutdown
        """
    
    def run_cisco(self, device):
        for interface in device.interfaces:
            with self.block(f"interface {interface.name}"):
                yield "no shutdown"

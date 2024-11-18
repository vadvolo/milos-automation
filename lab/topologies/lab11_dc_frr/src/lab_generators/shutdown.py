from annet.generators import PartialGenerator
from annet.storage import Device


class Shutdown(PartialGenerator):
    """Partial generator class of interfaces' shutdown"""

    TAGS = ["shutdown", "iface"]

    def acl_cisco(self, _: Device):
        """ACL for Cisco devices"""

        return """
        interface
            shutdown
        """

    def run_cisco(self, device: Device):
        """Generator for Cisco devices"""

        for interface in device.interfaces:
            with self.block(f"interface {interface.name}"):
                yield "no shutdown"

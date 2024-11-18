from annet.generators import PartialGenerator
from annet.storage import Device


class Hostname(PartialGenerator):
    """Partial generator class of hostname"""

    TAGS = ["hostname"]

    def acl_cisco(self, _: Device):
        """ACL for Cisco devices"""

        return """
        hostname
        """

    def run_cisco(self, device: Device):
        """Generator for Cisco devices"""

        yield "hostname", device.name.split(".")[0]

    def acl_arista(self, _: Device):
        """ACL for Arista devices"""

        return """
        hostname
        """

    def run_arista(self, device: Device):
        """Generator for Arista devices"""

        yield "hostname", device.name.split(".")[0]

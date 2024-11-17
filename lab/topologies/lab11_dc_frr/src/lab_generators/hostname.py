from annet.generators import PartialGenerator
from annet.storage import Device


class Hostname(PartialGenerator):

    TAGS = ["hostname"]

    def acl_cisco(self, _: Device):
        return """
        hostname
        """

    def run_cisco(self, device: Device):
        yield "hostname", device.name.split(".")[0]

    def acl_arista(self, _: Device):
        return """
        hostname
        """

    def run_arista(self, device: Device):
        yield "hostname", device.name.split(".")[0]

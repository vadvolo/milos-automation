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

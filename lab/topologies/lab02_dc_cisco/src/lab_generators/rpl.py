from annet.generators import PartialGenerator
from annet.storage import Device


class RoutePolicy(PartialGenerator):
    
    TAGS = ["rpl", "routing"]
    
    def acl_cisco(self, _: Device):
        return """
        ip bgp-community new-format
        ip community-list
        route-map
            ~ %global=1
        """
    
    def run_cisco(self, device: Device):
        yield "ip bgp-community new-format"
        yield "ip community-list standard TOR_NETS permit 65000:1"
        yield "ip community-list standard GRACEFUL_SHUTDOWN permit 65535:0"

        if device.device_role.name == "ToR":
            yield """
route-map SPINE_IMPORT permit 10
  match community TOR_NETS GRACEFUL_SHUTDOWN
  set local-preference 0
route-map SPINE_IMPORT permit 20
  match community TOR_NETS
  set local-preference 100
route-map SPINE_IMPORT deny 9999
route-map SPINE_EXPORT permit 10
  match community TOR_NETS
route-map SPINE_EXPORT deny 9999
route-map CONNECTED permit 10
  match interface Loopback0
  set community 65000:1
route-map CONNECTED deny 9999
"""
        elif device.device_role.name == "Spine":
            yield """
route-map TOR_IMPORT permit 10
  match community TOR_NETS
route-map TOR_IMPORT deny 9999
"""
            with self.block("route-map TOR_EXPORT permit 10"):
                yield "match community TOR_NETS"
                if is_drained_device(device):
                    yield "set community 65535:0 additive"
            yield "route-map TOR_EXPORT deny 9999"


def is_drained_device(device: Device) -> bool:
    if "maintenance" in [tag.name for tag in device.tags]:
        return True
    return False

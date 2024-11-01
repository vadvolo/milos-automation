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

        if device.device_role.name == "ToR":
            yield """
ip community-list standard TOR_NETS permit 65000:1
route-map SPINE_IMPORT permit 10
  match community TOR_NETS
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
ip community-list standard TOR_NETS permit 65000:1
route-map TOR_IMPORT permit 10
  match community TOR_NETS
route-map TOR_IMPORT deny 9999
route-map TOR_EXPORT permit 10
  match community TOR_NETS
route-map TOR_EXPORT deny 9999
"""
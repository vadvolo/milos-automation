from annet.generators import PartialGenerator
from annet.storage import Device

from .helpers.router import is_drained_device


class RoutePolicy(PartialGenerator):
    """Partial generator class of routing policies"""

    TAGS = ["rpl", "routing"]

    def acl_cisco(self, _: Device):
        """ACL for Cisco devices"""

        return """
        ip bgp-community new-format
        ip community-list
        route-map
            ~ %global=1
        """

    def run_cisco(self, device: Device):
        """Generator for Cisco devices"""

        yield "ip bgp-community new-format"
        yield "ip community-list standard GSHUT permit 65535:0"
        yield "ip community-list standard TOR_NETS permit 65000:1"

        if device.device_role.name == "ToR":
            yield """
route-map TOR_IMPORT_SPINE permit 10
 match community GSHUT
 set local-preference 0
route-map TOR_IMPORT_SPINE permit 20
 set local-preference 100

route-map TOR_EXPORT_SPINE permit 10
 match community TOR_NETS
route-map TOR_EXPORT_SPINE deny 9999

route-map IMPORT_CONNECTED permit 10
 match interface Loopback0
 set community 65000:1
route-map IMPORT_CONNECTED deny 9999
"""
        elif device.device_role.name == "Spine":
            yield """
route-map SPINE_IMPORT_TOR permit 10
 match community TOR_NETS
route-map SPINE_IMPORT_TOR deny 9999
"""

            with self.block("route-map SPINE_EXPORT_TOR permit 10"):
                yield " match community TOR_NETS"
                if is_drained_device(device):
                    yield " set community 65535:0 additive"
            yield "route-map SPINE_EXPORT_TOR deny 9999"

    def acl_arista(self, _: Device):
        """ACL for Arista devices"""

        return """
        ip community-list
        route-map
            ~ %global=1
        """

    def run_arista(self, device: Device):
        """Generator for Arista devices"""

        if device.device_role.name == "Spine":
            yield "ip community-list GSHUT permit GSHUT"
            yield "ip community-list TOR_NETS permit 65000:1"

            with self.block("route-map SPINE_IMPORT_TOR permit 10"):
                yield "match community TOR_NETS"
            yield "route-map SPINE_IMPORT_TOR deny 9999"

            with self.block("route-map SPINE_EXPORT_TOR permit 10"):
                yield "match community TOR_NETS"
                if is_drained_device(device):
                    yield "set community 65535:0 additive"
            yield "route-map SPINE_EXPORT_TOR deny 9999"

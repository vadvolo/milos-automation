from annet.generators import PartialGenerator
from annet.storage import Device

class IfaceDescriptions(PartialGenerator):
    
    TAGS = ["description"]
    
    def acl_cisco(self, device):
        return """
        interface
            description
        """
    
    def run_cisco(self, device):
        for interface in device.interfaces:
            with self.block(f"interface {interface.name}"):
                yield "description MILOS002"
    

class Ifaces(PartialGenerator):
    TAGS = ["mgmt", "lldp"]

    def acl_routeros(self, _):
        return """
        ip
          neighbor              %cant_delete
            discovery-settings  %cant_delete
              set               %cant_delete

        """

    def run_routeros(self, _: Device):
        with self.block("ip"):
            with self.block("neighbor"):
                with self.block("discovery-settings"):
                    yield "set discover-interface-list=LAN lldp-med-net-policy-vlan=2 protocol=cdp"


class Ntp(PartialGenerator):
    TAGS = ["mgmt", "ntp"]

    def acl_routeros(self, _):
        return """
        system
          ntp         %cant_delete
            client    %cant_delete
              set     %cant_delete
              server  %cant_delete
                add   %cant_delete
        """

    #     /system ntp client
    # set enabled=yes
    # /system ntp client servers
    # add address=ntp0.ntp-servers.net
    def run_routeros(self, device: Device):
        with self.block("system"):
            with self.block("ntp"):
                with self.block("client"):
                    yield "set enabled=yes"

        with self.block("system"):
            with self.block("ntp"):
                with self.block("client"):
                    with self.block("server"):
                        yield "add address=ntp0.ntp-servers.net"

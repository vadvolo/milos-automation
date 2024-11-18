from typing import Optional

from annet.bgp_models import ASN
from annet.generators import PartialGenerator
from annet.mesh.executor import MeshExecutionResult
from annet.storage import Device

from .helpers.router import (
    AutonomusSystemIsNotDefined,
    bgp_asnum,
    bgp_groups,
    bgp_mesh,
    router_id,
)


class Bgp(PartialGenerator):
    """Partial generator class of BGP process and neighbors"""

    TAGS = ["bgp", "routing"]

    def acl_cisco(self, _: Device) -> str:
        """ACL for Cisco devices"""

        return """
        router bgp
            bgp
            neighbor
            redistribute connected
            maximum-paths
        """

    def run_cisco(self, device: Device):
        """Generator for Cisco devices"""

        mesh_data: MeshExecutionResult = bgp_mesh(device)
        rid: Optional[str] = router_id(mesh_data)
        try:
            asnum: Optional[ASN] = bgp_asnum(mesh_data)
        except AutonomusSystemIsNotDefined as err:
            raise RuntimeError(f"Device {device.name} has more than one defined autonomus system: {err}")

        if not asnum or not rid:
            return
        with self.block("router bgp", asnum):
            yield "bgp router-id", rid
            yield "bgp log-neighbor-changes"

            if device.device_role.name == "ToR":
                yield "maximum-paths 16"

            if mesh_data.global_options and mesh_data.global_options.ipv4_unicast and mesh_data.global_options.ipv4_unicast.redistributes:
                for redistribute in mesh_data.global_options.ipv4_unicast.redistributes:
                    yield "redistribute", redistribute.protocol, "" if not redistribute.policy else f"route-map {redistribute.policy}"

            for peer in mesh_data.peers:
                # define peer group attrs
                yield "neighbor", peer.group_name, "peer-group"

                if peer.import_policy:
                    yield "neighbor", peer.group_name, "route-map", peer.import_policy, "in"
                if peer.export_policy:
                    yield "neighbor", peer.group_name, "route-map", peer.export_policy, "out"
                if peer.options.soft_reconfiguration_inbound:
                    yield "neighbor", peer.group_name, "soft-reconfiguration inbound"
                if peer.options.send_community:
                    yield "neighbor", peer.group_name, "send-community both"

                # define peers specific attrs
                yield "neighbor", peer.addr, "peer-group", peer.group_name
                yield "neighbor", peer.addr, "remote-as", peer.remote_as

    def acl_arista(self, _: Device) -> str:
        """ACL for Arista devices"""

        return """
        router bgp
            router-id
            neighbor
            redistribute connected
            maximum-paths
            address-family
                neighbor
        """

    def run_arista(self, device: Device):
        """Generator for Arista devices"""

        mesh_data: MeshExecutionResult = bgp_mesh(device)
        rid: Optional[str] = router_id(mesh_data)
        try:
            asnum: Optional[ASN] = bgp_asnum(mesh_data)
        except AutonomusSystemIsNotDefined as err:
            raise RuntimeError(f"Device {device.name} has more than one defined autonomus system: {err}")

        if not asnum or not rid:
            return
        with self.block("router bgp", asnum):
            yield "router-id", rid

            for group in bgp_groups(mesh_data):
                yield "neighbor", group.group_name, "peer group"
                yield "neighbor", group.group_name, "route-map", group.import_policy, "in"
                yield "neighbor", group.group_name, "route-map", group.export_policy, "out"
                if group.send_community:
                    yield "neighbor", group.group_name, "send-community"

            for peer in mesh_data.peers:
                yield "neighbor", peer.addr, "peer group", peer.group_name
                yield "neighbor", peer.addr, "remote-as", peer.remote_as

            with self.block("address-family ipv4"):
                for group in bgp_groups(mesh_data):
                    yield "neighbor", group.group_name, "activate"

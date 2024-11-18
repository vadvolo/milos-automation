from dataclasses import dataclass
from typing import Optional

from annet.bgp_models import ASN
from annet.mesh import MeshExecutor
from annet.mesh.executor import MeshExecutionResult
from annet.storage import Device

from ..mesh_views import registry


@dataclass
class BGPGroup:
    """Represents a BGP group with specific attributes for configuration."""

    group_name: str
    remote_as: int  # Assumed that ASN is an integer, update if required
    import_policy: str
    export_policy: str
    send_community: bool = False

    # Define key fields as a class attribute
    _key_fields = ("group_name", "import_policy", "export_policy", "send_community")

    def __eq__(self, other):
        """
        Check equality of two BGPGroup instances.
        """
        if not isinstance(other, BGPGroup):
            return NotImplemented
        return all(getattr(self, attr) == getattr(other, attr) for attr in BGPGroup._key_fields)

    def __hash__(self):
        """
        Compute the hash value of the BGPGroup instance.
        """
        return hash(tuple(getattr(self, attr) for attr in BGPGroup._key_fields))


def bgp_mesh(device: Device) -> MeshExecutionResult:
    """Return mesh result for the devcice"""
    executor: MeshExecutor = MeshExecutor(registry, device.storage)
    mesh_data: MeshExecutionResult = executor.execute_for(device)
    return mesh_data


def bgp_asnum(mesh_data: MeshExecutionResult) -> Optional[ASN]:
    """Return AS number parse mesh bgp peers"""
    if not mesh_data:
        return None

    # AS can be defined in global options
    if mesh_data.global_options.local_as:
        return mesh_data.global_options.local_as

    # If AS is not defined in global options, searing it in peers
    asnum: set[ASN] = set()
    for peer in mesh_data.peers:
        asnum.add(peer.options.local_as)

    if len(asnum) == 1:
        return asnum.pop()
    elif len(asnum) > 1:
        raise AutonomusSystemIsNotDefined(str(asnum))

    return None


def router_id(mesh_data: MeshExecutionResult) -> Optional[str]:
    """Return router id for the device"""
    if mesh_data.global_options.router_id:
        return mesh_data.global_options.router_id
    return None


def bgp_groups(mesh_data: MeshExecutionResult) -> list[BGPGroup]:
    """Return list of BGP groups"""
    groups: set[BGPGroup] = set()
    for peer in mesh_data.peers:
        groups.add(BGPGroup(
            group_name=peer.group_name,
            remote_as=peer.remote_as,
            import_policy=peer.import_policy,
            export_policy=peer.export_policy,
            send_community=peer.options.send_community,
        ))
    return list(groups)


def is_drained_device(device: Device) -> bool:
    """Definition of devices in maintenance mode"""

    if "maintenance" in [tag.name for tag in device.tags]:
        return True
    return False


class AutonomusSystemIsNotDefined(Exception):
    """Autonomus system is not defined exception"""

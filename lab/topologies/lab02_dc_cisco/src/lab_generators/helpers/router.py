from typing import Optional

from annet.adapters.netbox.common.models import IpAddress
from annet.bgp_models import ASN
from annet.mesh import MeshExecutor
from annet.mesh.executor import MeshExecutionResult
from annet.storage import Device
from ..mesh_views import registry


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


def router_id(mesh_data: MeshExecutionResult) ->Optional[str]:
    """Return router id for the device"""
    if mesh_data.global_options.router_id:
        return mesh_data.global_options.router_id
    return None


def deduplicate_ip_addr(ip_addresses: list[IpAddress]) -> list[IpAddress]:
    """Return deduplicated ip addresses"""
    data_set: set[IpAddress] = set()
    for ip_addr in ip_addresses:
        data_set.add(ip_addr)
    return list(data_set)


class AutonomusSystemIsNotDefined(Exception):
    pass
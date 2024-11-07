import re

from annet.generators import PartialGenerator
from annet.storage import Device
from dataclasses import dataclass

@dataclass
class BgpPeer:
    addr: str
    asnum: str


def parse_device_name(device_name: str) -> dict[str, int]:
    """Parse device name to get """
    re_tor = r"^tor\-(?P<pod>\d+)\-(?P<num>\d+)"
    re_spine = r"^spine\-(?P<pod>\d+)\-(?P<plane>\d+)"

    m = re.match(re_tor, device_name)
    if m:
        tor_attr: dict[str, str] = m.groupdict()
        return {
            "pod": int(tor_attr["pod"]),
            "num": int(tor_attr["num"]),
        }

    m = re.match(re_spine, device_name)
    if m:
        spine_attr: dict[str, str] = m.groupdict()
        return {
            "pod": int(spine_attr["pod"]),
            "plane": int(spine_attr["plane"]),
        }
    raise RuntimeError(f"Could parse name '{device_name}'")


def get_asnum(device: Device) -> str:
    """Return AS number follow by role and name of the device"""
    if device.device_role.name == "Spine":
        parsed_name_spine: dict[str, int] = parse_device_name(device.name)
        spine_pod: int = parsed_name_spine["pod"]
        return f"6520{spine_pod}"
    elif device.device_role.name == "ToR":
        parsed_name_tor: dict[str, int] = parse_device_name(device.name)
        tor_pod: int = parsed_name_tor["pod"]
        tor_num: int = parsed_name_tor["num"]
        return f"651{tor_pod}{tor_num}"
    elif device.device_role.name == "Unknown":
        return ""
    else:
        raise RuntimeError(f"Unknown device role '{device.device_role.name}' of '{device.name}'")


def get_rid(device: Device) -> str:
    """Return Router ID follow by role and name of the device"""
    if device.device_role.name == "Spine":
        parsed_name_spine: dict[str, int] = parse_device_name(device.name)
        spine_plane: int = parsed_name_spine["plane"]
        spine_pod: int = parsed_name_spine["pod"]
        return f"1.2.{spine_pod}.{spine_plane}"
    elif device.device_role.name == "ToR":
        parsed_name_tor: dict[str, int] = parse_device_name(device.name)
        tor_pod: int = parsed_name_tor["pod"]
        tor_num: int = parsed_name_tor["num"]
        return f"1.1.{tor_pod}.{tor_num}"
    elif device.device_role.name == "Unknown":
        return ""
    else:
        raise RuntimeError(f"Unknown device role '{device.device_role.name}' of '{device.name}'")


def bgp_peers(device: Device) -> list[BgpPeer]:
    """Return list of BGP peers"""

    res: list[BgpPeer] = []
    if device.device_role.name == "Spine":
        parsed_name_spine: dict[str, int] = parse_device_name(device.name)
        spine_plane: int = parsed_name_spine["plane"]
        for remote_device in device.neighbours:
            if remote_device.device_role.name == "ToR":
                parsed_name_tor: dict[str, int] = parse_device_name(remote_device.name)
                tor_num: int = parsed_name_tor["num"]
                res.append(
                    BgpPeer(
                        addr=f"192.168.{spine_plane}{tor_num}.2",
                        asnum=get_asnum(remote_device),
                    )
                )

    elif device.device_role.name == "ToR":
        parsed_name_tor: dict[str, int] = parse_device_name(device.name)
        tor_num: int = parsed_name_tor["num"]
        for remote_device in device.neighbours:
            if remote_device.device_role.name == "Spine":
                parsed_name_spine: dict[str, int] = parse_device_name(remote_device.name)
                spine_plane: int = parsed_name_spine["plane"]
                res.append(
                    BgpPeer(
                        addr=f"192.168.{spine_plane}{tor_num}.1",
                        asnum=get_asnum(remote_device),
                    )
                )

    return res


class Bgp(PartialGenerator):
    
    TAGS = ["bgp", "routing"]
    
    def acl_cisco(self, _: Device) -> str:
        return """
        router bgp
            bgp
            neighbor
            redistribute connected
            maximum-paths
        """
    
    def run_cisco(self, device: Device):
        asnum: str = get_asnum(device)
        rid: str = get_rid(device)
        if not asnum or not rid:
            return
        with self.block("router bgp", asnum):
            yield "bgp router-id", rid
            yield "bgp log-neighbor-changes"
            if device.device_role.name == "Spine":
                yield "neighbor TOR peer-group"
                yield "neighbor TOR route-map TOR_IMPORT in"
                yield "neighbor TOR route-map TOR_EXPORT out"
                yield "neighbor TOR soft-reconfiguration inbound"
                yield "neighbor TOR send-community both"

                for peer in bgp_peers(device):
                    yield "neighbor", peer.addr, "peer-group TOR"
                    yield "neighbor", peer.addr, "remote-as", peer.asnum

            if device.device_role.name == "ToR":
                yield "redistribute connected route-map CONNECTED"
                yield "maximum-paths 16"
                yield "neighbor SPINE peer-group"
                yield "neighbor SPINE route-map SPINE_IMPORT in"
                yield "neighbor SPINE route-map SPINE_EXPORT out"
                yield "neighbor SPINE soft-reconfiguration inbound"
                yield "neighbor SPINE send-community both"

                for peer in bgp_peers(device):
                    yield "neighbor", peer.addr, "remote-as", peer.asnum
                    yield "neighbor", peer.addr, "peer-group SPINE"

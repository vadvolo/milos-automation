from annet.bgp_models import Redistribute
from annet.mesh import DirectPeer, GlobalOptions, MeshRulesRegistry, MeshSession


registry = MeshRulesRegistry(match_short_name=True)

BASE_ASNUM = 65000


@registry.device("tor-{pod}-{num}")
def global_options(global_opts: GlobalOptions):
    """Define global options of ToR switches"""

    global_opts.router_id = f"1.1.{global_opts.match.pod}.{global_opts.match.num}"
    global_opts.ipv4_unicast.redistributes = (
        Redistribute(protocol="connected", policy="IMPORT_CONNECTED"),
    )


# pylint: disable=unused-argument
@registry.direct("tor-{pod}-{num}", "spine-{pod}-{plane}")
def tor_to_spine(tor: DirectPeer, spine: DirectPeer, session: MeshSession):
    """Define peering between Spines and ToRs for IPv4 unicast family"""

    tor.asnum = BASE_ASNUM + 100 + tor.match.pod * 10 + tor.match.num
    tor.addr = f"10.{spine.match.plane}.{tor.match.num}.12/24"
    tor.families = ["ipv4_unicast"]
    tor.group_name = "TOR"
    tor.import_policy = "TOR_IMPORT_SPINE"
    tor.export_policy = "TOR_EXPORT_SPINE"
    tor.send_community = True
    tor.soft_reconfiguration_inbound = True

    spine.asnum = BASE_ASNUM + 200 + spine.match.pod
    spine.addr = f"10.{spine.match.plane}.{tor.match.num}.11/24"
    spine.families = ["ipv4_unicast"]
    spine.group_name = "SPINE"
    spine.import_policy = "SPINE_IMPORT_TOR"
    spine.export_policy = "SPINE_EXPORT_TOR"
    spine.send_community = True
    spine.soft_reconfiguration_inbound = True

from annet.mesh import MeshRulesRegistry

from . import spine, tor


registry = MeshRulesRegistry(match_short_name=True)
registry.include(tor.registry)
registry.include(spine.registry)

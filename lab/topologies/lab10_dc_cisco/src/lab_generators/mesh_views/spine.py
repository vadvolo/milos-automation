from annet.mesh import GlobalOptions, MeshRulesRegistry


registry = MeshRulesRegistry(match_short_name=True)


@registry.device("spine-{pod}-{plane}")
def global_options(global_opts: GlobalOptions):
    """Define global options of Spine switches"""

    global_opts.router_id = f"1.2.{global_opts.match.pod}.{global_opts.match.plane}"

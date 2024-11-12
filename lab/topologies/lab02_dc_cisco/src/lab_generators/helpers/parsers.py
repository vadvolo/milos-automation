import re

def parse_device_name(device_name: str) -> dict[str, int]:
    """Parse device name to get pod, num and plane from the name"""
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

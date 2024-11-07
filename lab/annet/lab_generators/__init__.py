from annet.generators import BaseGenerator
from annet.storage import Storage

from . import interfaces


def get_generators(store: Storage) -> list[BaseGenerator]:
    return [
        interfaces.Ifaces(store),
        interfaces.Ntp(store),
        interfaces.IfaceDescriptions(store),
    ]

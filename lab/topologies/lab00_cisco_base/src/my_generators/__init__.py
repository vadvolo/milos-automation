from annet.generators import BaseGenerator
from annet.storage import Storage

from . import interfaces


def get_generators(store: Storage) -> list[BaseGenerator]:
    return [
        interfaces.IfaceDescriptions(store),
        interfaces.IfaceMtu(store),
    ]

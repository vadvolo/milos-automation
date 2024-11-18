from annet.generators import BaseGenerator
from annet.storage import Storage

from . import interfaces


def get_generators(store: Storage) -> list[BaseGenerator]:
    """All the generators should be returned by the function"""

    return [
        interfaces.IfaceDescriptions(store),
        interfaces.IfaceMtu(store),
    ]

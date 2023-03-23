from typing import List
from prompt_toolkit.completion import Completer, Completion, NestedCompleter
from cli.command import COMMANDS


def __none_values(d: dict):
    """
    Recursively convert all values to None, used for generating a nested completer
    """
    res = {}
    for k, v in d.items():
        if type(v) == dict:
            res[k] = __none_values(v)
        else:
            res[k] = None
    return res

FullCompleter = NestedCompleter.from_nested_dict(__none_values(COMMANDS))
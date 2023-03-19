from typing import List
from prompt_toolkit.completion import Completer, Completion, NestedCompleter

COMMANDS = {
    'help': 'Show this help message and exit',
    'show': {
        'implants': 'Show all implants',
        'operators': 'Show all operators',
        'tasks': 'Show all active tasks',
        'stats': 'Show basic statistics about the server and clients'
    },
    'interact': "<implant_id>",
    "exit": "Exit the client",

}

def __none_values(d: dict):
    # Recursively convert all values to None
    res = {}
    for k, v in d.items():
        if type(v) == dict:
            res[k] = __none_values(v)
        else:
            res[k] = None
    return res

FullCompleter = NestedCompleter.from_nested_dict(__none_values(COMMANDS))
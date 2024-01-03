from typing import Any, Dict, Iterable, List
from prompt_toolkit.completion import  NestedCompleter
from cli.commands import *
from config import OperatorConfig
import comms


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

InteractCompleter = NestedCompleter.from_nested_dict(__none_values(INTERACT_COMMANDS))


def sort_tasks_by_created_at(tasks: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    return sorted(tasks, key=lambda t: t.get("created_at"))  # type: ignore


def filter_by_implant(
    tasks: Iterable[Dict[str, Any]], implant_id: str
) -> List[Dict[str, Any]]:
    return list(filter(lambda t: t.get("implant_id") == implant_id, tasks))


def get_home_dynamic_completer(config: OperatorConfig) -> NestedCompleter:
    comms.ensure_token(config)
    current_implants = comms.list_implants(config)
    tasks = comms.get_tasks(config)

    completer_dict = {
        "help": None,
        "show": {"implants", "tasks", "stats"},
        "builder": {
            "set": {
                "initial_sleep_seconds",
                "schtask_persist",
                "use_antidebug",
                "kill_parent",
                "use_antivm",
                "scheduled_task_name",
                "register_max_retries",
            },
            "show": {
                "all",
                "initial_sleep_seconds",
                "schtask_persist",
                "use_antidebug",
                "kill_parent",
                "use_antivm",
                "scheduled_task_name",
                "register_max_retries",
            },
        },
        "broadcast": None,
        "build": None,
        "interact": {i.get("implant_id", "") for i in current_implants},
        "results": set(
            map(lambda t: t.get("task_id", ""), sort_tasks_by_created_at(tasks))
        ),
        "alias": {
            "set": {i.get("implant_id", "") for i in current_implants},
            "list": {i.get("implant_id", "") for i in current_implants},
            "delete": {i.get("implant_id", "") for i in current_implants},
        },
        "exit": None,
    }
    return NestedCompleter.from_nested_dict(completer_dict)


def get_interact_dynamic_completer(
    config: OperatorConfig, implant_id: str
) -> NestedCompleter:
    comms.ensure_token(config)
    tasks = comms.get_tasks(config)

    completer_dict = {
        "cmd": None,
        "download": None,
        "disable_defender": None,
        "upload": None,
        "inject": None,
        "sysinfo": None,
        "sleep": None,
        "selfdestruct": None,
        "chdir": None,
        "pwd": None,
        "getenv": None,
        "ls": None,
        "ps": None,
        "whoami": None,
        "clipboard": {
            "set": {"text"},
            "show": None,
        },
        "config": {
            "set": {
                "user_agent",
                "sleep_time",
                "kill_date",
                "jitter",
                "max_retries",
                "auto_self_destruct",
                "retry_wait",
                "retry_jitter",
                "tailoring_hash_function",
                "tailoring_hash_rounds",
            },
            "show": {
                "all",
                "user_agent",
                "sleep_time",
                "kill_date",
                "jitter",
                "max_retries",
                "auto_self_destruct",
                "retry_wait",
                "retry_jitter",
                "tailoring_hash_function",
                "tailoring_hash_rounds",
            },
        },
        "results": set(
            map(
                lambda t: t.get("task_id", ""),
                filter_by_implant(sort_tasks_by_created_at(tasks), implant_id),
            )
        ),
        "back": None,
        "exit": None,
    }
    return NestedCompleter.from_nested_dict(completer_dict)

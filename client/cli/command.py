from dataclasses import dataclass
from typing import Callable, Dict, List, Optional, Union
from tabulate import tabulate
import cli.interact

from comms import get_server_stats, list_implants, get_tasks, implant_exists
from config import OperatorConfig
from .logging import StyledLogger, get_styled_logger
from .commands import *
from opcodes import Opcodes

logger = get_styled_logger()


def handle(cmd: str, args: List[str], config: OperatorConfig) -> None:
    """
    Handle a command
    """
    if cmd in ["help", "h", "?"]:
        handle_help(args[0] if len(args) > 0 else None)
    elif cmd == "show":
        handle_show(args[0] if len(args) > 0 else None, config)
    elif cmd == "interact":
        handle_interact(args[0] if len(args) > 0 else None, config)
    elif cmd == "exit":
        handle_exit()
    elif cmd == "clear":
        print("\033c")
    elif cmd == "":
        return
    else:
        logger.error(f"Command {cmd} not found")


def handle_help(args: Optional[str]) -> None:
    """
    Handle the help command
    """
    if args is None:
        print("Available commands:\n")
        for k, v in COMMANDS.items():
            if type(v) == str:
                print(f"  {k}: {v}")
            else:
                print(f"  {k}:")
                for k2, v2 in v.items():
                    print(f"    - {k2}: {v2}")
        print()
    else:
        if args in COMMANDS:
            print(f"{args}: {COMMANDS[args]}")
        else:
            logger.error(f"Command {args} not found")


def handle_show(args: Optional[str], config: OperatorConfig) -> None:
    """
    Handle the show command
    """
    if args is None:
        print("Available commands:")
        for k, v in COMMANDS["show"].items():
            print(f"{k}: {v}")
    else:
        if args in COMMANDS["show"]:
            if args == "implants":
                show_implants(config)
            elif args == "tasks":
                show_tasks(config)
            elif args == "stats":
                show_stats(config)
        else:
            logger.error(f"Command {args} not found")


def handle_interact(implant_id: Optional[str], config: OperatorConfig) -> None:
    """
    Handle the interact command
    """
    if implant_id is None:
        print("Please specify an implant ID")
    else:
        if implant_exists(config, implant_id):
            print(f"Interacting with implant {implant_id}...")
            cli.interact.interact_prompt(config, implant_id)


def handle_exit() -> None:
    """
    Handle the exit command
    """
    print("Exiting...")
    exit(0)


############
# Functions to show the implants, operators, tasks, and stats
############


def show_implants(config: OperatorConfig) -> None:
    """
    Show all the implants
    """

    implants = list_implants(config)
    if implants is None:
        logger.error("Failed to get implants")
        return

    minified = []
    for implant in implants:
        minified.append(
            [
                implant["implant_id"][0:8],
                implant["hostname"],
                implant["ip"],
                implant["os"],
                implant["last_seen"],
                implant["created_at"],
            ]
        )

    print(
        tabulate(
            minified,
            headers=["ID", "Hostname", "IP", "OS", "Last Seen", "First Seen"],
            tablefmt="fancy_grid",
        )
    )


def show_stats(config: OperatorConfig) -> None:
    """
    Show the server stats
    """
    stats = get_server_stats(config)
    if stats is None:
        logger.error("Failed to get server stats")
        return

    # Headers are keys, values are in one row
    print(
        tabulate(
            [list(stats.keys()), list(stats.values())],
            headers="firstrow",
            tablefmt="fancy_grid",
        )
    )


def show_tasks(config: OperatorConfig) -> None:
    """
    Show the active tasks that belong to this operator
    """
    tasks = get_tasks(config)

    if tasks is None:
        logger.error("Failed to get tasks")
        return

    if len(tasks) == 0:
        logger.warning("No tasks found")
        return

    minified = []
    for job in tasks:
        if job["status"] != "COMPLETE":
            minified.append(
                [
                    job["task_id"][0:8],
                    Opcodes.get_by_value(job["opcode"]),
                    job["status"],
                    job["implant_id"][0:8],
                    job["args"],
                    job["created_at"],
                ]
            )

    print(
        tabulate(
            minified,
            headers=["ID", "Opcode", "Status", "Implant ID", "Args", "Created At"],
            tablefmt="fancy_grid",
        )
    )

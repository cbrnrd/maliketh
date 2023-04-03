COMMANDS = {
    "help": "Show this help message and exit",
    "show": {
        "implants": "Show all implants",
        "tasks": "Show all active tasks",
        "stats": "Show basic statistics about the server and clients",
    },
    "interact <implant_id>": "Interact with a given implant id",
    "exit": "Exit the client",
}

INTERACT_COMMANDS = {
    "cmd <args>": "Execute a shell command on the implant",
    "download <remote_path> <local_path>": "Download a file from the implant",
    "upload <local_path>": "Upload a file to the implant",
    "inject <local_dll>": "Inject a DLL into the implant",
    "sysinfo": "Get system information from the implant",
    "sleep <seconds>": "Sleep for a given number of seconds",
    "selfdestruct": "Remove and kill the implant",
    "config": {
        "set": {
            "kill_date <YYYY-MM-DD>": "Set the kill date for the implant",
            "jitter <percent>": "Set the jitter for the implant",
            "max_retries <number>": "Set the max retries for the implant",
        },
        "show": {
            "all": "Show all configuration options",
            "kill_date": "Show the kill date for the implant",
            "jitter": "Show the jitter for the implant",
            "max_retries": "Show the max retries for the implant",
        }
    },
    "exit": "Exit the interact menu",

}

def walk_dict(d,depth=0):
    for k,v in sorted(d.items(),key=lambda x: x[0]):
        if isinstance(v, dict):
            print("  "*depth + ("- %s" % k))
            walk_dict(v,depth+1)
        else:
            print("  "*depth + "- %s: %s" % (k, v))

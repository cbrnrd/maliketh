from enum import Enum

class Opcodes(Enum):
    CMD = 1           # Execute a command via cli
    SELFDESTRUCT = 2  # Self destruct the agent (delete self and shut down)
    SYSINFO = 3       # Get system information
    SLEEP = 4         # Sleep for a given amount of time

    def __str__(self):
        return self.name



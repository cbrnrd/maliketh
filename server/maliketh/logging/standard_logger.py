from io import StringIO
from typing import List, TextIO, Union
import colorama
from dataclasses import dataclass
from enum import Enum


class LogLevel(Enum):
    DEBUG = 0
    INFO = 1
    OK = 2
    WARNING = 3
    ERROR = 4
    CRITICAL = 5

    @staticmethod
    def get_names() -> List[str]:
        return [name for name, member in LogLevel.__members__.items()]

    def get_prelog(self):
        if self == LogLevel.DEBUG:
            return f"[{colorama.Fore.MAGENTA}?{colorama.Style.RESET_ALL}] "
        elif self == LogLevel.INFO:
            return f"[{colorama.Fore.BLUE}*{colorama.Style.RESET_ALL}] "
        elif self == LogLevel.OK:
            return f"[{colorama.Fore.GREEN}+{colorama.Style.RESET_ALL}] "
        elif self == LogLevel.WARNING:
            return f"[{colorama.Fore.YELLOW}!{colorama.Style.RESET_ALL}] "
        elif self == LogLevel.ERROR:
            return f"[{colorama.Fore.RED}!{colorama.Style.RESET_ALL}] "
        elif self == LogLevel.CRITICAL:
            return f"[{colorama.Fore.WHITE}{colorama.Back.RED}!!!{colorama.Style.RESET_ALL}]"


@dataclass
class StandardLogger:
    out_file: Union[StringIO, TextIO]  # The file or stream to write to
    err_file: Union[
        StringIO, TextIO
    ]  # The file or stream to write errors to (will write to out_file as well)
    level: LogLevel  # The minimum level to log

    def log(self, level: LogLevel, msg: str):
        if level.value < self.level.value:
            return
        self.out_file.write(f"{level.get_prelog()}{msg}\n")
        if level.value >= LogLevel.WARNING.value:
            self.err_file.write(f"{level.get_prelog()}{msg}\n")

    def debug(self, msg: str):
        self.log(LogLevel.DEBUG, msg)

    def info(self, msg: str):
        self.log(LogLevel.INFO, msg)

    def ok(self, msg: str):
        self.log(LogLevel.OK, msg)

    def warning(self, msg: str):
        self.log(LogLevel.WARNING, msg)

    def error(self, msg: str):
        self.log(LogLevel.ERROR, msg)

    def critical(self, msg: str):
        self.log(LogLevel.CRITICAL, msg)

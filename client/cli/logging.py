import abc
from dataclasses import dataclass
from prompt_toolkit.styles import Style
from prompt_toolkit.formatted_text import FormattedText
from prompt_toolkit import print_formatted_text
from enum import Enum
import config


class LogLevel(Enum):
    DEBUG = 0
    INFO = 1
    OK = 2
    WARNING = 3
    ERROR = 4
    CRITICAL = 5

    def get_icon(self) -> str:
        """
        Returns a string representation of the log level.
        """
        if self == LogLevel.DEBUG:
            return "?"
        elif self == LogLevel.INFO:
            return "+"
        elif self == LogLevel.OK:
            return "+"
        elif self == LogLevel.WARNING:
            return "!"
        elif self == LogLevel.ERROR:
            return "x"
        elif self == LogLevel.CRITICAL:
            return "xxx"

        return "?"

    def to_lower(self) -> str:
        return self.name.lower()


@dataclass
class StyledLogger:
    level: LogLevel = LogLevel.INFO
    _style: Style = Style.from_dict(
        {
            "debug": "#ffaa00",  # Orange
            "ok": "#00ff00",  # Green
            "info": "#0000ff",  # blue
            "warning": "#ffff00",
            "error": "#ff0000",
            "critical": "#ff0000 bold",
        }
    )

    def log(self, level: LogLevel, msg: str, *args, **kwargs):
        if level.value < self.level.value:
            return

        formatted = msg % args
        to_print = FormattedText(
            [
                ("", "["),
                (f"class:{level.to_lower()}", level.get_icon()),
                ("", "] "),
                ("", formatted),
            ]
        )
        print_formatted_text(to_print, style=self._style)

    def debug(self, msg: str, *args, **kwargs):
        self.log(LogLevel.DEBUG, msg, *args, **kwargs)

    def info(self, msg: str, *args, **kwargs):
        self.log(LogLevel.INFO, msg, *args, **kwargs)

    def ok(self, msg: str, *args, **kwargs):
        self.log(LogLevel.OK, msg, *args, **kwargs)

    def warning(self, msg: str, *args, **kwargs):
        self.log(LogLevel.WARNING, msg, *args, **kwargs)

    def error(self, msg: str, *args, **kwargs):
        self.log(LogLevel.ERROR, msg, *args, **kwargs)

    def critical(self, msg: str, *args, **kwargs):
        self.log(LogLevel.CRITICAL, msg, *args, **kwargs)


def get_styled_logger() -> StyledLogger:
    return StyledLogger(config.log_level)

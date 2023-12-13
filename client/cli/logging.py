from dataclasses import dataclass
from datetime import datetime
import logging
import logging.config
from prompt_toolkit.styles import Style
from prompt_toolkit.formatted_text import FormattedText
from prompt_toolkit import print_formatted_text
from prompt_toolkit.patch_stdout import patch_stdout
from enum import Enum
import structlog

RESET_ALL = "\033[0m"
BRIGHT = "\033[1m"
DIM = "\033[2m"
RED = "\033[31m"
BLUE = "\033[34m"
CYAN = "\033[36m"
MAGENTA = "\033[35m"
YELLOW = "\033[33m"
GREEN = "\033[32m"
RED_BACK = "\033[41m"


class StructlogStyle:
    reset = RESET_ALL
    bright = BRIGHT
    level_critical = RED
    level_exception = RED
    level_error = RED
    level_warn = YELLOW
    level_info = CYAN
    level_debug = GREEN
    level_notset = RED_BACK
    timestamp = DIM
    logger_name = BLUE
    kv_key = CYAN
    kv_value = MAGENTA


style_dict = {
    "critical": StructlogStyle.level_critical,
    "exception": StructlogStyle.level_exception,
    "error": StructlogStyle.level_error,
    "warn": StructlogStyle.level_warn,
    "warning": StructlogStyle.level_warn,
    "info": StructlogStyle.level_info,
    "debug": StructlogStyle.level_debug,
    "notset": StructlogStyle.level_notset,
    "ok": StructlogStyle.level_debug
}

def ok(self, msg, *args, **kwargs):
    return self.log(21, msg, *args, **kwargs)


def setup_structlog(level: int, with_timestamps: bool):
    
    # Set up our `OK` log level.
    OK = 21
    structlog.stdlib.OK = OK
    structlog.stdlib._NAME_TO_LEVEL["ok"] = OK
    structlog.stdlib._LEVEL_TO_NAME[OK] = "ok"
    structlog.stdlib._FixedFindCallerLogger.ok = ok
    structlog.stdlib.BoundLogger.ok = ok

    timestamper = structlog.processors.TimeStamper(fmt="iso")
    pre_chain = [
        # Add the log level and a timestamp to the event_dict if the log entry is not from structlog.
        structlog.stdlib.add_log_level,
        timestamper,
    ]

    processors = [
        structlog.stdlib.filter_by_level,
        structlog.stdlib.add_log_level,
        structlog.stdlib.PositionalArgumentsFormatter(),
        structlog.processors.StackInfoRenderer(),
        structlog.stdlib.ProcessorFormatter.wrap_for_formatter,
    ]

    if with_timestamps:
        processors.insert(2, timestamper)

    structlog.configure(
        processors=processors,
        context_class=structlog.threadlocal.wrap_dict(dict),
        logger_factory=structlog.stdlib.LoggerFactory(),
        wrapper_class=structlog.stdlib.BoundLogger,
        cache_logger_on_first_use=True,
    )

    logging.config.dictConfig(
        {
            "version": 1,
            "disable_existing_loggers": False,
            "formatters": {
                "console": {
                    "()": structlog.stdlib.ProcessorFormatter,
                    "processor": structlog.dev.ConsoleRenderer(
                        colors=True, level_styles=style_dict, exception_formatter=structlog.dev.rich_traceback),
                    "foreign_pre_chain": pre_chain,
                },
            },
            "handlers": {
                "development": {
                    "level": "DEBUG",
                    "class": "logging.StreamHandler",
                    "formatter": "console",
                },
                "production": {
                    "level": "DEBUG",
                    "class": "logging.StreamHandler",
                    "formatter": "console",
                },
            },
            "loggers": {
                "": {"handlers": ["production"], "level": "DEBUG", "propagate": True},
                "pika": {"handlers": ["production"], "level": "WARN", "propagate": True},
                "requests": {"handlers": ["production"], "level": "WARN", "propagate": True}
            },
        }
    )

    logging.addLevelName(OK, "OK")
    logging.getLogger().setLevel(level)


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
    level: LogLevel
    with_timestamps: bool = True
    _style: Style = Style.from_dict(
        {
            "dim": "#101010",
            "debug": "#ffaa00 bold",  # Orange
            "ok": "#00ff00 bold",  # Green
            "info": "#0000ff bold",  # blue
            "warning": "#ffff00 bold",
            "error": "#ff0000 bold",
            "critical": "#ff0000 bold",
        }
    )

    def log(self, level: LogLevel, msg: str, *args, **kwargs):
        if level.value < self.level.value:
            return

        formatted = msg % args
        to_print = FormattedText(
            [
                ("class:dim", f"{datetime.now().isoformat()}Z ") if self.with_timestamps else ("", ""),
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


def get_styled_logger(with_timestamps=False) -> StyledLogger:
    return StyledLogger(LogLevel.INFO, with_timestamps=with_timestamps)

from typing import Callable
from prompt_toolkit import print_formatted_text, HTML, PromptSession
from prompt_toolkit.shortcuts import prompt
from prompt_toolkit.auto_suggest import AutoSuggestFromHistory
from prompt_toolkit.styles import Style
from cli.logging import LogLevel, StyledLogger
from cli.command import handle
from config import OperatorConfig
from .completer import FullCompleter
from comms import get_server_stats

PROMPT_STYLE = Style.from_dict({
    
    # Home is red
    "home": "#ff0000 bold",
    "interact": "#00ff00 bold",  # Green
    "warning": "#ffff00 bold",  # Yellow
})

def bottom_bar(config: OperatorConfig) -> Callable:
    server_stats = get_server_stats(config)
    return lambda: HTML(f"User: <b bg='ansired'>{config.name}</b> | Implants: <b bg='ansired'>{server_stats['implants']}</b> | Operators: <b bg='ansired'>{server_stats['operators']}</b> | Uptime: <b bg='ansired'>{server_stats['uptime']}</b>")

def main_loop(config: OperatorConfig):



    logger = StyledLogger(LogLevel.INFO)

    while True:
        try:
            session = PromptSession(
                message=HTML(f"<warning>maliketh</warning> (<home>{config.name}</home>) > "),
                style=PROMPT_STYLE,
                enable_history_search=True,
                completer=FullCompleter,
                bottom_toolbar=bottom_bar(config),
                auto_suggest=AutoSuggestFromHistory(),
            )
            text = session.prompt()
            cmd, *args = text.split(" ")
            handle(cmd, args, config)

        except KeyboardInterrupt:
            continue
        except EOFError:
            break

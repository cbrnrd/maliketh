import sys
import threading
from typing import Callable
from prompt_toolkit import print_formatted_text, HTML, PromptSession
from prompt_toolkit.shortcuts import prompt
from prompt_toolkit.auto_suggest import AutoSuggestFromHistory
from cli.command import handle
from config import OperatorConfig
from .completer import FullCompleter, get_home_dynamic_completer
from comms import get_server_stats
from .style import PROMPT_STYLE
from .banner import get_full_banner


def bottom_bar(config: OperatorConfig) -> Callable:
    server_stats = get_server_stats(config)
    return lambda: HTML(
        f"User: <b bg='ansired'>{config.name}</b> | Implants: <b bg='ansired'>{server_stats['implants']}</b> | Operators: <b bg='ansired'>{server_stats['operators']}</b> | Uptime: <b bg='ansired'>{server_stats['uptime']}</b>"
    )


def main_loop(config: OperatorConfig):
    print_formatted_text(HTML(f"<home>{get_full_banner()}</home>"), style=PROMPT_STYLE)

    while True:
        try:
            session = PromptSession(
                message=HTML(
                    f"<warning>maliketh</warning> (<home>{config.name}</home>) > "
                ),
                style=PROMPT_STYLE,
                enable_history_search=True,
                completer=get_home_dynamic_completer(config),
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

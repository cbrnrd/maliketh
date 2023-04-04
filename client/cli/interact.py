from typing import List, Optional
from prompt_toolkit import print_formatted_text, HTML, PromptSession
from prompt_toolkit.shortcuts import prompt
from prompt_toolkit.auto_suggest import AutoSuggestFromHistory
from prompt_toolkit.styles import Style
from .style import PROMPT_STYLE
from config import OperatorConfig
from .completer import InteractCompleter
from .commands import INTERACT_COMMANDS, walk_dict
from .logging import get_styled_logger

logger = get_styled_logger()


def interact_prompt(config: OperatorConfig, implant_id: str):
    from .cli import bottom_bar

    session = PromptSession(
        message=HTML(
            f"<warning>maliketh</warning> (<home>{config.name}</home>) - <interact>{implant_id[0:8]}</interact> > "
        ),
        style=PROMPT_STYLE,
        enable_history_search=True,
        completer=InteractCompleter,
        bottom_toolbar=bottom_bar(config),
        auto_suggest=AutoSuggestFromHistory(),
    )

    while True:
        try:
            text = session.prompt()
            cmd, *args = text.split(" ")
            handle(cmd, args, config, implant_id)
        except KeyboardInterrupt:
            continue
        except EOFError:
            break


def handle(cmd: str, args: List[str], config: OperatorConfig, implant_id: str) -> None:
    """
    Handle a command
    """
    if cmd in ["help", "h", "?"]:
        handle_help(args[0] if len(args) > 0 else None)
    elif cmd == "exit":
        handle_exit(config)
    else:
        logger.error(f"Command {cmd} not found")


def handle_help(args: Optional[str]) -> None:
    """
    Handle the help command
    """
    if args is None:
        print("Available commands:\n")
        walk_dict(INTERACT_COMMANDS)

        print()
    else:
        if INTERACT_COMMANDS.get(args):
            walk_dict(INTERACT_COMMANDS[args])
        else:
            logger.error(f"Command {args} not found")


def handle_exit(config: OperatorConfig) -> None:
    """
    Handle the exit command, jump back to the main loop
    """
    from .cli import main_loop

    main_loop(config)

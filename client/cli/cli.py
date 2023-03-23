from prompt_toolkit import print_formatted_text, HTML, PromptSession
from prompt_toolkit.shortcuts import prompt
from prompt_toolkit.styles import Style
from cli.logging import LogLevel, StyledLogger
from cli.command import handle
from config import OperatorConfig
from .completer import FullCompleter

PROMPT_STYLE = Style.from_dict({
    
    # Home is red
    "home": "#ff0000 bold",
    "interact": "#00ff00 bold",  # Green
    "warning": "#ffff00 bold",  # Yellow
})

def main_loop(config: OperatorConfig):

    session = PromptSession(
        message=HTML(f"<warning>maliketh</warning> (<home>{config.name}</home>) > "),
        style=PROMPT_STYLE,
        enable_history_search=True,
        completer=FullCompleter,
    )

    logger = StyledLogger(LogLevel.INFO)

    while True:
        try:
            text = session.prompt()
            cmd, *args = text.split(" ")
            handle(cmd, args, config)

        except KeyboardInterrupt:
            continue
        except EOFError:
            break

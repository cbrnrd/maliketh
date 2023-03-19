from prompt_toolkit import print_formatted_text, HTML, PromptSession
from prompt_toolkit.shortcuts import prompt
from prompt_toolkit.styles import Style
from cli.logging import LogLevel, StyledLogger

from config import OperatorConfig
from .completer import FullCompleter

PROMPT_STYLE = Style.from_dict({
    
    # Home is red
    "home": "#ff0000 bold",
    "interact": "#00ff00 bold",  # Green
})

def main_loop(config: OperatorConfig):

    session = PromptSession(
        message=HTML(f"<home>{config.name}</home> > "),
        style=PROMPT_STYLE,
        enable_history_search=True,
        completer=FullCompleter,
    )

    logger = StyledLogger(LogLevel.INFO)

    while True:
        try:
            text = session.prompt()
        except KeyboardInterrupt:
            continue
        except EOFError:
            break
        else:
            logger.ok(text)

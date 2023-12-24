from prompt_toolkit.styles import Style

TABULATE_STYLE = "fancy_grid"

PROMPT_STYLE = Style.from_dict(
    {
        # Home is red
        "home": "#ff0000 bold",
        "interact": "#00ff00 bold",  # Green
        "warning": "#ffff00 bold",  # Yellow
    }
)

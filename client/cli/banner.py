import random


BANNER = """
                         )
                        (            ((        ( (
                        )   (          ))       ))
              />        /  ((            )//    /)
 (           //-------------------------------------\\
(*)OXOXOXOXO(*>======================================>
 (           \\\\-------------------------------------/
              \\>
           
"""

QUOTES = [
    "I shall not part with it again.",
    "Stay away from Destined Death.",
    "1337 h4x0r",
    "Praise the sun!",
]


def get_full_banner():
    return BANNER + f'\t\t"{random.choice(QUOTES)}"'

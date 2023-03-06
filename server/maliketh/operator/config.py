import json
import os

from maliketh.crypto.utils import random_string
from maliketh.crypto.ec import generate_b64_ecc_keypair, generate_b64_signing_keypair
from maliketh.config import CONFIG_DIR

from typing import Any, Dict, Optional

from maliketh.config import SERVER_PUB_KEY_PATH


def generate_config(name: str, outfile: Optional[str]) -> Dict[str, Any]:
    """
    Generate a new (generated) config for an operator with the given name

    @param name - the name of the operator
    @param outfile - the file to write the config to. If None, the config will be printed to stdout
    """
    keypair = generate_b64_ecc_keypair()
    signing_keypair = generate_b64_signing_keypair()
    with open(SERVER_PUB_KEY_PATH, 'r') as f:
      server_pub = f.read()
      config = {
          "name": name,
          "c2": "localhost",
          "c2_port": 5000,
          "login_secret": random_string(32),
          "secret": keypair[0],
          "public": keypair[1],
          "signing_key": signing_keypair[0],
          "verify_key": signing_keypair[1],
          "server_pub": server_pub
      }

    if outfile is not None:
        with open(outfile, "w") as f:
            json.dump(config, f, indent=4)
        return config
    else:
        return config


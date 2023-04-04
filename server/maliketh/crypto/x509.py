import os
from typing import Optional

"""
Generate a CA certificate and private key for use with mutual TLS
Note: We should probably use a crypto library for this, but for now use openssl
"""


def generate_ca_cert() -> tuple:

    # check if openssl is installed
    if os.system("command -v openssl") != 0:
        raise Exception("openssl is not installed")

    ca_key = os.popen("openssl genrsa 2048").read()
    ca_cert = os.popen(
        "openssl req -new -x509 -nodes -days 365 -key - -subj '/CN=maliketh CA'"
    ).read()
    return ca_key, ca_cert

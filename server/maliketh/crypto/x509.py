import os

"""
Generate a CA certificate and private key for use with mutual TLS
Note: We should probably use a crpto library for this, but for now use openssl
"""
def generate_ca_cert() -> tuple:

    # check if openssl is installed
    if os.system("command -v openssl") != 0:
        raise Exception("openssl is not installed")

    ca_key = os.popen("openssl genrsa 2048").read()
    ca_cert = os.popen("openssl req -new -x509 -nodes -days 365 -key - -subj '/CN=maliketh CA'").read()
    return ca_key, ca_cert

"""
Generate an ECC keypair for use with mutual TLS
"""
def generate_ecc_keypair() -> tuple:
    ecc_key = os.popen("openssl ecparam -genkey -name secp384r1").read()
    ecc_public_key = os.popen("openssl ec -in - -pubout").read()
    return ecc_key, ecc_public_key
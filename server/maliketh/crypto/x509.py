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
    ca_cert = os.popen("openssl req -new -x509 -nodes -days 365 -key - -subj '/CN=maliketh CA'").read()
    return ca_key, ca_cert

"""
Generate an ECC keypair for use with mutual TLS.

@param write_dir: If specified, write the keypair to the specified directory
"""
def generate_ecc_keypair(write_dir: Optional[str]=None) -> tuple:
    ecc_key = os.popen("openssl ecparam -genkey -name secp384r1").read()
    if write_dir:
        with open(os.path.join(write_dir, "ecc_key.pem"), "w") as f:
            f.write(ecc_key)

    if write_dir:
        ecc_public_key = os.popen("openssl ec -in {} -pubout".format(os.path.join(write_dir, "ecc_key.pem"))).read()
        with open(os.path.join(write_dir, "ecc_key.pub"), "w") as f:
            f.write(ecc_public_key)
    else:
        # Write ecc_key to temp file
        with open("/tmp/ecc_key.pem", "w") as f:
            f.write(ecc_key)
        ecc_public_key = os.popen("openssl ec -in {} -pubout".format("/tmp/ecc_key.pem")).read()
        os.remove("/tmp/ecc_key.pem")
    return ecc_key, ecc_public_key

"""
Check if an ECC keypair exists in the specified directory
"""
def check_ecc_keypair(dir: str) -> bool:
    if not os.path.exists(os.path.join(dir, "ecc_key.pem")):
        return False
    if not os.path.exists(os.path.join(dir, "ecc_key.pub")):
        return False
    return True

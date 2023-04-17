from typing import Tuple
from nacl.public import PrivateKey, PublicKey, Box
from nacl.signing import SigningKey, VerifyKey
from nacl.encoding import Base64Encoder, RawEncoder
from cryptography.hazmat.primitives.asymmetric.ed25519 import (
    Ed25519PrivateKey,
    Ed25519PublicKey,
)
from cryptography.hazmat.primitives.serialization import (
    Encoding,
    PrivateFormat,
    PublicFormat,
    NoEncryption,
)

from maliketh.crypto.ec import *
from maliketh.crypto.utils import random_hex, random_string

with open(SERVER_PRIV_KEY_PATH, "rb") as f:
    server_priv_key = f.read()

with open(SERVER_PUB_KEY_PATH, "rb") as f:
    server_pub_key = f.read()

server_priv_key = PrivateKey(server_priv_key, encoder=Base64Encoder)
server_pub_key = PublicKey(server_pub_key, encoder=Base64Encoder)

config = {
    "name": "admin",
    "c2": "localhost",
    "c2_port": 5000,
    "login_secret": "RtEl.usKG@+@Ey/8LU[eEqH`uVC*FC_M",
    "secret": "yalsrURU0oxU+qPBSF7vv41bNX7tG3Bb8zi4caXM/is=",
    "public": "7Yi3+TvHC7G8iq28p7Y2dn/0dF9XjGjRJ5K3mcHZJic=",
    "signing_key": "8jE0g9YP9zX6NlmhBqUrAXlf4ovzGIPi87leZYIJHkA=",
    "verify_key": "L3d1Rd5d7V3mAlLZc5I3yy/TGHXS9UuStk91hZorP9U=",
    "server_pub": "Wum5IVMIbhuj6817C8UCl5BHuRHpL3QlsnuiE4yMRzw=",
    "rmq_queue": "pK6K5{,LWp6OG-szU>Wg)M5i5^M&^Sxm"
}
admin_signing_key = SigningKey(
    config["signing_key"].encode("utf-8"), encoder=Base64Encoder
)
admin_vk = admin_signing_key.verify_key
admin_pk = PublicKey(config["public"].encode("utf-8"), encoder=Base64Encoder)
admin_secret = PrivateKey(config["secret"].encode("utf-8"), encoder=Base64Encoder)


message = config["login_secret"]

### Client side ###
signed = admin_signing_key.sign(message.encode("utf-8"))
print(f"Signed message: {signed.signature}")

box = Box(admin_secret, server_pub_key)
b64encrypted = encrypt(server_pub_key, admin_secret, signed, encoder=Base64Encoder)
print(f"Base64 Encrypted message: {b64encrypted}")

print("=" * 80)

### Server side ###
decrypted = decrypt(admin_pk, server_priv_key, b64encrypted, encoder=Base64Encoder)
print(f"Decrypted message: {decrypted}")
print(f"Decrypted matches original: {decrypted == bytes(signed)}")

verify = verify_signature(admin_vk, decrypted)
print(f"Verified message: {verify}")

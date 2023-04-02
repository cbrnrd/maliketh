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
    "login_secret": ">=BieisD0W.!,6smFXa_Em9o1[(ax5+`",
    "secret": "+lLdlHCJOb7QMoFi45pIwJ+h7lw3KGgwYWS6t1ngrZo=",
    "public": "mOWFpZ82fpD9VSfR/Xvyj76iRxkE2AeOl9Uu9oFJiFE=",
    "signing_key": "KJ6askmfJ7TArJLQOxQVG/M75iFxydsUgI/X3PlU7LA=",
    "verify_key": "RvwVJCLw6R086uzYXOF3N1W3i/OfpKXDxACWjc3/4/Y=",
    "server_pub": "5JMcHUe6vVLUDJNEw8oviucvecH8jbPj51+jZydPlHQ="
}
admin_signing_key = SigningKey(config["signing_key"].encode("utf-8"), encoder=Base64Encoder)
admin_vk = admin_signing_key.verify_key
admin_pk = PublicKey(config["public"].encode("utf-8"), encoder=Base64Encoder)
admin_secret = PrivateKey(config["secret"].encode("utf-8"), encoder=Base64Encoder)


message = config['login_secret']

### Client side ###
signed = admin_signing_key.sign(message.encode("utf-8"))
print(f'Signed message: {signed.signature}')

box = Box(admin_secret, server_pub_key)
b64encrypted = encrypt(server_pub_key, admin_secret, signed, encoder=Base64Encoder)
print(f'Base64 Encrypted message: {b64encrypted}')

print('=' * 80)

### Server side ###
decrypted = decrypt(admin_pk, server_priv_key, b64encrypted, encoder=Base64Encoder)
print(f'Decrypted message: {decrypted}')
print(f'Decrypted matches original: {decrypted == bytes(signed)}')

verify = verify_signature(admin_vk, decrypted)
print(f'Verified message: {verify}')


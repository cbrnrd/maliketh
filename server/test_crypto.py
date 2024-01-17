from typing import Tuple
from nacl.public import PrivateKey, PublicKey, Box
from nacl.signing import SigningKey
from nacl.encoding import Base64Encoder

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
    "login_secret": "BrlhYmq'I>Py*]+oOaw9e5?o1ALYTv43",
    "secret": "Y3YybGWeYwWzxzzz3c6jmxTbQwmUJBKY1EEIGfQvplk=",
    "public": "P1rTHOOwPetIVdrO00Tu6oJupAUHv+grY5srhayXdgo=",
    "signing_key": "4tL7xJ4iy73Lu+WWS2aDe9jKzfQsgOVkniTTq++pdVY=",
    "verify_key": "31uKxgpMfEoyRt8cQN9+vo6KjlXKogr4rKKZwXm+Uw0=",
    "server_pub": "N4rWJU3k+891Ce72iglugk0b+biOZ9zped/45728FTo=",
    "rmq_queue": "oUA6~wDTk1oNJR^B4@el3z]2w{XKNb`("
}
admin_signing_key = SigningKey(
    config["signing_key"].encode("utf-8"), encoder=Base64Encoder
)

print(f"signing priv key: {[int(v) for v in admin_signing_key.encode()]}")
print(f"signing verify key: {[int(v) for v in admin_signing_key.verify_key.encode()]}")
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

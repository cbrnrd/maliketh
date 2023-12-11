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
    "login_secret": "}W+%v#I<EH6P04Dg6nB@pUeFN7&+<~}p",
    "secret": "+nIhLhABpqlFTfj7tw6r/AmWyT+OiD5kjZFkrYyOAn8=",
    "public": "VuMyC4/Cv3vuv8EF9gYJPbcGsPix8Wg5EHXgI4/HYAc=",
    "signing_key": "me2uRpl8lG6K4TdOlZmG99ioO11PF/IyS3q960kImmE=",
    "verify_key": "a3NU5gVTTgO7Jz56l0wviq9IQAqYSm3Wl2bJcHWK1u4=",
    "server_pub": "1qlKnFqqzKoKpj1ULTMWBBTTPVVz6DKVD9icm9hg0Vk=",
    "rmq_queue": "MTDa1e<p9gg$#Jq&Dra'NaXqSu(3MEP9"
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

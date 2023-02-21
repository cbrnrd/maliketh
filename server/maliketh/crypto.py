
import base64
import os

"""
Generate a random AES key and return it as a base64 encoded string
"""
def generate_aes_key() -> str:
    return base64.b64encode(random_bytes(32)).decode("utf-8")

"""
Generate a random AES IV and return it as a base64 encoded string
"""
def generate_aes_iv() -> str:
    return base64.b64encode(random_bytes(16)).decode("utf-8")

def random_bytes(length: int) -> bytes:
    return os.urandom(length)

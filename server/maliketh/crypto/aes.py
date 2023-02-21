
import base64
import os
from Crypto.Cipher import AES

from maliketh.crypto.utils import pad, random_bytes


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

"""
Encrypt a given byte array using AES
"""
def encrypt_aes(data: bytes, key: bytes, iv: bytes) -> bytes:
    cipher = AES.new(key, AES.MODE_CBC, iv)
    return cipher.encrypt(pad(data))

"""
Decrypt a given byte array using AES
"""
def decrypt_aes(data: bytes, key: bytes, iv: bytes) -> bytes:
    cipher = AES.new(key, AES.MODE_CBC, iv)
    return cipher.decrypt(data)

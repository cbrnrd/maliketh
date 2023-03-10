import base64
import os
from Crypto.Cipher import AES

from maliketh.crypto.utils import pad, random_bytes


def generate_aes_key() -> str:
    """
    Generate a random AES key and return it as a base64 encoded string
    """
    return base64.b64encode(random_bytes(32)).decode("utf-8")


def generate_aes_iv() -> str:
    """
    Generate a random AES IV and return it as a base64 encoded string
    """
    return base64.b64encode(random_bytes(16)).decode("utf-8")


def encrypt_aes(data: bytes, key: bytes, iv: bytes) -> bytes:
    """
    Encrypt a given byte array using AES
    """
    cipher = AES.new(key, AES.MODE_CBC, iv)
    return cipher.encrypt(pad(data))


def decrypt_aes(data: bytes, key: bytes, iv: bytes) -> bytes:
    """
    Decrypt a given byte array using AES
    """
    cipher = AES.new(key, AES.MODE_CBC, iv)
    return cipher.decrypt(data)

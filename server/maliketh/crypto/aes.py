import base64
import os
from Crypto.Cipher import AES
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

from maliketh.crypto.utils import pad, random_bytes


class GCM:
    @staticmethod
    def encrypt(data: bytes, aad: bytes, key: bytes) -> bytes:
        """
        Encrypt a given byte array using AES-GCM
        """
        aesgcm = AESGCM(key)
        nonce = os.urandom(12)
        return aesgcm.encrypt(nonce, data, aad) + nonce

    @staticmethod
    def encrypt_b64(data: bytes, aad: bytes, key: bytes) -> str:
        """
        Encrypt a given byte array using AES-GCM and return it as a base64 encoded string
        """
        return base64.b64encode(GCM.encrypt(data, aad, key)).decode("utf-8")

    @staticmethod
    def decrypt(data: bytes, aad: bytes, key: bytes) -> bytes:
        """
        Decrypt a given byte array using AES-GCM
        """
        aesgcm = AESGCM(key)
        nonce = data[-12:]
        data = data[:-12]
        return aesgcm.decrypt(nonce, data, aad)

    @staticmethod
    def decrypt_b64(data: str, aad: bytes, key: bytes) -> bytes:
        """
        Decrypt a given base64 encoded string using AES-GCM
        """
        return GCM.decrypt(base64.b64decode(data), aad, key)

    @staticmethod
    def gen_key(bit_length=256) -> bytes:
        """
        Generate a random AES key
        """
        return AESGCM.generate_key(bit_length=bit_length)

    @staticmethod
    def gen_key_b64(bit_length=256) -> str:
        """
        Generate a random AES key and return it as a base64 encoded string
        """
        return base64.b64encode(GCM.gen_key(bit_length)).decode("utf-8")

    @staticmethod
    def is_valid_key(key: bytes) -> bool:
        try:
            AESGCM(key)
            return True
        except:
            return False


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

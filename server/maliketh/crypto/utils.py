import os


"""
Pad a byte array to a multiple of 16 bytes
"""
def pad(data: bytes) -> bytes:
    return data + (16 - len(data) % 16) * chr(16 - len(data) % 16).encode("utf-8")


"""
Generate a random string of bytes of the specified length
"""
def random_bytes(length: int) -> bytes:
    return os.urandom(length)

"""
Generates a random string of bytes of the specified length and returns it as a hex string
"""
def random_hex(length: int) -> str:
    return os.urandom(length).hex()

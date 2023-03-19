import os
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
from typing import Optional

from maliketh.config import CONFIG_DIR, SERVER_PRIV_KEY_PATH, SERVER_PUB_KEY_PATH
from maliketh.db import db
from maliketh.models import Operator


def generate_ecc_keypair() -> Tuple[PrivateKey, PublicKey]: 
    """
    Generate a new ECC keypair for use with NaCl
    """
    private_key = PrivateKey.generate()
    public_key = private_key.public_key
    return private_key, public_key


def generate_b64_ecc_keypair() -> Tuple[str, str]:
    """
    Generate a new ECC keypair for use with NaCl and return the keys as base64 encoded strings
    """
    private_key, public_key = generate_ecc_keypair()
    return private_key.encode(encoder=Base64Encoder).decode("utf-8"), public_key.encode(
        encoder=Base64Encoder
    ).decode("utf-8")


def load_ecc_keypair(private_key: str, public_key: str) -> Tuple[PrivateKey, PublicKey]:
    """
    Load a keypair from a base64 encoded string
    """
    return PrivateKey(private_key.encode("utf-8"), encoder=Base64Encoder), PublicKey(
        public_key.encode("utf-8"), encoder=Base64Encoder
    )


def encrypt(
    public_key: PublicKey, private_key: PrivateKey, data: bytes, encoder=Base64Encoder
) -> bytes:
    """
    Encrypt a message using NaCl
    """
    box = Box(private_key, public_key)
    return box.encrypt(data, encoder=encoder)


def decrypt(
    public_key: PublicKey, private_key: PrivateKey, data: bytes, encoder=Base64Encoder
) -> bytes:
    """
    Decrypt a message using NaCl
    """
    box = Box(private_key, public_key)
    return box.decrypt(data, encoder=encoder)


def decrypt_with_server_key(priv: str, data: bytes) -> bytes:
    """
    Decrypt `data` using the server's private key

    :return the decrypted data
    """
    with open(SERVER_PUB_KEY_PATH, "r") as f:
        server_pub = f.read()

    if type(data) is str:
        data = data.encode("utf-8")

    private_key, public_key = load_ecc_keypair(priv, server_pub)
    return decrypt(public_key, private_key, data)


def private_key_to_pem(sk: PrivateKey) -> Tuple[str, str]:
    """
    Convert a private key to PEM format
    """
    sk_bytes = bytes(sk)
    sk_pem = (
        Ed25519PrivateKey.from_private_bytes(sk_bytes)
        .private_bytes(
            encoding=Encoding.PEM,
            format=PrivateFormat.PKCS8,
            encryption_algorithm=NoEncryption(),
        )
        .decode("utf-8")
    )
    pk_pem = (
        Ed25519PublicKey.from_public_bytes(bytes(sk.public_key))
        .public_bytes(encoding=Encoding.PEM, format=PublicFormat.SubjectPublicKeyInfo)
        .decode("utf-8")
    )
    return sk_pem, pk_pem


def generate_signing_keypair() -> Tuple[SigningKey, VerifyKey]:
    """
    Generate a new signing keypair for use with NaCl
    """
    signing_key = SigningKey.generate()
    verify_key = signing_key.verify_key
    return signing_key, verify_key


def generate_b64_signing_keypair() -> Tuple[str, str]:
    """
    Generate a new signing keypair for use with NaCl
    """
    signing_key = SigningKey.generate()
    verify_key = signing_key.verify_key
    return signing_key.encode(encoder=Base64Encoder).decode("utf-8"), verify_key.encode(
        encoder=Base64Encoder
    ).decode("utf-8")


def verify_signature(key: VerifyKey, data: bytes) -> Optional[bytes]:
    """
    Verify a signature using NaCl

    :return the data if the signature is valid, None otherwise
    """
    try:
        return key.verify(data)
    except:
        return None


def decrypt_and_verify(
    data: bytes, operator: Operator, encoder=Base64Encoder
) -> Optional[bytes]:
    """
    Decrypt and verify a message using PyNaCl. The message must have been encrypted using the
    server's public key and signed using the operator's signing key.

    :param data: the data to decrypt
    :param operator: the operator's username
    :param encoder: the encoder to use
    :return Optional[bytes]: the decrypted data if the signature is valid, None otherwise
    """
    with open(SERVER_PRIV_KEY_PATH, "rb") as f:
        server_priv_key = PrivateKey(f.read(), encoder=Base64Encoder)

    operator_verify_key = VerifyKey(
        operator.verify_key.encode("utf-8"), encoder=Base64Encoder
    )
    operator_pub_key = PublicKey(
        operator.public_key.encode("utf-8"), encoder=Base64Encoder
    )

    decrypted = decrypt(
        operator_pub_key,
        server_priv_key,
        data,
        encoder=Base64Encoder,
    )

    og_msg = verify_signature(operator_verify_key, decrypted)
    return og_msg

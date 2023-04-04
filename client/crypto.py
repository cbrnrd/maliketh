from nacl.public import PrivateKey, PublicKey, Box
from nacl.signing import SigningKey
from nacl.encoding import Base64Encoder


def encrypt(
    public_key: PublicKey, private_key: PrivateKey, data: bytes, encoder=Base64Encoder
) -> bytes:
    """
    Encrypt a message using NaCl
    """
    box = Box(private_key, public_key)
    return box.encrypt(data, encoder=encoder)


def base64_encrypt_and_sign_str(
    encryption_key_b64: str, recv_pub_key_b64: str, signing_key_b64: str, data: str
) -> str:
    signing_key = SigningKey(signing_key_b64.encode("utf-8"), encoder=Base64Encoder)
    encryption_key = PrivateKey(
        encryption_key_b64.encode("utf-8"), encoder=Base64Encoder
    )
    recv_pub_key = PublicKey(recv_pub_key_b64.encode("utf-8"), encoder=Base64Encoder)

    # Sign data first
    signed_data = signing_key.sign(data.encode("utf-8"))

    # Encrypt data
    encrypted_data = encrypt(
        recv_pub_key, encryption_key, signed_data, encoder=Base64Encoder
    )

    return encrypted_data.decode("utf-8")

from datetime import datetime, timedelta
from typing import Any
from flask import Blueprint, jsonify, request
from maliketh.db import db
from maliketh.models import *
from maliketh.crypto.ec import *
from maliketh.crypto.utils import random_hex, random_string

admin = Blueprint("admin", __name__)


@admin.route("/op/auth/token/request", methods=["GET"])  # type: ignore
def request_token() -> Any:
    # Check if X-ID and X-Signature header is present
    if any(x not in request.headers for x in ["X-ID", "X-Signature"]):
        return jsonify({"status": False}), 400

    # Check if they have content
    if any(len(request.headers[x]) == 0 for x in ["X-ID", "X-Signature"]):
        return jsonify({"status": False}), 400

    # X-Signature is in format base64(enc_and_sign(pub_key, operator_signing_key, server_pub_key))
    # So we need to decrypt it with the server public key,
    # Then verify the signature with the operator signing key

    # Get the operator from X-ID
    operator = Operator.query.filter_by(username=request.headers["X-ID"]).first()
    if operator is None:
        return jsonify({"status": False}), 400

    operator_pub_key = PublicKey(
        operator.public_key.encode("utf-8"), encoder=Base64Encoder
    )
    operator_verify_key = VerifyKey(
        operator.verify_key.encode("utf-8"), encoder=Base64Encoder
    )

    # Decrypt the message
    with open(SERVER_PRIV_KEY_PATH, "rb") as f:
        server_priv_key = PrivateKey(f.read(), encoder=Base64Encoder)
    decrypted = decrypt(
        operator_pub_key,
        server_priv_key,
        bytes(request.headers["X-Signature"], "utf-8"),
        encoder=Base64Encoder,
    )
    if decrypted is None:
        return jsonify({"status": False, "msg": "Empty decrypted message"}), 400

    # Verify the signature
    if not verify_signature(operator_verify_key, decrypted):
        return jsonify({"status": False, "msg": "Couldn't verify signature"}), 400

    # Generate a new token
    token = random_hex(128)
    operator.auth_token = token
    operator.auth_token_expiry = datetime.now() + timedelta(hours=6)
    db.session.commit()

    # TODO rabbitmq stuff

    return (
        jsonify(
            {
                "status": True,
                "token": token,
                "rmq_queue": "",
                "rmq_host": "",
                "rmq_port": 1,
            }
        ),
        200,
    )

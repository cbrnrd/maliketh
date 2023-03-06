from datetime import datetime, timedelta
from typing import Any
from flask import Blueprint, jsonify, request
from maliketh.db import db
from maliketh.models import *
from maliketh.crypto.ec import *
from maliketh.crypto.utils import random_hex, random_string

admin = Blueprint("admin", __name__)


def verify_auth_token(request) -> Optional[str]:
    """
    Given a request, verify its Authentication header

    :param request: The request to verify
    :return The operator username if the token is valid, otherwise None
    """
    # Get bearer token
    token = request.headers.get("Authorization", None)
    if token is None:
        return None

    # Check if it's a bearer token
    if not token.startswith("Bearer "):
        return None

    # Get the token
    token = token[7:]

    # Get the operator
    operator = Operator.query.filter_by(auth_token=token).first()
    if operator is None:
        return None

    # Check if the token is still valid
    token_exp = datetime.strptime(operator.auth_token_expiry, "%Y-%m-%d %H:%M:%S")
    if token_exp < datetime.now():
        return None

    return operator


def verified(func):
    """
    Decorator to check if the request is authenticated
    :param func: The function to decorate
    :return The decorated function with the operator as a keyword argument
    """
    def wrapper(*args, **kwargs):
        # Check if the request is authenticated
        operator = verify_auth_token(request)
        if operator is None:
            return jsonify({"status": False, "msg": "Not authenticated"}), 401
        kwargs = {**kwargs, "operator": operator}
        return func(*args, **kwargs)

    return wrapper


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

    original_message = decrypt_and_verify(
        bytes(request.headers["X-Signature"], "utf-8"), operator
    )

    if original_message is None:
        return jsonify({"status": False, "msg": "Couldn't verify signature"}), 400

    original_message = original_message.decode("utf-8")

    if original_message != operator.login_secret:
        return (
            jsonify(
                {
                    "status": False,
                    "msg": "Unable to verify signature",
                }
            ),
            400,
        )

    # If we get here, the operator is authenticated

    # Check if the token is still valid

    if (
        operator.auth_token_expiry is not None
        and datetime.strptime(operator.auth_token_expiry, "%Y-%m-%d %H:%M:%S")
        > datetime.now()
    ):
        token = operator.auth_token
    else:
        # Generate a new token
        token = random_hex(128)
        operator.auth_token = token
        operator.auth_token_expiry = (datetime.now() + timedelta(hours=6)).strftime(
            "%Y-%m-%d %H:%M:%S"
        )
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


@admin.route("/op/auth/token/revoke", methods=["GET"])  # type: ignore
@verified
def revoke_token(operator: Operator) -> Any:
    operator.auth_token = None
    operator.auth_token_expiry = None
    db.session.commit()
    return jsonify({"status": True}), 200
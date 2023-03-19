from typing import Union
import requests
from dataclasses import dataclass


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

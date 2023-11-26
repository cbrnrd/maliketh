from typing import Callable, Tuple
import uuid
from flask import jsonify, Response, Blueprint
from functools import wraps
import flask

import structlog
from maliketh.config import OP_ROUTES


def basic_status_json(status: bool, msg: str, code: int = 401) -> Tuple[Response, int]:
    return jsonify({"status": status, "msg": msg}), code


def error_json(msg: str, code: int = 401) -> Tuple[Response, int]:
    return basic_status_json(False, msg, code)


def success_json(msg: str, code: int = 200) -> Tuple[Response, int]:
    return basic_status_json(True, msg, code)


def create_route(bp: Blueprint, route_name: str) -> Callable:
    """
    Utility function to create a route in a blueprint for a given endpoint.

    :param bp: The flask blueprint to add the route to
    :param route_name: The name of the route in the admin config file
    """

    def inner(func):
        @bp.route(
            OP_ROUTES[route_name]["path"], methods=OP_ROUTES[route_name]["methods"]
        )
        @wraps(func)
        def wrap(*args, **kwargs):
            return func(*args, **kwargs)

        return wrap

    return inner

def setup_logger(func):
    @wraps(func)
    def inner(*args, **kwargs):
        structlog.contextvars.clear_contextvars()
        structlog.contextvars.bind_contextvars(
            view=flask.request.path,
            request_id=str(uuid.uuid4()),
            peer=flask.request.access_route[0],
        )
        return func(*args, **kwargs)
    
    return inner

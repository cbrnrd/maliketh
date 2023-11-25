from typing import Callable, Tuple
from flask import jsonify, Response, Blueprint
from functools import wraps
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

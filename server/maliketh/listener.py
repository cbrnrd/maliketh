from crypt import methods
from flask import Blueprint, jsonify, request
from maliketh.db import db
from maliketh.models import Task, CREATED, TASKED, COMPLETE

c2 = Blueprint('c2', __name__)



@c2.route("/index.html")
def hello_c2():
    return "Nothing to see here"
from datetime import datetime
from flask import Blueprint, jsonify, request
from maliketh.db import db
import maliketh.crypto.aes
from maliketh.models import *

admin = Blueprint('admin', __name__)

"""
How operators work:
Operators are users of the C2 that can create tasks for implants.
When an operator is created, it is given an ECC certificate signed by the C2's CA certificate.
"""

"""
Body should look like:
{
    "username": "admin",
    "cert": "---BEGIN CERTIFICATE---\n...\n---END CERTIFICATE---"
}
"""
@admin.route("/auth", methods=["POST"])
def auth():
    # TODO
    pass
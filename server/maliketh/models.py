from dataclasses import dataclass, asdict
from maliketh.db import db
import os
import sys
import json

CREATED = "CREATED"
TASKED = "TASKED"
COMPLETE = "COMPLETE"
ERROR = "ERROR"


def random_id():
    return os.urandom(16).hex()


@dataclass
class TestDB(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    foo: str = db.Column(db.String)
    bar: str = db.Column(db.String)

    def toJSON(self):
        return asdict(self)


"""
An implant is a remote computer that is running the maliketh agent
"""


@dataclass
class Implant(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    implant_id: str = db.Column(db.String)
    hostname: str = db.Column(db.String)
    ip: str = db.Column(db.String)
    os: str = db.Column(db.String)
    arch: str = db.Column(db.String)
    user: str = db.Column(db.String)
    aes_key: str = db.Column(db.String)  # Base64 encoded AES key
    created_at: str = db.Column(db.String)
    last_seen: str = db.Column(db.String)

    def toJSON(self):
        return asdict(self)

    def __repr__(self):
        return f"Implant: {self.implant_id} {self.hostname} {self.ip} {self.os} {self.arch} {self.user} {self.created_at} {self.last_seen}"


def get_implant_by_id(implant_id: str):
    return Implant.query.filter_by(implant_id=implant_id).first()


"""
A task is a job given by an operator to an implant.
A job has "owned by" an operator and "executed by" an implant
"""
@dataclass
class Task(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    operator_id: str = db.Column(db.String)
    task_id: str = db.Column(db.String)
    implant_id: str = db.Column(db.String)
    cmd: str = db.Column(db.String)
    args: str = db.Column(db.String)
    status: str = db.Column(db.String)
    output: str = db.Column(db.String)
    created_at: str = db.Column(db.String)
    read_at: str = db.Column(db.String)
    executed_at: str = db.Column(db.String)

    def toJSON(self):
        return asdict(self)


def get_task_by_id(task_id: str):
    return Task.query.filter_by(task_id=task_id).first()


def get_oldest_task_for_implant(implant_id: str):
    return Task.query.filter_by(implant_id=implant_id, status=CREATED).first()


"""
An operator is a user of this server that is able to issue commands to implants
"""
@dataclass
class Operator(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username: str = db.Column(db.String)
    auth_token: str = db.Column(db.String)
    # public_key : str = db.Column(db.String)  # The public key of the operator to authenticate to server
    created_at: str = db.Column(db.String)
    last_login: str = db.Column(db.String)

    def toJSON(self):
        return asdict(self)

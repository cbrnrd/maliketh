from dataclasses import dataclass, asdict
from datetime import datetime
from typing import Any
from maliketh.db import db
from maliketh.profile import *
import base64
import os
import sys
import json
import yaml


# Job statuses
CREATED = "CREATED"
TASKED = "TASKED"
COMPLETE = "COMPLETE"
ERROR = "ERROR"


def random_id(n=16):
    return os.urandom(n).hex()[0:n]


@dataclass
class TestDB(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    foo: str = db.Column(db.String)
    bar: str = db.Column(db.String)

    def toJSON(self):
        return asdict(self)


@dataclass
class Implant(db.Model):
    """
    An implant is a remote computer that is running the maliketh agent"""

    id = db.Column(db.Integer, primary_key=True)
    implant_id: str = db.Column(db.String)
    hostname: str = db.Column(db.String)
    ip: str = db.Column(db.String)
    os: str = db.Column(db.String)
    arch: str = db.Column(db.String)
    user: str = db.Column(db.String)
    server_sk: str = db.Column(
        db.String
    )  # The Base64 encoded private key (of the server, EC) to use for encryption
    implant_pk: str = db.Column(
        db.String
    )  # The Base64 encoded public key (of the implant, EC) to use for encryption
    created_at: str = db.Column(db.String)
    last_seen: str = db.Column(db.String)

    def toJSON(self):
        return asdict(self)

    def __repr__(self):
        return f"Implant: {self.implant_id} {self.hostname} {self.ip} {self.os} {self.arch} {self.user} {self.created_at} {self.last_seen}"


def get_implant_by_id(implant_id: str):
    return Implant.query.filter_by(implant_id=implant_id).first()


@dataclass
class ImplantConfig(db.Model):
    """
    An implant config is a changeable configuration for an implant to use.
    """

    id = db.Column(db.Integer, primary_key=True)
    implant_id: str = db.Column(
        db.String
    )  # The implant ID this config belongs to (note this is not the `id` field of Implant, rather the `implant_id` field))
    cookie: str = db.Column(
        db.String
    )  # The cookie name to use for implant identification
    kill_date: str = db.Column(db.String)  # The timestamp of the kill date
    sleep_time: int = db.Column(
        db.Integer
    )  # The number of seconds to sleep between tasks/checkin
    jitter: float = db.Column(
        db.Float
    )  # The percentage of jitter to add to the sleep time
    max_retries: int = db.Column(
        db.Integer
    )  # The number of times to retry a task before giving up. -1 for infinite
    enc_key: str = db.Column(
        db.String
    )  # The Base64 encoded public key (of the server, EC) to use for encryption
    tailoring_hash_function: str = db.Column(
        db.String
    )  # The hash function to use for payload tailoring
    tailoring_hash_rounds: str = db.Column(
        db.Integer
    )  # The number of rounds to use for the given hash function
    tailoring_hashes: List[str] = db.Column(
        db.ARRAY(db.String)
    )  # The hashes to use for payload tailoring

    def toJSON(self):
        dicted = asdict(self)

        del dicted["implant_id"]

        return dicted

    def toYAML(self):
        return yaml.dump(self.toJSON())

    @staticmethod
    def from_profile(profile: MaleableProfile, implant_id: str, enc_key: str) -> "ImplantConfig":
        implant_profile = profile.implant_profile
        return ImplantConfig(
            implant_id=implant_id,
            cookie=profile.globals.implant_id_cookie,
            kill_date=implant_profile.kill_date,
            sleep_time=implant_profile.sleep,
            jitter=implant_profile.jitter,
            max_retries=implant_profile.max_retries,
            enc_key=enc_key,
            tailoring_hash_function=implant_profile.tailoring_hash_function,
            tailoring_hash_rounds=implant_profile.tailoring_hash_rounds,
            tailoring_hashes=implant_profile.tailoring_hashes,
        )



    @staticmethod
    def create_min_config(
        implant_id: str, cookie: str, b64_server_pub: str
    ) -> "ImplantConfig":
        """
        Helper to create a new config with the minimum required fields, and add it to the database.
        """

        config = ImplantConfig(
            implant_id=implant_id,
            kill_date="",
            sleep_time=60,
            jitter=0.1,
            max_retries=-1,
            enc_key=b64_server_pub,
            tailoring_hash_function="",
            tailoring_hash_rounds=100,
            tailoring_hashes=[],
            cookie=cookie,
        )
        db.session.add(config)
        db.session.commit()
        return config

    def add_hash(self, hash: str) -> None:
        """
        Add a hash to the tailoring hashes
        """
        self.tailoring_hashes.append(hash)
        db.session.commit()

    def remove_hash(self, hash: str) -> None:
        """
        Remove a hash from the tailoring hashes. If `hash` is 
        not in the list, do nothing.
        """
        try:
            self.tailoring_hashes.remove(hash)
        except ValueError:
            pass
        db.session.commit()

@dataclass
class Task(db.Model):
    """
    A task is a job given by an operator to an implant.
    A job has "owned by" an operator and "executed by" an implant"""

    id = db.Column(db.Integer, primary_key=True)
    operator_name: str = db.Column(
        db.String
    )  # The username of the operator that created this task
    task_id: str = db.Column(db.String)  # The task ID
    implant_id: str = db.Column(
        db.String
    )  # The ID of the implant that this task is assigned to
    opcode: int = db.Column(db.Integer)  # The opcode of the task
    args: str = db.Column(db.JSON)  # The arguments of the task
    status: str = db.Column(db.String)  # The status of the task
    output: str = db.Column(db.String)  # The output of the task
    created_at: str = db.Column(db.String)  # The datetime the task was created
    read_at: str = db.Column(db.String)  # The datetime the task was read by the implant
    executed_at: str = db.Column(
        db.String
    )  # The datetime the task was executed by the implant

    def toJSON(self):
        return asdict(self)

    def to_filtered_json(self) -> str:
        """
        Returns a JSON string with only the fields needed by the implant.
        """
        return json.dumps(
            {
                "task_id": self.task_id,
                "opcode": self.opcode,
                "args": self.args,
            }
        )

    @staticmethod
    def new_task(operator_name: str, implant_id: str, opcode: int, args: Any):
        """
        Helper to create a new task with the minimum required fields, and add it to the database.
        """
        task_id = random_id()
        task = Task(
            operator_name=operator_name,
            task_id=task_id,
            implant_id=implant_id,
            opcode=opcode,
            args=args,
            status=CREATED,
            created_at=datetime.now(),
        )
        db.session.add(task)
        db.session.commit()
        return task


def get_task_by_id(task_id: str):
    return Task.query.filter_by(task_id=task_id).first()


def get_oldest_task_for_implant(implant_id: str):
    return Task.query.filter_by(implant_id=implant_id, status=CREATED).first()


@dataclass
class Operator(db.Model):
    """
    An operator is a user of this server that is able to issue commands to implants"""

    id = db.Column(db.Integer, primary_key=True)
    username: str = db.Column(db.String)
    public_key: str = db.Column(db.String)
    verify_key: str = db.Column(db.String)
    login_secret: str = db.Column(db.String)
    auth_token: str = db.Column(db.String)
    auth_token_expiry: str = db.Column(db.String)  # The time the auth token expires
    created_at: str = db.Column(db.String)
    last_login: str = db.Column(db.String)
    rmq_queue: str = db.Column(db.String)  # The RabbitMQ queue name for this operator

    def toJSON(self):
        return asdict(self)

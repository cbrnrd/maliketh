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
    id  = db.Column(db.Integer, primary_key=True)
    foo: str = db.Column(db.String)
    bar: str = db.Column(db.String)
    def toJSON(self):
        return asdict(self)

@dataclass
class Task(db.Model):
    id =  db.Column(db.Integer, primary_key=True)
    task_id : str = db.Column(db.String)
    implant_id : str = db.Column(db.String)
    cmd : str = db.Column(db.String)
    args : str = db.Column(db.String)
    
    status : str =  db.Column(db.String)
    
    # @TODO it would be nice to know when the job was created! when it was pulled down and when it was executed
    ##timestamp =  
    def toJSON(self):
        return asdict(self)
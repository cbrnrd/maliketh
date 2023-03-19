from datetime import datetime
import json
import os
import sys
from flask import Flask
from maliketh.config import CONFIG_DIR
from maliketh.operator.config import generate_config 
from maliketh.db import db
from maliketh.models import Operator
from maliketh.logging.standard_logger import StandardLogger, LogLevel
from maliketh.crypto.ec import generate_b64_ecc_keypair

def build_operator_app():
    app = Flask("operator")
    app.config.from_mapping(
        # default secret that should be overridden in environ or config
        SQLALCHEMY_DATABASE_URI="postgresql://postgres:Vv4QZnP7eS#K7Z!4HWzx@postgres:5432",
    ) 
    from .listeners.admin import admin as admin_blueprint
    app.register_blueprint(admin_blueprint)
    
    db.init_app(app)   

    return app


def build_c2_app():
    app = Flask("c2")
    app.config.from_mapping(
        # default secret that should be overridden in environ or config
        SQLALCHEMY_DATABASE_URI="postgresql://postgres:Vv4QZnP7eS#K7Z!4HWzx@postgres:5432",
    )

    from .listeners.c2 import c2 as c2_blueprint
    app.register_blueprint(c2_blueprint)

    db.init_app(app)

    return app

def init_db():
    logger = StandardLogger(sys.stdout, sys.stderr, LogLevel.INFO)
    logger.info("Initializing database")

    db.drop_all()
    db.create_all()

    logger.info("Generating server keypair")
    server_sk, server_pk = generate_b64_ecc_keypair()
    with open(os.path.join(CONFIG_DIR, 'admin', 'certs', 'server_priv'), 'w') as f:
        f.write(server_sk)
    with open(os.path.join(CONFIG_DIR, 'admin', 'certs', 'server_pub'), 'w') as f:
        f.write(server_pk)
    
    logger.info("Done, server keypair written to config/admin/certs/")
    
    logger.info("Generating admin credentials:")
    # Add admin operator
    admin_config = generate_config("admin", None)
    admin = Operator(
        username=admin_config['name'],
        public_key=admin_config['public'],
        verify_key=admin_config['verify_key'],
        login_secret=admin_config['login_secret'],
        auth_token=None,
        auth_token_expiry=None,
        created_at=datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
    )

    db.session.add(admin)
    db.session.commit()
    print(json.dumps(admin_config, indent=4))

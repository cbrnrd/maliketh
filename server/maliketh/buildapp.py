import datetime
from flask import Flask
from maliketh.config import get_admin_creds 
from maliketh.db import db
from maliketh.models import Operator

def build_app():
    app = Flask(__name__)
    app.config.from_mapping(
        # default secret that should be overridden in environ or config
        SQLALCHEMY_DATABASE_URI="sqlite:///c2.db",
    ) 
    # from .admin import admin as admin_blueprint
    # app.register_blueprint(admin_blueprint)

    from .listeners.c2 import c2 as c2_blueprint
    app.register_blueprint(c2_blueprint)
    db.init_app(app)   

    return app

def init_db():
    db.drop_all()
    db.create_all()

    # Add admin operator
    admin_creds = get_admin_creds()
    admin = Operator(
        username=admin_creds[0],
        auth_token=admin_creds[1],
        created_at=datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
    )

    db.session.add(admin)
    db.session.commit()

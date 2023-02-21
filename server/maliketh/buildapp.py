from flask import Flask 
from maliketh.db import db

def build_app():
    app = Flask(__name__)
    app.config.from_mapping(
        # default secret that should be overridden in environ or config
        SQLALCHEMY_DATABASE_URI="sqlite:///c2.db",
    ) 
    # from .admin import admin as admin_blueprint
    # app.register_blueprint(admin_blueprint)

    from .listener import c2 as c2_blueprint
    app.register_blueprint(c2_blueprint)
    db.init_app(app)   

    return app

def init_db():
    db.drop_all()
    db.create_all()

from app import operator_app
from maliketh.buildapp import init_db

with operator_app.app_context():
    init_db()

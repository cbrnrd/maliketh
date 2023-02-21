from maliketh.app import app
from maliketh.buildapp import init_db

with app.app_context():
    init_db()

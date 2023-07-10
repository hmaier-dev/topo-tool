from flask import Flask, app


class Web:
    def __init__(self):
        self.app = Flask(__name__)

    @app.route('/')
    def base(self):
        return 'hello!'

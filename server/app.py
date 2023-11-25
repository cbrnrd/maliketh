import os
import sys
import gunicorn.app.base
from maliketh.buildapp import build_operator_app, build_c2_app, init_db
from maliketh.logging.standard_logger import StandardLogger, LogLevel
from maliketh.operator.rmq import rmq_setup
from optparse import OptionParser

from maliketh.config import set_c2_profile, DEFAULT_C2_PROFILE

# init_db()

operator_app = build_operator_app()
c2_app = build_c2_app()


class StandaloneApplication(gunicorn.app.base.BaseApplication):
    def __init__(self, app, options=None):
        self.options = options or {}
        self.application = app
        super().__init__()

    def load_config(self):
        config = {
            key: value
            for key, value in self.options.items()
            if key in self.cfg.settings and value is not None
        }
        for key, value in config.items():
            self.cfg.set(key.lower(), value)

    def load(self):
        return self.application


def parse_options():
    opts = {}
    parser = OptionParser()

    parser.description = "Maliketh Server"
    parser.usage = "usage: %prog <--start-operator> [options]"

    operator_group = parser.add_option_group("Operator Listener Options")
    operator_group.add_option(
        "",
        "--start-operator",
        dest="start_operator",
        default=False,
        help="Start operator listener",
        action="store_true",
    )
    operator_group.add_option(
        "-p",
        "--port",
        dest="port",
        default=5000,
        help="Port to listen on for operator connections",
    )
    operator_group.add_option(
        "-a",
        "--address",
        dest="address",
        default="0.0.0.0",
        help="Address to listen on for operator connections",
    )
    operator_group.add_option(
        "-l", "--log-level", dest="log_level", default="INFO", help="Log level"
    )

    c2_group = parser.add_option_group("C2 Listener Options")
    c2_group.add_option(
        "", "--start-c2", dest="start_c2", default=True, help="Start C2 listener"
    )
    c2_group.add_option(
        "",
        "--profile",
        dest="profile",
        default=DEFAULT_C2_PROFILE,
        help="C2 maleable profile to use",
    )
    c2_group.add_option(
        "-c",
        "--c2-port",
        dest="c2_port",
        default=8080,
        help="Port to listen on for C2 connections",
    )
    c2_group.add_option(
        "-b",
        "--c2-address",
        dest="c2_address",
        default="0.0.0.0",
        help="Address to listen on for C2 connections",
    )
    c2_group.add_option(
        "-f", "--c2-log-level", dest="c2_log_level", default="INFO", help="Log level"
    )

    misc_group = parser.add_option_group("Miscellaneous Options")
    misc_group.add_option(
        "-i", "--init-db", dest="init_db", default=False, help="Initialize database"
    )
    misc_group.add_option(
        "-d",
        "--debug",
        dest="debug",
        default=False,
        help="Enable debug mode",
        action="store_true",
    )
    opts, args = parser.parse_args()
    return opts


def validate_args(opts) -> None:
    if opts.log_level not in LogLevel.get_names():
        print("Invalid log level")
        sys.exit(1)
    if opts.c2_log_level not in LogLevel.get_names():
        print("Invalid log level")
        sys.exit(1)
    if opts.profile is not None and not os.path.exists(opts.profile):
        print(f"Profile {opts.profile} does not exist")
        sys.exit(1)


def init_simple_logger(level: LogLevel = LogLevel.INFO) -> StandardLogger:
    return StandardLogger(sys.stdout, sys.stderr, level)  # type: ignore


def main():
    global operator_app, c2_app
    opts = parse_options()
    if opts.init_db:
        choice = input(
            "Are you sure you want to initialize the database? You will lose all operators and implants (y/n): "
        )
        if choice != "y" or choice != "Y":
            sys.exit(0)
        init_db()
        sys.exit(0)

    logger = init_simple_logger(LogLevel[opts.log_level])

    validate_args(opts)
    set_c2_profile(opts.profile)

    if opts.start_operator:
        logger.info("Starting operator listener on %s:%s" % (opts.address, opts.port))
        operator_app.run(host=opts.address, port=opts.port, debug=opts.debug)
        # StandaloneApplication(operator_app, {}).run(host=opts.address, port=opts.port, debug=opts.debug)
    elif opts.start_c2:
        logger.info("Starting C2 listener on %s:%s" % (opts.c2_address, opts.c2_port))
        c2_app.run(host=opts.c2_address, port=opts.c2_port, debug=opts.debug)
        # StandaloneApplication(c2_app, {}).run(host=opts.c2_address, port=opts.c2_port, debug=opts.debug)


if __name__ == "__main__":
    main()

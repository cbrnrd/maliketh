from datetime import datetime
import json
import sys
from maliketh.db import db
from app import operator_app as app

from optparse import OptionParser

from maliketh.models import Operator
from maliketh.operator.config import generate_config
from maliketh.logging.standard_logger import LogLevel, StandardLogger


def main():
    logger = StandardLogger(sys.stdout, sys.stderr, LogLevel.INFO)
    parser = OptionParser()
    parser.usage = "%prog -n <operator name> [options]"
    required = parser.add_option_group("Required arguments")
    required.add_option(
        "-n",
        "--name",
        action="store",
        dest="name",
        default=None,
        help="Name of the operator",
    )

    actions = parser.add_option_group("Actions")
    actions.add_option(
        "-d",
        "--delete",
        action="store_true",
        dest="delete",
        default=False,
        help="Delete an operator from the database",
    )
    actions.add_option(
        "-c",
        "--create",
        action="store_true",
        dest="create",
        default=False,
        help="Create an operator in the database",
    )
    actions.add_option(
        "-l",
        "--list",
        action="store_true",
        dest="list",
        default=False,
        help="List all operators in the database",
    )

    other = parser.add_option_group("Other")
    other.add_option(
        "-v",
        "--verbose",
        action="store_true",
        dest="verbose",
        default=False,
        help="Verbose output",
    )
    other.add_option(
        "-q",
        "--quiet",
        action="store_true",
        dest="quiet",
        default=False,
        help="Quiet output",
    )
    other.add_option(
        "-o",
        "--output",
        action="store",
        dest="output",
        default=None,
        help="Output file (default: stdout)",
    )

    if len(sys.argv) == 1:
        parser.print_help()
        sys.exit(1)

    (options, args) = parser.parse_args()

    if not (options.list or options.delete):
        if not options.create:
            options.create = True

    with app.app_context():
        if options.list:
            operators = Operator.query.all()
            for operator in operators:
                print(operator.username)
            sys.exit(0)

        if options.create:
            if options.name is None:
                logger.error("You must specify a name for the operator")
                sys.exit(1)
            else:
                operator = Operator.query.filter_by(username=options.name).first()
                if operator is not None:
                    logger.error("Operator already exists")
                    sys.exit(1)
                else:
                    operator_config = generate_config(options.name, None)
                    operator = Operator(
                        username=operator_config["name"],
                        public_key=operator_config["public"],
                        login_secret=operator_config["login_secret"],
                        verify_key=operator_config["verify_key"],
                        auth_token=None,
                        auth_token_expiry=None,
                        created_at=datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
                        rmq_queue=operator_config["rmq_queue"],
                    )
                    db.session.add(operator)
                    db.session.commit()
                    if options.output is not None:
                        with open(options.output, "w") as f:
                            json.dump(operator_config, f, indent=4)
                    else:
                        print(json.dumps(operator_config, indent=4))

        if options.delete:
            if options.name is None:
                logger.error("You must specify a name for the operator")
                sys.exit(1)
            else:
                operator = Operator.query.filter_by(username=options.name).first()
                if operator is None:
                    logger.error("Operator does not exist")
                    sys.exit(1)
                else:
                    db.session.delete(operator)
                    db.session.commit()
                    logger.ok("Operator deleted")


if __name__ == "__main__":
    main()

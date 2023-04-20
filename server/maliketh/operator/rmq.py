from datetime import datetime
import time
import pika
import pika.exceptions
from maliketh.models import Operator


def rmq_setup(max_retry=5, retry_delay=5):
    """
    Creates the RabbitMQ exchange and queues for the operators.
    """
    connection = None
    for i in range(max_retry):
        try:
            connection = pika.BlockingConnection(pika.ConnectionParameters(host="rabbitmq"))
            break
        except pika.exceptions.AMQPConnectionError:
            print(f"Failed to connect to RabbitMQ. Retrying in {retry_delay} seconds...")
            time.sleep(retry_delay)
    if connection is None:
        raise pika.exceptions.AMQPConnectionError

    channel = connection.channel()

    # Define an exchange for global logs (to be sent to all operators)
    channel.exchange_declare(exchange="announcements", exchange_type="fanout")

    # Define an exchange for operator logs (to be sent to a specific operator)
    channel.exchange_declare(exchange="logs", exchange_type="direct")

    # Create a queue for each operator
    # for op in Operator.query.all():
    #     channel.queue_declare(queue=op.rmq_queue)
    #     channel.queue_bind(exchange="logs", queue=op.rmq_queue, routing_key=op.rmq_queue)

    connection.close()


def send_message_to_operator(op: Operator, msg: str):
    """
    Sends a message to the operator's RabbitMQ queue.
    """
    connection = pika.BlockingConnection(pika.ConnectionParameters(host="rabbitmq"))
    channel = connection.channel()
    channel.basic_publish(exchange="logs", routing_key=op.rmq_queue, body=msg)
    connection.close()


def send_message_to_all_queues(msg: str):
    """
    Sends a message to all RabbitMQ queues using fanout.
    """
    connection = pika.BlockingConnection(pika.ConnectionParameters(host="rabbitmq"))
    channel = connection.channel()
    channel.basic_publish(exchange="announcements", routing_key="", body=msg)
    connection.close()


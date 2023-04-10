from datetime import datetime
import pika
from maliketh.models import Operator

def rmq_setup():
    """
    Creates the RabbitMQ exchange and queues for the operators.
    """
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host="rabbitmq")
    )
    channel = connection.channel()
    channel.exchange_declare(exchange="logs", exchange_type="fanout")
    connection.close()

def send_message_to_operator(op: Operator, msg: str):
    """
    Sends a message to the operator's RabbitMQ queue.
    """
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host="rabbitmq")
    )
    channel = connection.channel()
    channel.queue_declare(queue=op.rmq_queue)
    channel.basic_publish(exchange="logs", routing_key=op.rmq_queue, body=msg)
    connection.close()

def send_message_to_all_queues(msg: str):
    """
    Sends a message to all RabbitMQ queues using fanout.
    """
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host="rabbitmq")
    )
    msg = f"[{datetime.now()}] {msg}"
    channel = connection.channel()
    channel.exchange_declare(exchange="logs", exchange_type="fanout")
    channel.basic_publish(exchange="logs", routing_key="", body=msg)
    connection.close()

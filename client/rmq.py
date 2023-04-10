from config import OperatorConfig
from threading import Thread
from cli.logging import get_styled_logger
import pika

def listen_for_messages_in_thread(op: OperatorConfig):
    logger = get_styled_logger()
    def callback(ch, method, properties, body):
        logger.ok(f"{body.decode()}")

    def listen_for_messages():
        connection = pika.BlockingConnection(
            pika.ConnectionParameters(host=op.c2)
        )
        channel = connection.channel()

        channel.queue_bind(exchange='logs', queue=op.rmq_queue)

        channel.basic_consume(
            queue=op.rmq_queue, on_message_callback=callback, auto_ack=True
        )
        channel.start_consuming()

    t = Thread(target=listen_for_messages, daemon=True)
    t.start()

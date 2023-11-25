from typing import Any, Dict
from config import OperatorConfig
from threading import Thread
from cli.logging import get_styled_logger
import pika
import datetime


def listen_for_messages_in_thread(op: OperatorConfig, cli_opts: Dict[str, Any]):
    logger = get_styled_logger()

    def callback(ch, method, properties, body):
        msg = f"[{method.exchange}] {f'[{datetime.datetime.now()}] ' if cli_opts.with_timestamps else ''}{body.decode()}"
        logger.ok(msg)

    def listen_for_messages():
        connection = pika.BlockingConnection(pika.ConnectionParameters(host=op.c2))
        channel = connection.channel()

        channel.exchange_declare(exchange="announcements", exchange_type="fanout")
        channel.exchange_declare(exchange="logs", exchange_type="direct")

        announcement_queue = channel.queue_declare(queue="", exclusive=True)
        logs_queue = channel.queue_declare(queue="", exclusive=False)
        announcement_queue_name = announcement_queue.method.queue
        logs_queue_name = logs_queue.method.queue

        channel.queue_bind(exchange="announcements", queue=announcement_queue_name)
        channel.queue_bind(
            exchange="logs", queue=logs_queue_name, routing_key=op.rmq_queue
        )

        channel.basic_consume(
            queue=announcement_queue_name, on_message_callback=callback, auto_ack=True
        )
        channel.basic_consume(
            queue=logs_queue_name, on_message_callback=callback, auto_ack=True
        )

        channel.start_consuming()

    t = Thread(target=listen_for_messages, daemon=True)
    t.start()

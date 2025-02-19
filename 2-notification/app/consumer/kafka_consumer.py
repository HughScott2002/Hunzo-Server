from confluent_kafka import Consumer
from app.services.email import send_email
from app.core.config import settings

def start_kafka_consumer():
    conf = {
        'bootstrap.servers': settings.kafka_bootstrap_servers,
        'group.id': 'email-service',
        'auto.offset.reset': 'earliest'
    }
    consumer = Consumer(conf)
    consumer.subscribe(['notifications'])

    while True:
        msg = consumer.poll(1.0)
        if msg is None:
            continue
        if msg.error():
            print(f"Consumer error: {msg.error()}")
            continue
        # Process message
        send_email(...)  # Your email logic here
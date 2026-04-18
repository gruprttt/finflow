from kafka import KafkaConsumer
import time

consumer = None

while True:
    try:
        print("Trying to connect to Kafka...")
        consumer = KafkaConsumer(
            'order_created',
            bootstrap_servers='kafka:9092',
            auto_offset_reset='earliest',
            group_id=None,
        )
        print("Connected to Kafka ✅")
        break
    except Exception as e:
        print("Kafka not ready, retrying...", e)
        time.sleep(5)

print("Consumer started...")

for message in consumer:
    print("Received event:", message.value.decode())

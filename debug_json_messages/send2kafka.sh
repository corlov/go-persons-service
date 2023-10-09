#!/bin/bash

echo "send correct message"
cat msg.json | /home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic FIO > /dev/null
sleep 3

echo "send correct message #2"
cat msg2.json | /home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic FIO > /dev/null
sleep 3


echo "send incorrect message"
cat err_msg.json | /home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic FIO > /dev/null
sleep 3

echo "send incorrect message #2"
cat err_msg2.json | /home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic FIO > /dev/null
sleep 3
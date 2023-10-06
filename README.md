CREATE TABLE "Population".Person (
	id bigserial NOT NULL, 
	name varchar(100) NOT NULL, 
	surname varchar(100) NOT NULL,
	patronymic varchar(100),
	age int,
	-- todo: сделать внешним ключом
	country_id  varchar(100),
	-- todo: сделать внешним ключом
	gender_id varchar(100),
		
	creaed_at timestamptz  default now() NOT NULL,
	
	CONSTRAINT pk_keyword_id PRIMARY KEY (id)
);
COMMENT ON TABLE Person.keyword IS 'граждане';
---------------------------------------------------------------------------------



cat ./debug_json_messages/msg.json | /home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic FIO > /dev/nul


/home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic FIO_FAILED --from-beginning

-------------------------------------------------------------
Сатьи, которые использоваляиь для написания проекта:

https://www.digitalocean.com/community/tutorials/how-to-install-apache-kafka-on-ubuntu-20-04
https://www.sohamkamani.com/golang/working-with-kafka/
------------------

sudo nano /etc/systemd/system/zookeeper.service

sudo nano /etc/systemd/system/kafka.service




sudo systemctl enable zookeeper

sudo systemctl start zookeeper

sudo systemctl status zookeeper


sudo systemctl start kafka

sudo systemctl status kafka

sudo systemctl enable kafka



sudo systemctl daemon-reload


Это на версиях Кафки старше чем 2.2 не будет работать (устаревшая опция):
/home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic TutorialTopic

поэтому создаем топик так:
/home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-topics.sh --create --topic TutorialTopic --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1                     
(https://stackoverflow.com/questions/69297020/exception-in-thread-main-joptsimple-unrecognizedoptionexception-zookeeper-is)




echo "Hello, World" | /home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic TutorialTopic > /dev/null


/home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic TutorialTopic --from-beginning

=========================
1) установить Kafka cluster на локальной машине

2) создать очередь FIO

/home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-topics.sh --create --topic FIO --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1                     

3) Проверка входящих сообщений, публикация в очередь с ошибками тех что не прошли верификацию

/home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-topics.sh --create --topic FIO_FAILED --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1  
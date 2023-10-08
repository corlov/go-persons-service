
**GraphQL service API**

  

  

https://betterprogramming.pub/building-a-graphql-server-using-the-schema-first-approach-in-golang-a8da71d7e5b7

  

  

Exec commands bellow to generate your own GraphQL handlers by GraphQL schema:

  
  
  

export GOPATH=/home/kostya/go/src

go mod init

go get github.com/99designs/gqlgen

  

go run github.com/99designs/gqlgen init

  

go mod tidy

go run github.com/99designs/gqlgen generate

  

  

Sometimes some errors happened, you need to try to install manually:

  

go get github.com/99designs/gqlgen/codegen/config@v0.17.39

go get github.com/99designs/gqlgen/internal/imports@v0.17.39

go get github.com/99designs/gqlgen@v0.17.39

go run github.com/99designs/gqlgen generate

  

  

-------------------------------------------------------------------------------------------------------

  

  

**Kafka**

  

  

Kafka version > 2.2:

  

/home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic TutorialTopic

create topic:

  

/home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-topics.sh --create --topic TutorialTopic --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1

  

(https://stackoverflow.com/questions/69297020/exception-in-thread-main-joptsimple-unrecognizedoptionexception-zookeeper-is)

  

  

send the message to topic:

  

cat msg.json | /home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic FIO > /dev/null

  

  

read a message from topic:

  

/home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic TutorialTopic --from-beginning

  

  

---------------------------------------------------------------------------------

  

Manuals and articles thats helped and have been notable useful for me:

  

 - https://www.digitalocean.com/community/tutorials/how-to-install-apache-kafka-on-ubuntu-20-04
   
     
   
   https://www.sohamkamani.com/golang/working-with-kafka/
   
     
   
   https://www.digitalocean.com/community/tutorials/how-to-write-unit-tests-in-go-using-go-test-and-the-testing-package
   
     
   
   https://betterprogramming.pub/building-a-graphql-server-using-the-schema-first-approach-in-golang-a8da71d7e5b7
   
     
   
   https://hevodata.com/learn/postgresql-partitions/#t8
   
     
   
   https://golangbyexample.com/load-env-fiie-golang/
   
     
   
   https://go.dev/doc/tutorial/web-service-gin

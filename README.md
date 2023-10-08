Осталось:

8. Покрыть код логами

9. Покрыть бизнес-логику unit-тестами

https://www.digitalocean.com/community/tutorials/how-to-write-unit-tests-in-go-using-go-test-and-the-testing-package

10. Вынести все конфигурационные данные в .env

----

  

12. разнести все по модулям

12. code review, DRY, KISS, SOLID

  

-------------------------------------------------------------------------------------------------------

  

**GraphQL**

  

https://betterprogramming.pub/building-a-graphql-server-using-the-schema-first-approach-in-golang-a8da71d7e5b7

  

Для того чтобы сгенерировать схему нужно выполнить команды:

  

    export GOPATH=/home/kostya/go/src
    
      
    
    go mod init
    
      
    
    go get github.com/99designs/gqlgen
    
      
    
    go run github.com/99designs/gqlgen init
    
      
    
    go mod tidy
    
      
    
    go run github.com/99designs/gqlgen generate

  

При генерации могут полезть ошибки, нужно установить пакеты:

    2062 go get github.com/99designs/gqlgen/codegen/config@v0.17.39
    
    2063 go get github.com/99designs/gqlgen/internal/imports@v0.17.39
    
    2064 go get github.com/99designs/gqlgen@v0.17.39
    
    2065 go run github.com/99designs/gqlgen generate

  

-------------------------------------------------------------------------------------------------------

  

**Kafka**

  

На версиях Кафки старше чем 2.2:

    /home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic TutorialTopic
    
      

поэтому создаем топик так:

    /home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-topics.sh --create --topic TutorialTopic --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1

(https://stackoverflow.com/questions/69297020/exception-in-thread-main-joptsimple-unrecognizedoptionexception-zookeeper-is)

  

сообщение отправить:

    cat msg.json | /home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic FIO > /dev/null

  

сообщение принять:

    /home/kostya/kafka/kafka_2.12-3.6.0/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic TutorialTopic --from-beginning

  
  

---------------------------------------------------------------------------------

Сатьи, которые использоваляиь для написания проекта:

  

https://www.digitalocean.com/community/tutorials/how-to-install-apache-kafka-on-ubuntu-20-04

https://www.sohamkamani.com/golang/working-with-kafka/

https://www.digitalocean.com/community/tutorials/how-to-write-unit-tests-in-go-using-go-test-and-the-testing-package

https://betterprogramming.pub/building-a-graphql-server-using-the-schema-first-approach-in-golang-a8da71d7e5b7

https://hevodata.com/learn/postgresql-partitions/#t8


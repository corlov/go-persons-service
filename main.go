package main

import (
	"data_enrichment"
	"graphql"
	"rest"
	"time"
)

func main() {    
    go rest.ServiceRun()
    go graphql.ServiceRun()
    go data_enrichment.ServiceRun()

    // TODO: проверяем живы сервисы или нет. Если нет, то перезапуск умершего сервиса
    for {        
        time.Sleep(time.Second) 
    }
}


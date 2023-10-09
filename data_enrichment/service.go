package data_enrichment

import (
	"context"
	"database/sql"
	"db_utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"types"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"

	_ "github.com/lib/pq"
)


const (
	fioTopicName         	= "FIO"
	errrorTopicName 		= "FIO_FAILED"
)
const END_POINT_NATION = "https://api.nationalize.io"
const END_POINT_GENDER = "https://api.genderize.io"
const END_POINT_AGE    = "https://api.agify.io"

var Log *log.Logger
var brokerAddress string


func ServiceRun() {
	file, err := os.Create("enrichment.log")
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags | log.Lshortfile)
	Log.Println("started")

	err = godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	brokerAddress = os.Getenv("KAFKA_ADDR")
    

	ctx := context.Background()	
	consume(ctx)
}


func consume(ctx context.Context) {
	// create a new logger that outputs to stdout
	// and has the `kafka reader` prefix
	
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   fioTopicName,
		GroupID: "my-group",
		Logger: Log,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			Log.Println("could not read message " + err.Error())
			continue
		}
		Log.Println("received: ", string(msg.Value))

		
		var newPerson types.Person
		err = json.Unmarshal([]byte(string(msg.Value)), &newPerson)
		if err != nil {
			Log.Println("json format error:", err)
			go produceToErrorQueue(string(msg.Value))
		} else {
			Log.Println("name: '" + newPerson.Name + "'\n", "surname: '" + newPerson.Surname + "'\n", "Patronymic: '" + newPerson.Patronymic + "'\n",)

			if (newPerson.Name == "") || (newPerson.Surname == "") {  
				go produceToErrorQueue(string(msg.Value))			
				Log.Println("Empty fields error!")
			} else {
				go expandData(&newPerson)
			}
		}
	}
}



func expandData(p *types.Person) {		
	var age types.Age
	jsonStr := httpRes(END_POINT_AGE, p.Name)
	err := json.Unmarshal([]byte(jsonStr), &age)
	if err != nil {
		Log.Println(err.Error())
		return
	}
	// FIXME: может вернуть в случае экзотического имени пустое значение
	// поэтому проверку типов пройти, если null, то тогда выдать сообщение о некорректности имени и т.д.
	p.Age = age.Age

	var gender types.Gender
	jsonStr = httpRes(END_POINT_GENDER, p.Name)
	err = json.Unmarshal([]byte(jsonStr), &gender)
	if err != nil {
		Log.Println(err.Error())
		return
	}
	// FIXME: может вернуть в случае экзотического имени пустое значение
	// поэтому проверку типов пройти, если null, то тогда выдать сообщение о некорректности имени и т.д.
	p.Gender = gender.Gender


	var nationality types.Nationality
	jsonStr = httpRes(END_POINT_NATION, p.Name)
	err = json.Unmarshal([]byte(jsonStr), &nationality)
	if err != nil {
		Log.Println(err.Error())
		return
	}
	p.Nationality = "??"
	max := 0.0
	for _, nation := range nationality.Country { 
		if nation.Probability > max {
			max = nation.Probability
			p.Nationality = nation.CountryId
		}
	}
	
	Log.Println(p)
	insertDb(p)
}


func insertDb(p *types.Person) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    						db_utils.Host, db_utils.Port, db_utils.User, db_utils.Password, db_utils.DbName)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()
     
	
    sqlStatement := `INSERT INTO "Population".Person (name, surname, patronymic, age, country_id, gender_id) VALUES ($1, $2, $3, $4, $5, $6)`
    _, err = db.Exec(sqlStatement, p.Name, p.Surname, p.Patronymic, p.Age, p.Nationality, p.Gender)
    if err != nil {
        panic(err)
    } else {
        Log.Println("\nRow inserted successfully!")
    }
}


func httpRes(baseURL string, personName string) string {
	params := url.Values{}
	params.Add("name", personName)
	u, _ := url.ParseRequestURI(baseURL)				
	u.RawQuery = params.Encode()
	
	urlStr := fmt.Sprintf("%v", u) 
	Log.Println("query to ", urlStr)

	resp, err := http.Get(urlStr)
	defer resp.Body.Close()
	Log.Println(resp.Status)
	if err != nil {
		Log.Println(err.Error())
		return ""
	}
	
	if resp.Status != "200 OK" {
		Log.Println(resp.Status)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil {
		Log.Println(err.Error())
		return ""
	}
	Log.Println(string(body), err) 

	return string(body)
}



func produceToErrorQueue(message string) {
	ctx := context.Background()

	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   errrorTopicName,
		// assign the logger to the writer
		Logger: Log,
	})

	// each kafka message has a key and value. The key is used
	// to decide which partition (and consequently, which broker)
	// the message gets published on
	err := w.WriteMessages(ctx, kafka.Message{ Key: []byte("0"), Value: []byte(message), })
	if err != nil {
		panic("could not write message " + err.Error())
	}

	w.Close()	
}

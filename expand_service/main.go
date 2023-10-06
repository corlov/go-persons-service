package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/segmentio/kafka-go"

	_ "github.com/lib/pq"
)

const (
    host = "127.0.0.1"
    port = 5432
    user = "postgres"
    password = "postgres"
    dbname = "WorldPopulation"
)

const (
	fioTopicName         = "FIO"
	errrorTopicName = "FIO_FAILED"
	brokerAddress = "localhost:9092"
)

type Person struct {
	Name  string `json:"name"`
	Surname string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Age uint8
	Gender string
	Nationality string
}

type Age struct {
	Count  int `json:"count"`
	Name string `json:"name"`
	Age uint8 `json:"age"`	
}

type Gender struct {
	Count  int `json:"count"`
	Name string `json:"name"`
	Gender string `json:"gender"`
	Probability float32 `json:"probability"`
}

type Country struct {
	CountryId  string `json:"country_id"`
	Probability float64 `json:"probability"`
}

type Nationality struct {
	Count  int `json:"count"`
	Name string `json:"name"`
	Country []Country
}



func main() {
	ctx := context.Background()	
	consume(ctx)

	// p := Person{Name: "Dmitry", Surname: "Petrov", Patronymic: "Olegovich"}
	// expandData(&p)

}


func consume(ctx context.Context) {
	// create a new logger that outputs to stdout
	// and has the `kafka reader` prefix
	l := log.New(os.Stdout, "kafka reader: ", 0)
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   fioTopicName,
		GroupID: "my-group",
		// assign the logger to the reader
		Logger: l,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		fmt.Println("received: ", string(msg.Value))

		
		var newPerson Person
		err = json.Unmarshal([]byte(string(msg.Value)), &newPerson)
		if err != nil {
			fmt.Println("json format error:", err)
			go produceToErrorQueue(string(msg.Value))
		} else {
			fmt.Println("name: '" + newPerson.Name + "'\n", "surname: '" + newPerson.Surname + "'\n", "Patronymic: '" + newPerson.Patronymic + "'\n",)

			if (newPerson.Name == "") || (newPerson.Surname == "") {  
				go produceToErrorQueue(string(msg.Value))			
				fmt.Println("Empty fields error!")
			} else {
				go expandData(&newPerson)
			}
		}
	}
}



func expandData(p *Person) {	
	var age Age
	jsonStr := httpRes("https://api.agify.io", p.Name)
	err := json.Unmarshal([]byte(jsonStr), &age)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// FIXME: проверку типов пройти, если null, то тогда выдать сообщение о некорректности имени и т.д.
	p.Age = age.Age

	var gender Gender
	jsonStr = httpRes("https://api.genderize.io", p.Name)
	err = json.Unmarshal([]byte(jsonStr), &gender)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	p.Gender = gender.Gender


	var nationality Nationality
	jsonStr = httpRes("https://api.nationalize.io", p.Name)
	err = json.Unmarshal([]byte(jsonStr), &nationality)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	p.Nationality = "unknown"
	max := 0.0
	for _, nation := range nationality.Country { 
		if nation.Probability > max {
			max = nation.Probability
			p.Nationality = nation.CountryId
		}
	}
	
	fmt.Println(p)

	insertDb(p)
}


func insertDb(p *Person) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()
     
	
    sqlStatement := `INSERT INTO "Population".Person (name, surname, patronymic, age, country_id, gender_id) VALUES ($1, $2, $3, $4, $5, $6)`
    _, err = db.Exec(sqlStatement, p.Name, p.Surname, p.Patronymic, p.Age, p.Nationality, p.Nationality)
    if err != nil {
        panic(err)
    } else {
        fmt.Println("\nRow inserted successfully!")
    }
}


func httpRes(baseURL string, personName string) string {
	params := url.Values{}
	params.Add("name", personName)
	u, _ := url.ParseRequestURI(baseURL)				
	u.RawQuery = params.Encode()
	
	urlStr := fmt.Sprintf("%v", u) 
	fmt.Println("query to ", urlStr)

	resp, err := http.Get(urlStr)
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	
	if resp.Status != "200 OK" {
		fmt.Println(resp.Status)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	fmt.Println(string(body), err) 

	return string(body)
}



func produceToErrorQueue(message string) {
	ctx := context.Background()

	l := log.New(os.Stdout, "kafka writer: ", 0)
	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   errrorTopicName,
		// assign the logger to the writer
		Logger: l,
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


//
// https://stackoverflow.com/questions/49585021/how-to-check-if-a-json-matches-a-struct-struct-fields
//
// func (p *Person) UnmarshalJSON(data []byte) error {
//     var m map[string]interface{}
//     err := json.Unmarshal(data, &m)
//     if err != nil {
//         return err
//     }

//     v := reflect.ValueOf(p).Elem()
//     t := v.Type()

//     var missing []string
//     for i := 0; i < t.NumField(); i++ {
//         field := t.Field(i)
//         val, ok := m[field.Name]
//         delete(m, field.Name)
//         if !ok {
//             missing = append(missing, field.Name)
//             continue
//         }

//         switch field.Type.Kind() {
//         // TODO: if the field is an integer you need to transform the val from float
//         default:
//             v.Field(i).Set(reflect.ValueOf(val))
//         }
//     }

//     if len(missing) > 0 {
//         return errors.New("missing fields: " + strings.Join(missing, ", "))
//     }

//     if len(m) > 0 {
//         extra := make([]string, 0, len(m))
//         for field := range m {
//             extra = append(extra, field)
//         }
//         // TODO: consider sorting the output to get deterministic errors:
//         // sort.Strings(extra)
//         return errors.New("unknown fields: " + strings.Join(extra, ", "))
//     }

//     return nil
// }
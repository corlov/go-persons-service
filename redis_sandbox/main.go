package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis"
)

// FIXME: в общие структурные файлы структуры вынести
const (
    host = "127.0.0.1"
    port = 5432
    user = "postgres"
    password = "postgres"
    dbname = "WorldPopulation"
)

// FIXME: в общие структурные файлы структуры вынести
type Person struct {
	Id  uint64 `json:"id"`
	Name  string `json:"name"`
	Surname string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Age uint8 `json:"age"`
	Gender string `json:"gender"`
	Nationality string `json:"nationality"`
}

func main() {
	loadDb2Redis()

    // client := redis.NewClient(&redis.Options{
    //     Addr:	  "localhost:6379",
    //     Password: "", // no password set
    //     DB:		  0,  // use default DB
    // })

	// ctx := context.Background()

	// err := client.Set(ctx, "foo", "bar", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// val, err := client.Get(ctx, "foo").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("foo", val)
}


func loadDb2Redis() {

	client := redis.NewClient(&redis.Options{
        Addr:	  "localhost:6379",
        Password: "", // no password set
        DB:		  0,  // use default DB
    })

	ctx := context.Background()


	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
    if err != nil {
		return
    }
    defer db.Close()
    
	// Avoiding SQL injection risk
	rows, err := db.Query(`
		select  
			id,
			name, 
			surname, 
			coalesce(patronymic, '') as patronymic,
			coalesce(age, 0) as age, 
			coalesce(gender_id, '') as gender,
			coalesce(country_id, '') as nationality
		from "Population".Person`)

    if err != nil {
		return
    }
    defer rows.Close()

    for rows.Next() {
		var p Person
        err := rows.Scan(&p.Id, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
        if err != nil {
			return
        }
		
		jsonText, err := json.Marshal(p)
		err = client.Set(ctx, strconv.Itoa(int(p.Id)), jsonText, 0).Err()
		if err != nil {
			panic(err)
		}
    }
    err = rows.Err()
    if err != nil {
		return
    }
	fmt.Println("Loaded!")

	// reading
	
	// val, err := client.Get(ctx, "68").Result()
	// if err != nil {
	// 	fmt.Println("Not found")
	// 	panic(err)
	// }
	// fmt.Println("6", val)

	// var p Person
	// err = json.Unmarshal([]byte(val), &p)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(p)

}
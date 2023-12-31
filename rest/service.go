package rest

import (
	"database/sql"
	"db_utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"types"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Log *log.Logger
var redisAddr string



func ServiceRun() {
	file, err := os.Create("rest_service.log")
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags | log.Lshortfile)
	Log.Println("started")

	err = godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	redisAddr = os.Getenv("REDIS_ADDR")
	restPort := os.Getenv("REST_PORT")

	loadDb2Redis()
	
    router := gin.Default()    
	router.GET("/get_persons", getPersons)	
	router.POST("/add_person", addPerson)
	router.GET("/remove_person", removePerson)
	router.POST("/update_person", updatePerson)
    router.Run("localhost:" + restPort)
}


/*
curl http://localhost:8080/remove_person?id=4
*/
func removePerson(c *gin.Context) {
	argId := 0

	// TODO: предусмотреть удаление по другим параметрам также
	paramPairs := c.Request.URL.Query()
    for key, val := range paramPairs {
		switch key {
			case "id":	
				var err error
				argId, err = strconv.Atoi(val[0])
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}	
		}
    }

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable", 
							db_utils.Host, db_utils.Port, db_utils.User, db_utils.Password, db_utils.DbName)
	db, err := sql.Open("postgres", connStr)
    if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
    }
    defer db.Close()
    
	// Avoiding SQL injection risk
	sqlQueryText := `delete from "Population".Person where id = $1`
	rows, err := db.Query(sqlQueryText, argId)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
    }
    defer rows.Close()
	c.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})   
}



/*
curl http://localhost:8080/update_person \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{ "id": 2, "name": "Alex", "surname": "Ivanov", "age": 67, "gender": "male", "nationality": "ES"}'
*/
func updatePerson(c *gin.Context) {
	var updPerson types.Person
    // Call BindJSON to bind the received JSON to newPerson
    if err := c.BindJSON(&updPerson); err != nil {
		Log.Println("BIND ERROR")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }
	Log.Println(updPerson)

	errorMsg := update(updPerson)
	if errorMsg != "" {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": errorMsg})
	}

	c.IndentedJSON(http.StatusCreated, updPerson)
}


func update(p types.Person) string {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    						db_utils.Host, db_utils.Port, db_utils.User, db_utils.Password, db_utils.DbName)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
		return err.Error()
    }
    defer db.Close()
     	
    sqlStatement := `
		UPDATE "Population".Person SET 
			name = $1,
			surname = $2,
			patronymic = $3, 
			age = $4,
			country_id = $5,
			gender_id = $6
		WHERE (id = $7)`
    res, err := db.Exec(sqlStatement, p.Name, p.Surname, p.Patronymic, p.Age, p.Nationality, p.Gender, p.Id)
    if err != nil {
        return err.Error()
    } 

	n, err := res.RowsAffected()
	if err != nil {
		return err.Error()
	}

	if n < 1 {
		return "not found"
	}
	
	return ""
}

/*
curl http://localhost:8080/add_person \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{ "name": "Alex", "surname": "Ivanov", "age": 31, "gender": "male", "nationality": "RU" }'
*/
func addPerson(c *gin.Context) {
	var newPerson types.Person
    // Call BindJSON to bind the received JSON to newPerson
    if err := c.BindJSON(&newPerson); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    	db_utils.Host, db_utils.Port, db_utils.User, db_utils.Password, db_utils.DbName)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
    }
    defer db.Close()
     	
    sqlStatement := `INSERT INTO "Population".Person (name, surname, patronymic, age, country_id, gender_id) VALUES ($1, $2, $3, $4, $5, $6)`
    _, err = db.Exec(sqlStatement, newPerson.Name, newPerson.Surname, newPerson.Patronymic, newPerson.Age, newPerson.Nationality, newPerson.Gender)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
    } 
    c.IndentedJSON(http.StatusCreated, newPerson)
}



func getPersons(c *gin.Context) {
	// для пагинации результатов
	blockSize := 10
	page := 0

	var (
		argId uint64 = 0
		argName string
		argSurname string
		argPatronymic string
		argAge int = 0
		argGender string
		argNationality string
	)

	// если запрос только по ид, то считываем из кеша (Redis) иначе из БД
	requestByIdOnly := true

	paramPairs := c.Request.URL.Query()
    for key, val := range paramPairs {
		switch key {
			case "id":
				id, err := strconv.Atoi(val[0])
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
				argId = uint64(id)

			case "name":
				argName = val[0]		
				requestByIdOnly = false	

			case "surname":
				argName = val[0]
				requestByIdOnly = false

			case "patronymic":
				argPatronymic = val[0]
			case "age":
				var err error
				argAge, err = strconv.Atoi(val[0])
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
				requestByIdOnly = false

			case "gender":
				argGender = val[0]
				requestByIdOnly = false

			case "nationality":
				argNationality = val[0]
				requestByIdOnly = false

			case "page":
				pageVal, err := strconv.Atoi(string(val[0]))
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
				page = 	blockSize * pageVal
		}
    }

	if requestByIdOnly {
		fmt.Println(requestByIdOnly)
		
		client := redis.NewClient(&redis.Options{
			Addr:	  redisAddr,
			Password: "", // no password set
			DB:		  0,  // use default DB
		})

		val, err := client.Get(strconv.Itoa(int(argId))).Result()
		if err != nil {
			Log.Println("Not found, read from Db")
		} else {
			var p types.Person
			err = json.Unmarshal([]byte(val), &p)
			if err != nil {
				Log.Println(err.Error())
				return
			}
			Log.Println(p)

			var persons []types.Person
			persons = append(persons, p)
			c.IndentedJSON(http.StatusOK,  persons) 
			Log.Println("found into Redis")
			return
		}
	}

	sqlWhere := `
		(name = $1       or NULLIF($1, '') is null) and
		(surname = $2    or NULLIF($2, '') is null) and
		(patronymic = $3 or NULLIF($3, '') is null) and
		(age = $4        or NULLIF($4,  0) is null) and
		(gender_id = $5  or NULLIF($5, '') is null) and
		(country_id = $6 or NULLIF($6, '') is null) and
		(id = $7         or NULLIF($7, 0) is null)`

	sqlLimit := " limit " + strconv.Itoa(blockSize)
	sqlOffset := " offset " + strconv.Itoa(page)
	
	sqlQueryText := `
		select  
			id,
			name, 
			surname, 
			coalesce(patronymic, '') as patronymic,
			coalesce(age, 0) as age, 
			coalesce(gender_id, '') as gender,
			coalesce(country_id, '') as nationality
		from "Population".Person where ` + sqlWhere + sqlLimit + sqlOffset


	// fmt.Println(sqlQueryText)
	
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    						db_utils.Host, db_utils.Port, db_utils.User, db_utils.Password, db_utils.DbName)
	db, err := sql.Open("postgres", connStr)
    if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
    }
    defer db.Close()
    
	// Avoiding SQL injection risk
	rows, err := db.Query(sqlQueryText, argName, argSurname, argPatronymic, argAge, argGender, argNationality, argId)

    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
    }
    defer rows.Close()

	var persons []types.Person
    for rows.Next() {
		var p types.Person
        err := rows.Scan(&p.Id, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
        if err != nil {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
        }
		persons = append(persons, p)
    }
    err = rows.Err()
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
    }

	c.IndentedJSON(http.StatusOK,  persons)   
}


func loadDb2Redis() {

	client := redis.NewClient(&redis.Options{
        Addr:	  redisAddr,
        Password: "", // no password set
        DB:		  0,  // use default DB
    })

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    					db_utils.Host, db_utils.Port, db_utils.User, db_utils.Password, db_utils.DbName)
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
		Log.Println("Load from RedisDb failed!")
		return
    }
    defer rows.Close()

    for rows.Next() {
		var p types.Person
        err := rows.Scan(&p.Id, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
        if err != nil {
			return
        }
		
		jsonText, err := json.Marshal(p)
		if err != nil {
			Log.Println("Load from RedisDb failed!", err.Error())
			return
		}
		err = client.Set(strconv.Itoa(int(p.Id)), jsonText, 0).Err()
		if err != nil {
			panic(err)
		}
		val, err := client.Get(strconv.Itoa(int(p.Id))).Result()
		if err != nil {
			Log.Println("Load from RedisDb failed!", err.Error())
			return
		}
		fmt.Println("Loaded ", val)
    }
    err = rows.Err()
    if err != nil {
		Log.Println("Load from RedisDb failed!", err.Error())
		return
    }
	Log.Println("All records have loaded!")

	// reading from Redis

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
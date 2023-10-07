package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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
    router := gin.Default()
    
	router.GET("/get_persons", getPersons)	
	router.POST("/add_person", addPerson)
	router.GET("/remove_person", removePerson)
	router.POST("/update_person", updatePerson)

	// FIXME: в константы env
    router.Run("localhost:8080")
}


/*
curl http://localhost:8080/remove_person?id=4
*/
func removePerson(c *gin.Context) {
	argId := 0

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

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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
	var updPerson Person
    // Call BindJSON to bind the received JSON to newPerson
    if err := c.BindJSON(&updPerson); err != nil {
		fmt.Println("BIND ERROR")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }
	fmt.Println(updPerson)

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
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
    _, err = db.Exec(sqlStatement, updPerson.Name, updPerson.Surname, updPerson.Patronymic, updPerson.Age, updPerson.Nationality, updPerson.Gender, updPerson.Id)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
    } 
    c.IndentedJSON(http.StatusCreated, updPerson)
}



/*
curl http://localhost:8080/add_person \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{ "name": "Alex", "surname": "Ivanov", "age": 31, "gender": "male", "nationality": "RU" }'
*/
func addPerson(c *gin.Context) {
	var newPerson Person
    // Call BindJSON to bind the received JSON to newPerson
    if err := c.BindJSON(&newPerson); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
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
	// FIXME: в константы
	blockSize := 10
	page := 0

	var (
		argName string
		argSurname string
		argPatronymic string
		argAge int = 0
		argGender string
		argNationality string
	)

	paramPairs := c.Request.URL.Query()
    for key, val := range paramPairs {
		switch key {
			case "name":
				argName = val[0]			
			case "surname":
				argName = val[0]
			case "patronymic":
				argPatronymic = val[0]
			case "age":
				var err error
				argAge, err = strconv.Atoi(val[0])
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
			case "gender":
				argGender = val[0]
			case "nationality":
				argNationality = val[0]
			case "page":
				pageVal, err := strconv.Atoi(string(val[0]))
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
				page = 	blockSize * pageVal
		}
    }

	sqlWhere := `
		(name = $1       or NULLIF($1, '') is null) and
		(surname = $2    or NULLIF($2, '') is null) and
		(patronymic = $3 or NULLIF($3, '') is null) and
		(age = $4        or NULLIF($4,  0) is null) and
		(gender_id = $5  or NULLIF($5, '') is null) and
		(country_id = $6 or NULLIF($6, '') is null)`

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
    host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
    if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
    }
    defer db.Close()
    
	// Avoiding SQL injection risk
	rows, err := db.Query(sqlQueryText, argName, argSurname, argPatronymic, argAge, argGender, argNationality)

    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
    }
    defer rows.Close()

	var persons []Person
    for rows.Next() {
		var p Person
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

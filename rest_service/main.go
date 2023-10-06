package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
    host = "127.0.0.1"
    port = 5432
    user = "postgres"
    password = "postgres"
    dbname = "WorldPopulation"
)


type Person struct {
	Name  string `json:"name"`
	Surname string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Age uint8 `json:"age"`
	Gender string `json:"gender"`
	Nationality string `json:"nationality"`
}

// album represents data about a record album.
type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
    router := gin.Default()
    router.GET("/get_persons", getPersons)
	router.GET("/get_person/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

    router.Run("localhost:8080")
}



func NewNullString(s string) sql.NullString {
    if len(s) == 0 {
        return sql.NullString{}
    }
    return sql.NullString{
         String: s,
         Valid: true,
    }
}

func getPersons(c *gin.Context) {
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
				fmt.Println("AGE: ", argAge)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
			case "gender":
				argGender = val[0]
			case "nationality":
				argNationality = val[0]
			case "page":
				pageVal, _ := strconv.Atoi(string(val[0]))
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
        err := rows.Scan(&p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
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


// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
    id := c.Param("id")

    // Loop over the list of albums, looking for
    // an album whose ID value matches the parameter.
    for _, a := range albums {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
    var newAlbum album

    // Call BindJSON to bind the received JSON to
    // newAlbum.
    if err := c.BindJSON(&newAlbum); err != nil {
        return
    }

    // Add the new album to the slice.
    albums = append(albums, newAlbum)
    c.IndentedJSON(http.StatusCreated, newAlbum)
}
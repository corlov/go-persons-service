package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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


func getPersons(c *gin.Context) {
	blockSize := 10
	page := 0

	sqlWhere := " (1 = 1) "

	paramPairs := c.Request.URL.Query()
    for key, val := range paramPairs {
        fmt.Printf("key: %v, value: %v\n", key, val)
		switch key {
			case "age":
				ageVal, err := strconv.Atoi(string(val[0]))
				if (err != nil) {
					c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error"})
					return 
				}
				sqlWhere += "and (age = " + strconv.Itoa(ageVal) + ") "	

				// FIXME: защита от инъекции SQL
			case "name":
				sqlWhere += "and (name = '" + val[0] + "') "
			
			case "surname":
				sqlWhere += "and (surname = '" + val[0] + "') "	

			case "page":
				pageVal, _ := strconv.Atoi(string(val[0]))
				page = 	blockSize * pageVal
		}
    }

	sqlLimit := " limit " + strconv.Itoa(blockSize)
	sqlOffset := " offset " + strconv.Itoa(page)

	sqlQueryText := "select name, surname, age from \"Population\".Person where " + sqlWhere + sqlLimit + sqlOffset


	fmt.Println(sqlQueryText)
	

	// делаем запрос к БД

    //c.IndentedJSON(http.StatusOK, albums)
	// c.JSON(http.StatusOK, albums)

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": sqlQueryText})
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
package graph

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"graphql_service/graphql_service/graph/model"
	"strconv"

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


// AddPerson is the resolver for the add_person field.
func (r *mutationResolver) AddPerson(ctx context.Context, input model.PersonInput) (*model.PostStatus, error) {	
	description := ""	
	newPerson := model.Person{Name: input.Name, Surname: input.Surname, Patronymic: input.Patronymic, Age: input.Age, Gender: input.Gender, Nationality: input.Nationality}
   
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        description = err.Error()
		return &model.PostStatus{Iserror: true, Description: &description}, errors.New(description)
    }
    defer db.Close()
     	
    sqlStatement := `INSERT INTO "Population".Person (name, surname, patronymic, age, country_id, gender_id) VALUES ($1, $2, $3, $4, $5, $6)`
    _, err = db.Exec(sqlStatement, newPerson.Name, newPerson.Surname, newPerson.Patronymic, newPerson.Age, newPerson.Nationality, newPerson.Gender)
    if err != nil {
        description = err.Error()
		return &model.PostStatus{Iserror: true, Description: &description}, errors.New(description)
    } 
    
	return &model.PostStatus{Iserror: false, Description: &description}, errors.New(description)
}

// UpdatePerson is the resolver for the update_person field.
func (r *mutationResolver) UpdatePerson(ctx context.Context, input *model.UpdatePersonInput) (*model.PostStatus, error) {
	description := ""
	newPerson := model.UpdatePersonInput{ID: input.ID, Name: input.Name, Surname: input.Surname, Patronymic: input.Patronymic, Age: input.Age, Gender: input.Gender, Nationality: input.Nationality}
   
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        description = err.Error()
		return &model.PostStatus{Iserror: true, Description: &description}, errors.New(description)
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
    _, err = db.Exec(sqlStatement, newPerson.Name, newPerson.Surname, newPerson.Patronymic, newPerson.Age, newPerson.Nationality, newPerson.Gender, newPerson.ID)
    if err != nil {
        description = err.Error()
		return &model.PostStatus{Iserror: true, Description: &description}, errors.New(description)
    } 
    
	return &model.PostStatus{Iserror: false, Description: &description}, errors.New(description)
}

// DeletePerson is the resolver for the delete_person field.
func (r *mutationResolver) DeletePerson(ctx context.Context, personID string) (*model.PostStatus, error) {
	description := ""
	
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
    if err != nil {
		description = err.Error()
		return &model.PostStatus{Iserror: true, Description: &description}, errors.New(description)
    }
    defer db.Close()
    
	personIDAsInt, err := strconv.Atoi(personID)
	if err != nil {
        description = err.Error()
		return &model.PostStatus{Iserror: true, Description: &description}, errors.New(description)
    }

	sqlQueryText := `delete from "Population".Person where id = $1`
	rows, err := db.Query(sqlQueryText, personIDAsInt)
    if err != nil {
        description = err.Error()
		return &model.PostStatus{Iserror: true, Description: &description}, errors.New(description)
    }
    defer rows.Close()
	return &model.PostStatus{Iserror: false, Description: &description}, errors.New(description)
}

// FIXME: добавить в параметры пагинацию, другие поля по которым делать выборку
// GetPersons is the resolver for the get_persons field.
func (r *queryResolver) GetPersons(ctx context.Context, id string) ([]*model.Person, error) {	
	personIDAsInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, nil
    }

	sqlQueryText := `
		select  
			id,
			name, 
			surname, 
			coalesce(patronymic, '') as patronymic,
			coalesce(age, 0) as age, 
			coalesce(gender_id, '') as gender,
			coalesce(country_id, '') as nationality
		from "Population".Person where id = $1`

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
    if err != nil {
		return nil, nil
    }
    defer db.Close()

	// Avoiding SQL injection risk
	rows, err := db.Query(sqlQueryText, personIDAsInt)

    if err != nil {
        return nil, nil
    }
    defer rows.Close()

	var persons []*model.Person
    for rows.Next() {
		var p model.Person
        err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
        if err != nil {
            return nil, nil
        }
		persons = append(persons, &p)
    }
    err = rows.Err()
    if err != nil {
        return nil, nil
    }

	return persons, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

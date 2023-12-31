package rest

import (
	"fmt"
	"testing"
	"types"
)

func TestUpdate_1(t *testing.T){

	person := types.Person{Id: 9, Name: "Igor", Surname: "Smirnoff", Age: 67, Gender: "male", Nationality: "IZ"}
    got := update(person)
    want := ""
    if got != want {
        t.Errorf("got %q, wanted %q", got, want)
    }
}


func TestUpdate_2(t *testing.T){
	person := types.Person{Id: 19, Name: "Igor", Surname: "Smirnoff", Age: 67, Gender: "male", Nationality: "IZ"}
    got := update(person)
    want := "not found"
    if got != want {
        t.Errorf("got %q, wanted %q", got, want)
    }
}


func TestUpdate_3(t *testing.T){
	fmt.Println("test3")
	person := types.Person{Age: 67, Gender: "male", Nationality: "IZ"}
    got := update(person)
    want := "not found"
    if got != want {
        t.Errorf("got %q, wanted %q", got, want)
    }
}
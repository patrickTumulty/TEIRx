package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	Name string
	Age  int
}

// Method to display details of the Person
func (p Person) Greet() {
	fmt.Printf("Hello, my name is %s and I am %d years old.\n", p.Name, p.Age)
}

type User struct {
	FirstName   string
	LastName    string
	DateOfBirth string // yyyy-MM-dd
	ZipCode     string
	Email       string
}

func (usr User) ToSQLString() string {
	return fmt.Sprintf("'%s', '%s', '%s', '%s'", usr.FirstName, usr.LastName, usr.DateOfBirth, usr.ZipCode)
}

func main() {

	dsn := "ptumulty:pass@tcp(127.0.0.1:3306)/reviewdb"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	usr := User{"Patrick", "Tumulty", "1994-05-06", "91602", "website@email.com"}

	insert, err := db.Query(fmt.Sprintf("INSERT INTO users (firstname, lastname, date_of_birth, zipcode, email) VALUES (%s);", usr.ToSQLString()))

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}

	// be careful deferring Queries if you are using transactions
	defer insert.Close()
}

package config

import (
	controller "go_prog/application_stucture/controllers"
	"log"
	"os"

	"github.com/go-pg/pg/v9"
)

func Connect() *pg.DB {
	opts := &pg.Options{
		User:     "postgres",
		Password: "postgres",
		Addr:     "localhost:5432",
		Database: "students",
	}

	var db *pg.DB = pg.Connect(opts)
	controller.CreateTodoTable(db)
	controller.InitiateDB(db)
	if db == nil {
		log.Printf("Failed to Connect")
		os.Exit(100)
	}
	log.Printf("Connected hurray")

	return db

}

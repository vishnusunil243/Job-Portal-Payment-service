package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vishnusunil243/Job-Portal-Payment-service/db"
	"github.com/vishnusunil243/Job-Portal-Payment-service/initializer"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("error loading env file")
	}
	addr := os.Getenv("DB_KEY")
	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatal("error initialising database")
	}
	initializer.Initializer(DB).Start(":8089")

}

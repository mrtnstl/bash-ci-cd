package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	port := getEnvString("PORT")
	api := NewApplication(port)
	
	if err := api.Run(); err != nil {
		log.Fatalf("%s", err)
	}
}
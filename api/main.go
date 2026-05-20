package main

import (
	"log"

	"example.com/api/cmd"
)

func main(){
	api := cmd.NewApplication(":8080")
	
	if err := api.Run(); err != nil {
		log.Fatalf("%s", err)
	}
}
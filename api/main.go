package main

import (
	"embed"
	"log"

	"example.com/api/cmd/api"
	"example.com/api/internals/utils"
	"github.com/joho/godotenv"
)

//go:embed swagger.html
//go:embed openapi.yaml
var SwaggerUI embed.FS

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	port := utils.GetEnvString("PORT")
	api := api.NewApplication(api.Config{
		Addr: port,
		Static: SwaggerUI,
	})

	if err := api.Run(); err != nil {
		log.Fatalf("%s", err)
	}
}

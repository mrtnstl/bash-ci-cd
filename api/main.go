package main

import (
	"embed"
	"log"
	"os"

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
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error while getting current working directory: %v", err)
	}

	api := api.NewApplication(api.Config{
		Addr: port,
		Static: SwaggerUI,
		AccessLogLocation: pwd + "/access_log.txt",
	})

	if err := api.Run(); err != nil {
		log.Fatalf("%s", err)
	}
}

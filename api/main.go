package main

import (
	"context"
	"embed"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"example.com/api/cmd/api"
	"example.com/api/internals/utils"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

//go:embed swagger.html
//go:embed openapi.yaml
var SwaggerUI embed.FS

func main() {
	if utils.GetEnvString(utils.GO_ENV) != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading .env file")
		}
	}

	port := utils.GetEnvString(utils.PORT)
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error while getting current working directory: %v", err)
	}

	// chanel for notifying goroutines of server shutdown
	shutdownChan := make(chan bool)

	var wg sync.WaitGroup

	api := api.NewApplication(
		api.Config{
			Addr:              port,
			Static:            SwaggerUI,
			AccessLogLocation: pwd + "/access_log.txt",
			RlLimit: rate.Limit(1),
			RlBurst: 3,
		},
		&wg,
		&shutdownChan,
	)

	server := api.CreateServer()

	go func() {
		if err := api.Run(server); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
		log.Println("http server stopped serving new connections")
	}()

	// sigint, sigterm signal listener
	shutdownSig, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-shutdownSig.Done()

	// shutdown pipeline
	log.Println("notifying background jobs of shutdown")

	select {
	case shutdownChan <- true:
	default:
	}

	log.Println("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		if err == context.DeadlineExceeded {
			log.Fatalf("timeout reached, force closing server: %v", err)
		} else {
			log.Fatalf("error during shutdown: %v", err)
		}
	} else {
		log.Println("server stopped")
	}
	log.Println("server.Shutdown() returned")
	time.Sleep(time.Second * 5)
	log.Println("waiting for goroutines to stop")
	wg.Wait()
	log.Println("goroutines done")
}

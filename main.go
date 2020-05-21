package main

import (
	"github.com/doubtnut/handler"
	"github.com/doubtnut/logging"
	"github.com/doubtnut/redis"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

var Logger = logging.NewLogger()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	} else {
		log.Println(".env file loaded")
	}
}

func main() {
	Logger.Infof("starting service ")
	redis.Init()
	//scheduler.Init()
	router := mux.NewRouter()
	sub := router.PathPrefix("/api/v1").Subrouter()
	sub.Methods(http.MethodPost).Path("/send/{id:[0-9]+}").HandlerFunc(handler.Send)

	if err := http.ListenAndServe(":6000", router); err != nil {
		log.Fatal("exiting")
	}

}

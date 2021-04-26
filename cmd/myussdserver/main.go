package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	myussd "github.com/yeboahnanaosei/x/testussd"
	redislocal "github.com/yeboahnanaosei/x/testussd/storage/redis"
)

var userService myussd.UserService
var caseService myussd.CaseService

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading environment: %v", err)
	}

	options, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("could not parse redis url: %v", err)
	}

	redisClient := redis.NewClient(options)
	userService = &redislocal.UserService{DB: redisClient}
	caseService = &redislocal.CaseService{DB: redisClient}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.HandleFunc("/", handleIndex)
	fmt.Println("listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		body := []byte("Method Not Allowed")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.Write(body)
		return
	}

	activeUser := myussd.User{
		ProfileName: strings.TrimSpace(r.FormValue("ProfileName")),
		WaId:        strings.TrimSpace(r.FormValue("WaId")),
		Input:       strings.TrimSpace(r.FormValue("Body")),

		// Africa's Talking API
		// ProfileName: strings.TrimSpace(r.FormValue("sessionId")),
		// WaId:        strings.TrimSpace(r.FormValue("phoneNumber")),
		// Input:       strings.TrimSpace(r.FormValue("text")),
	}

	if activeUser.ProfileName == "" || activeUser.WaId == "" {
		w.WriteHeader(http.StatusBadRequest)
		body := []byte("END Invalid request. No 'ProfileName' or 'WaId'")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.Write(body)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*8)
	defer cancel()

	stage, err := userService.GetStage(ctx, activeUser)
	if !errors.Is(err, myussd.ErrStageNotFound) && err != nil {
		fmt.Println("error getting user stage: ", err) // TODO: Log error
		w.WriteHeader(http.StatusInternalServerError)
		body := []byte("END we suffered an internal error")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		w.Write(body)
		return
	}
	activeUser.LastStage = stage
	handleResponse(ctx, activeUser, w)
}

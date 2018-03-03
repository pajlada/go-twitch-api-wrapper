package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dankeroni/gotwitch"
	"github.com/gorilla/mux"
)

var api *gotwitch.TwitchAPI

func getUsernameByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keys, _ := vars["userid"]

	c := make(chan string)
	statusCode := http.StatusOK

	keyString := strings.Split(keys, ",")

	go api.GetUsers(keyString, func(users []gotwitch.User) {
		if len(users) == 0 {
			c <- "No user found forsenT"
			statusCode = http.StatusNotFound
			return
		}

		result := make([]string, len(users))
		for i := 0; i < len(users); i++ {
			result[i] = users[i].Login
		}

		c <- strings.Join(result, ",")
	}, func(code int, statusMessage, errorMessage string) {
		c <- errorMessage
		statusCode = code
	}, func(err error) {
		c <- err.Error()
		statusCode = http.StatusInternalServerError
	})

	time.AfterFunc(5*time.Second, func() {
		// Deadline timer
		c <- "Deadline timer met"
		statusCode = http.StatusInternalServerError
	})

	response := <-c

	w.WriteHeader(statusCode)
	w.Write([]byte(response))
}

func getIDByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keys, _ := vars["username"]

	c := make(chan string)
	statusCode := http.StatusOK

	keyString := strings.Split(keys, ",")

	go api.GetUsersByLogin(keyString, func(users []gotwitch.User) {
		if len(users) == 0 {
			c <- "No user found forsenT"
			statusCode = http.StatusNotFound
			return
		}

		results := make([]string, len(users))
		for i := 0; i < len(users); i++ {
			results[i] = users[i].ID
		}

		c <- strings.Join(results, ",")
	}, func(code int, statusMessage, errorMessage string) {
		c <- errorMessage
		statusCode = code
	}, func(err error) {
		c <- err.Error()
		statusCode = http.StatusInternalServerError
	})

	time.AfterFunc(5*time.Second, func() {
		// Deadline timer
		c <- "Deadline timer met"
		statusCode = http.StatusInternalServerError
	})

	response := <-c

	w.WriteHeader(statusCode)
	w.Write([]byte(response))
}

func about(w http.ResponseWriter, r *http.Request) {
	str := "Contact me on twitter for concerns twitter.com/pajlada or through twitch whispers twitch.tv/pajlada"
	w.Write([]byte(str))
}

func main() {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	api = gotwitch.New(clientID)

	router := mux.NewRouter()

	router.HandleFunc("/api/twitch/getusernamebyid/{userid}", getUsernameByID).Methods("GET")
	router.HandleFunc("/api/twitch/getidbyusername/{username}", getIDByUsername).Methods("GET")
	router.HandleFunc("/api/twitch/", about)

	fmt.Println("online using client id", clientID)

	log.Fatal(http.ListenAndServe(":8080", router))
}

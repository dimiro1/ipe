package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pusher/pusher-http-go"
)

var client pusher.Client

func init() {
	client = pusher.Client{
		AppId:  "1",
		Key:    "278d525bdf162c739803",
		Secret: "7ad3753142a6693b25b9",
		Host:   ":8080",
	}
}

func pusherPresenceAuth(res http.ResponseWriter, req *http.Request) {
	log.Println("Presence Request")
	presenceData := pusher.MemberData{
		UserId:   "1",
		UserInfo: map[string]string{},
	}

	params, _ := ioutil.ReadAll(req.Body)
	response, err := client.AuthenticatePresenceChannel(params, presenceData)

	if err != nil {
		panic(err)
	}

	fmt.Fprint(res, string(response))
}

func pusherPrivateAuth(res http.ResponseWriter, req *http.Request) {
	params, _ := ioutil.ReadAll(req.Body)
	response, err := client.AuthenticatePrivateChannel(params)

	log.Printf("Private Request %s", params)
	log.Printf("Auth %s", response)

	if err != nil {
		panic(err)
	}

	fmt.Fprint(res, string(response))
}

func triggerMessage(res http.ResponseWriter, _ *http.Request) {
	client.Trigger("private-messages", "messages", "The message from server")

	fmt.Fprint(res, "OK")
}

func main() {
	http.HandleFunc("/pusher/presence/auth", pusherPresenceAuth)
	http.HandleFunc("/pusher/private/auth", pusherPrivateAuth)
	http.HandleFunc("/trigger", triggerMessage)
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":5000", nil)
}

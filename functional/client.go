package main

import (
	"fmt"
	"io/ioutil"
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
	presenceData := pusher.MemberData{
		UserId:   "1",
		UserInfo: map[string]string{},
	}

	params, _ := ioutil.ReadAll(req.Body)
	response, err := client.AuthenticatePresenceChannel(params, presenceData)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(res, string(response))
}

func pusherPrivateAuth(res http.ResponseWriter, req *http.Request) {
	params, _ := ioutil.ReadAll(req.Body)
	response, err := client.AuthenticatePrivateChannel(params)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(res, string(response))
}

func triggerMessage(res http.ResponseWriter, req *http.Request) {
	client.Trigger("private-messages", "messages", "The message from server")

	fmt.Fprintf(res, "OK")
}

func main() {
	http.HandleFunc("/pusher/presence/auth", pusherPresenceAuth)
	http.HandleFunc("/pusher/private/auth", pusherPrivateAuth)
	http.HandleFunc("/trigger", triggerMessage)
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":5000", nil)
}

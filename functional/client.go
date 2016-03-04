package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pusher/pusher-http-go"
)

var (
	client1 pusher.Client
	client2 pusher.Client
)

func init() {
	client1 = pusher.Client{
		AppId:  "1",
		Key:    "278d525bdf162c739803",
		Secret: "7ad3753142a6693b25b9",
		Host:   ":8080",
	}
	client2 = pusher.Client{
		AppId:  "2",
		Secure: true,
		Key:    "c8b30f611ffb13202976",
		Secret: "d6824d2fa32888931504",
		Host:   ":8090",
	}
}

func pusherPresenceAuth(res http.ResponseWriter, req *http.Request) {
	presenceData := pusher.MemberData{
		UserId:   "1",
		UserInfo: map[string]string{},
	}

	params, _ := ioutil.ReadAll(req.Body)
	response, err := client1.AuthenticatePresenceChannel(params, presenceData)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(res, string(response))
}

func pusherPrivateAuth(res http.ResponseWriter, req *http.Request) {
	params, _ := ioutil.ReadAll(req.Body)
	response, err := client1.AuthenticatePrivateChannel(params)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(res, string(response))
}

func triggerMessage(res http.ResponseWriter, req *http.Request) {
	client1.Trigger("private-messages", "messages", "The message from server")

	fmt.Fprintf(res, "OK")
}

func main() {
	http.HandleFunc("/pusher/presence/auth", pusherPresenceAuth)
	http.HandleFunc("/pusher/private/auth", pusherPrivateAuth)
	http.HandleFunc("/trigger", triggerMessage)

	http.Handle("/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":5000", nil)
}

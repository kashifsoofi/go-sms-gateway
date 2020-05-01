package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type sendSmsRequest struct {
	Recipient string `json:"recipient"`
	Body      string `json:"body"`
	Sender    string `json:"sender"`
}

var recipientAndCodeMap map[string]string

type server struct{}

func (s *server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		rw.WriteHeader(http.StatusOK)
		recipient := strings.TrimPrefix(r.URL.Path, "/")
		code := recipientAndCodeMap[recipient]
		response := fmt.Sprint("{\"code\": \"", code, "\"}")
		rw.Write([]byte(response))
	case "POST":
		var sendSms sendSmsRequest
		if err := json.NewDecoder(r.Body).Decode(&sendSms); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		code := strings.Trim(sendSms.Body, "Your verification code is: ")
		recipientAndCodeMap[sendSms.Recipient] = code
		rw.WriteHeader(http.StatusOK)
	default:
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{"message": "method not supported"}`))
	}
}

func main() {
	s := &server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

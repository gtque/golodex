package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/idtoken"
	"internal/golodex"
	"io/ioutil"
	"net/http"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	// Declare a new Person struct.
	var a golodex.Auth

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(req.Body).Decode(&a)
	if err != nil {
		fmt.Println("problem with json: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload, err := idtoken.Validate(context.Background(), a.Token, golodex.Property("GOOGLE_CLIENT_ID","710486342711-3g14joh4shppp62in94pbpktfpqhpeoc.apps.googleusercontent.com"))
	if err != nil {
		fmt.Println("problem checking token: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(a)
	//err := json.Unmarshal(payload.Claims, &g)
	gmail := payload.Claims["email"].(string)
	if a.Email == gmail {
		fmt.Println("emails matched, authenticated in.")
		//login to golodex
		dResponse,err := http.Post(golodex.Property("GOLODEX_DATA_HOST", "http://localhost:8095")+"/api/login", "application/json", bytes.NewBuffer(b))
		if err != nil {
			fmt.Println("problem checking token: " + err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer dResponse.Body.Close()
		body, err := ioutil.ReadAll(dResponse.Body)
		if err != nil {
			fmt.Println("problem checking token: " + err.Error())
			http.Error(w, "invalid request", http.StatusInternalServerError)
			return
		}
		fmt.Println("got a response from the data server.")
		fmt.Println(string(body))
		fmt.Fprintf(w, string(body))
	} else {
		//fmt.Println("problem checking token: " + err.Error())
		http.Error(w, "invalid login", http.StatusBadRequest)
		return
	}
	//fmt.Println(string(b))
	//fmt.Println(string(pb))
	//json.NewEncoder().Encode()
	//fmt.Fprintf(w, "{\"response\":\"greetings\"}")
}

func (p pager) Methods() map[string]bool {
	return map[string]bool{"POST":true}
}

func (p pager) Authenticate(req *http.Request) bool {
	return true
}

var Pager pager
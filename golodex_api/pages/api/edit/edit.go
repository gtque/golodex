package main

import (
	"internal/golodex"
	"io/ioutil"
	"log"
	"net/http"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	log.Println("editing...")
	var cardId = req.URL.Query().Get("card")
	//log.Println(cardId)
	//.String()//.Query().Get("card")
	//log.Println("url: " + cardId)
	//var cardId = ""
	//defer req.Body.Close()
	bodyCU, errCU := ioutil.ReadAll(req.Body)
	if errCU != nil {
		log.Println("error reading body trying to edit a card")
		log.Println(errCU)
		golodex.WriteResponse(w, "{\"error\": \"adding new card\"}", http.StatusInternalServerError)
		return
	}
	//log.Println("trying to edit: " + string(bodyCU))
	var bodyS = string(bodyCU)
	if bodyS == "" {
		log.Println("no body sent, nothing to do.")
		golodex.WriteResponse(w, "{\"error\": \"no date sent\"}", http.StatusBadRequest)
		return
	} else {
		log.Println("got body?")
	}

	if cardId == "" || cardId == "_empty" {
		log.Println("trying to add: " + bodyS)
		dResponse, err := golodex.ApiQueryPost(req, "add", string(bodyCU))
		if err != nil {
			log.Println("failed to add a new card")
			golodex.WriteResponse(w, "{\"error\": \"adding new card\"}", http.StatusInternalServerError)
		} else {
			golodex.WriteResponse(w, string(dResponse), http.StatusOK)
		}
	} else {
		log.Println("trying to edit: " + cardId)
		dResponse, err := golodex.ApiQueryPut(req, "edit?id="+cardId, string(bodyCU))
		if err != nil {
			log.Println("failed to get the edit the card...")
			golodex.WriteResponse(w, "{\"error\": \"getting all rolodex entries failed\"}", http.StatusInternalServerError)
		} else {
			//log.Println("got the rolodex: " + string(dResponse))
			log.Println("edited: " + string(dResponse))
			golodex.WriteResponse(w, string(dResponse), http.StatusOK)
		}
	}
	return
}

func (p pager) Methods() map[string]bool {
	return map[string]bool{"PUT": true}
}

func (p pager) Authenticate(req *http.Request) bool {
	log.Println("authenticating for edit...")
	return golodex.CheckAuthorization(req)
}

var Pager pager

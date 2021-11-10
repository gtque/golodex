package main

import (
	"internal/golodex"
	"log"
	"net/http"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	var cardId = req.URL.Query().Get("card")
	dResponse, err := golodex.ApiQueryGet(req, "card?id="+cardId)
	if err != nil {
		log.Println("failed to get the card: " + cardId)
		golodex.WriteResponse(w, "{\"error\": \"getting all rolodex entries failed\"}", http.StatusInternalServerError)
		return
	}
	//log.Println("got the rolodex: " + string(dResponse))
	golodex.WriteResponse(w, string(dResponse), http.StatusOK)
	return
}

func (p pager) Methods() map[string]bool {
	return map[string]bool{"GET": true}
}

func (p pager) Authenticate(req *http.Request) bool {
	return golodex.CheckAuthorization(req)
}

var Pager pager

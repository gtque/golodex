package main

import (
	"encoding/json"
	"internal/golodexdata"
	"log"
	"net/http"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	log.Println("deleting a card...")
	var cardId = req.URL.Query().Get("id")
	updatedCard,errCard := golodexdata.CouchQueryGet(req, cardId)
	if errCard != nil {
		golodexdata.WriteResponse(w, "{\"error\": \"card not found\"}", http.StatusInternalServerError)
		return
	}
	var card golodexdata.GolodexCard
	errUm := json.Unmarshal(updatedCard,&card)
	if errUm != nil {
		log.Println("problem unmarshalling to card")
		log.Println(errUm)
		golodexdata.WriteResponse(w, "{\"error\": \"parsing request\"}", http.StatusInternalServerError)
		return
	}
	deleted,errAll := golodexdata.CouchQueryDelete(req, card.Id+"?rev="+card.Revision)
	if errAll != nil {
		golodexdata.WriteResponse(w, "{\"error\": \"reading from couch\"}", http.StatusInternalServerError)
		return
	}
	golodexdata.WriteResponse(w, string(deleted), http.StatusOK)
}

func (p pager) Methods() map[string]bool {
	return map[string]bool{"DELETE":true}
}

func (p pager) Authenticate(req *http.Request) bool {
	return golodexdata.CheckAuthorization(req)
}

var DataPager pager
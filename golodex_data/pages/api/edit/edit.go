package main

import (
	"encoding/json"
	"internal/golodexdata"
	"log"
	"net/http"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	log.Println("editing a card...")
	var card golodexdata.GolodexCard
	errUm := json.NewDecoder(req.Body).Decode(&card)
	if errUm != nil {
		log.Println("problem unmarshalling to card")
		log.Println(errUm)
		golodexdata.WriteResponse(w, "{\"error\": \"parsing request\"}", http.StatusInternalServerError)
		return
	}
	body,errM := json.Marshal(card)
	if errM != nil {
		log.Println("problem marshalling card")
		log.Println(errM)
		golodexdata.WriteResponse(w, "{\"error\": \"problem generating body for update\"}", http.StatusInternalServerError)
		return
	}
	_,errAll := golodexdata.CouchQueryPut(req, card.Id, string(body))
	if errAll != nil {
		golodexdata.WriteResponse(w, "{\"error\": \"reading from couch\"}", http.StatusInternalServerError)
		return
	}
	updatedCard,errCard := golodexdata.CouchQueryGet(req, card.Id)
	if errCard != nil {
		golodexdata.WriteResponse(w, "{\"error\": \"getting updated card from couch\"}", http.StatusInternalServerError)
		return
	}
	golodexdata.WriteResponse(w, string(updatedCard), http.StatusOK)
}

func (p pager) Methods() map[string]bool {
	return map[string]bool{"PUT":true}
}

func (p pager) Authenticate(req *http.Request) bool {
	return golodexdata.CheckAuthorization(req)
}

var DataPager pager
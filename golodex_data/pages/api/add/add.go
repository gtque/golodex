package main

import (
	"encoding/json"
	"internal/golodexdata"
	"log"
	"net/http"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	log.Println("adding new card...")
	var cardNew golodexdata.GolodexCard
	var card golodexdata.GolodexNewCard
	errUm := json.NewDecoder(req.Body).Decode(&cardNew)
	if errUm != nil {
		log.Println("problem unmarshalling to card")
		log.Println(errUm)
		golodexdata.WriteResponse(w, "{\"error\": \"parsing request\"}", http.StatusInternalServerError)
		return
	}
	card.Name = cardNew.Name
	card.PhoneNumbers = cardNew.PhoneNumbers
	card.Addresses = cardNew.Addresses
	log.Println("new card built, ")
	body,errM := json.Marshal(card)
	if errM != nil {
		log.Println("problem marshalling card")
		log.Println(errM)
		golodexdata.WriteResponse(w, "{\"error\": \"problem generating body for update\"}", http.StatusInternalServerError)
		return
	}
	log.Println("new card: " + string(body))
	_new,errAll := golodexdata.CouchQueryPost(req, "", string(body))
	if errAll != nil {
		golodexdata.WriteResponse(w, "{\"error\": \"reading from couch\"}", http.StatusInternalServerError)
		return
	}
	log.Println("status of create: " + string(_new))
	var createdCard golodexdata.GolodexCreate
	errNew := json.Unmarshal(_new, &createdCard)
	if errNew != nil {
		log.Println("error creating new entry")
		log.Println(errNew)
		golodexdata.WriteResponse(w, "{\"error\": \"creating new entry\"}", http.StatusInternalServerError)
		return
	}
	log.Println("new card id: " + createdCard.Id)
	updatedCard,errCard := golodexdata.CouchQueryGet(req, createdCard.Id)
	if errCard != nil {
		golodexdata.WriteResponse(w, "{\"error\": \"getting newly created card from couch\"}", http.StatusInternalServerError)
		return
	}
	log.Println("created card: " + string(updatedCard))
	golodexdata.WriteResponse(w, string(updatedCard), http.StatusOK)
}

func (p pager) Methods() map[string]bool {
	return map[string]bool{"POST":true}
}

func (p pager) Authenticate(req *http.Request) bool {
	return golodexdata.CheckAuthorization(req)
}

var DataPager pager
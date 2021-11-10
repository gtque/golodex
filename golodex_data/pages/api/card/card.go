package main

import (
	"encoding/json"
	"internal/golodexdata"
	"log"
	"net/http"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	var card golodexdata.GolodexCard
	var updatedCard []byte
	var cardId = req.URL.Query().Get("id")
	if cardId == "_empty" {
		var phone golodexdata.GolodexPhoneNumber
		var address golodexdata.GolodexAddress
		card.PhoneNumbers = []golodexdata.GolodexPhoneNumber{phone}
		card.Addresses = []golodexdata.GolodexAddress{address}
		uc, errM := json.Marshal(card)
		if errM != nil {
			log.Println("problem marshalling to empty card")
			log.Println(errM)
			golodexdata.WriteResponse(w, "{\"error\": \"problem with empty card\"}", http.StatusInternalServerError)
			return
		}
		updatedCard = uc
		log.Println("returning the empty card: " + string(updatedCard))
	} else {
		uc, errCard := golodexdata.CouchQueryGet(req, cardId)
		if errCard != nil {
			log.Println(errCard)
			golodexdata.WriteResponse(w, "{\"error\": \"getting card from couch\"}", http.StatusInternalServerError)
		}
		updatedCard = uc
	}
	golodexdata.WriteResponse(w, string(updatedCard), http.StatusOK)
}

func (p pager) Methods() map[string]bool {
	return map[string]bool{"GET":true}
}

func (p pager) Authenticate(req *http.Request) bool {
	return golodexdata.CheckAuthorization(req)
}

var DataPager pager
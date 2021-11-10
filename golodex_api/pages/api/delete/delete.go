package main

import (
	"internal/golodex"
	"log"
	"net/http"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	log.Println("editing...")
	var cardId = req.URL.Query().Get("card")

	if cardId == "" || cardId == "_empty" {
		log.Println("nothing to delete, so I am pretending I did.")
		golodex.WriteResponse(w, "{}", http.StatusOK)
		return
	} else {
		log.Println("trying to delete: " + cardId)
		dResponse, err := golodex.ApiQueryDelete(req, "delete?id="+cardId)
		if err != nil {
			log.Println("failed to get the edit the card...")
			golodex.WriteResponse(w, "{\"error\": \"getting deleting card\"}", http.StatusInternalServerError)
		} else {
			//log.Println("got the rolodex: " + string(dResponse))
			log.Println("deleted: " + string(dResponse))
			golodex.WriteResponse(w, string(dResponse), http.StatusOK)
		}
	}
	return
}

func (p pager) Methods() map[string]bool {
	return map[string]bool{"DELETE": true}
}

func (p pager) Authenticate(req *http.Request) bool {
	log.Println("authenticating for edit...")
	return golodex.CheckAuthorization(req)
}

var Pager pager

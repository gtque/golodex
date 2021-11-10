package main

import (
	"internal/golodex"
	"log"
	"net/http"
	"net/url"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	var sort = req.URL.Query().Get("sort")
	var search = req.URL.Query().Get("search")
	var params = url.Values{}
	params.Add("sort", sort)
	params.Add("search", search)
	var all = "all?"+params.Encode()
	log.Println("all: " + all)
	dResponse, err := golodex.ApiQueryGet(req, all)
	if err != nil {
		log.Println("failed to get the rolodex")
		golodex.WriteResponse(w, "{\"error\": \"getting all rolodex entries failed\"}", http.StatusInternalServerError)
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

package main

import (
	"net/http"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	f := "./pages/api/about/about.json"
	http.ServeFile(w, req, f)
}

func (p pager) Methods() map[string]bool {
	return map[string]bool{"GET":true}
}

func (p pager) Authenticate(req *http.Request) bool {
	return true
}

var DataPager pager
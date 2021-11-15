package main

import (
	"net/http"
	"testing"
)

func simple_get_all(req *http.Request, query string)([]byte,error){

	data := "{"+
		"}"
	return []byte(data),nil
}
func TestSimpleAll(t *testing.T) {
	var w http.ResponseWriter
	var req *http.Request
	ThePage(w, req, simple_get_all)
}

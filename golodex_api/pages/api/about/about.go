package main

import (
	"fmt"
	"net/http"
	"strings"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	//var path = req.URL.Path
	//var param = req.URL.Query().Get("var1")
	//fmt.Fprintf(w, "greetings: " + path +":" + param +"\n")
	fmt.Println(req.URL.Path)
	f := req.URL.Path
	if f == "/" {
		f = "./static/index.html"
	} else {
		fn := strings.Replace(f, "/api", "", 1)
		f = "./pages"+f+fn+".html"
	}
	http.ServeFile(w, req, f)
}

func (p pager) Methods() map[string]bool {
	return map[string]bool{"GET":true}
}

func (p pager) Authenticate(req *http.Request) bool {
	return true
}

var Pager pager
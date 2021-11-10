package golodex

import (
	"log"
	"net/http"
)

type Session struct {
	Id string
	Pw string
}

func CheckAuthorization(req *http.Request)bool {
	log.Println("testing authorization...")
	return true
}

func SessionInformation(req *http.Request)Session {
	var session Session
	//session.Auth = base64.StdEncoding.EncodeToString([]byte(req.Header.Get("X-GOLODEX-ID")+":"+GetSessionPassword(req.Header.Get("X-GOLODEX-TOKEN"))))
	session.Id = req.Header.Get("X-GOLODEX-ID")
	session.Pw = req.Header.Get("X-GOLODEX-TOKEN")
	return session
}

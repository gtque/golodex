package golodexdata

import (
	"encoding/hex"
	"log"
	"net/http"
)

type Session struct {
	DB string
	Id string
	Pw string
}

func CheckAuthorization(req *http.Request)bool {
	log.Println("testing authorization...")
	return true
}

func SessionInformation(req *http.Request)Session {
	var session Session
	session.DB = "userdb-" + hex.EncodeToString([]byte(req.Header.Get("X-GOLODEX-ID")))
	//session.Auth = base64.StdEncoding.EncodeToString([]byte(req.Header.Get("X-GOLODEX-ID")+":"+GetSessionPassword(req.Header.Get("X-GOLODEX-TOKEN"))))
	session.Id = req.Header.Get("X-GOLODEX-ID")
	session.Pw = GetSessionPassword(req.Header.Get("X-GOLODEX-TOKEN"))
	return session
}

func GetSessionPassword(token string)string {
	return token
}
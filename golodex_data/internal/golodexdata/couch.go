package golodexdata

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CouchError struct {
	When time.Time
	What string
}

func (c CouchError) Error() string {
	return fmt.Sprintf("%v: %v", c.When, c.What)
}

type GolodexCouchRows struct {
	Rows []GolodexCouchRow `json:"rows"`
}

type GolodexCouchViewRows struct {
	Rows []GolodexCouchViewRow `json:"rows"`
}

type GolodexCouchRow struct {
	Doc GolodexCard `json:"doc"`
}

type GolodexCouchViewRow struct {
	Doc GolodexCard `json:"value"`
}

type GolodexRolodex struct {
	Cards []GolodexCard `json:"cards"`
}

type GolodexCard struct {
	Id string `json:"_id"`
	Revision string `json:"_rev"`
	Name GolodexName `json:"name"`
	PhoneNumbers []GolodexPhoneNumber `json:"phone_numbers"`
	Addresses []GolodexAddress `json:"addresses"`
}

type GolodexNewCard struct {
	Name GolodexName `json:"name"`
	PhoneNumbers []GolodexPhoneNumber `json:"phone_numbers"`
	Addresses []GolodexAddress `json:"addresses"`
}

type GolodexName struct {
	First string `json:"first"`
	Middle string `json:"middle"`
	Last string `json:"last"`
}

type GolodexPhoneNumber struct {
	Number string `json:"number"`
	Type string `json:"type"`
}

type GolodexAddress struct {
	Street1 string `json:"street_1"`
	Street2 string `json:"street_2"`
	City string `json:"city"`
	State string `json:"state"`
	ZipCode string `json:"zip_code"`
	Type string `json:"type"`
}

type GolodexCreate struct {
	Id string `json:"id"`
	Revision string `json:"rev"`
	Ok bool `json:"ok"`
}

func couchQuery(req *http.Request, query string, data []byte, method string)([]byte,error) {
	session := SessionInformation(req)
	//setup client to make the call
	client := http.Client{Timeout: 5 * time.Second}
	log.Println("building new couch request...")

	var couchQ *http.Request
	var errNR error
	if query != "" {
		query = "/" + query
	}
	if data == nil {
		couchQ, errNR = http.NewRequest(method, Property("COUCHDB_HOST", "http://localhost:5984")+"/"+session.DB+query, nil)
	} else {
		couchQ, errNR = http.NewRequest(method, Property("COUCHDB_HOST", "http://localhost:5984")+"/"+session.DB+query, bytes.NewBuffer(data))
	}
	if errNR != nil {
		log.Println("invalid session.")
		log.Println(errNR)
		return nil,errNR
	}
	couchQ.SetBasicAuth(session.Id, session.Pw)
	couchQ.Header.Set("Content-Type", "application/json")
	couchResponse, errR := client.Do(couchQ)
	if errR != nil {
		//problem checking user.
		log.Println("error calling couch")
		log.Println(errR)
		return nil,errR
	}
	log.Printf("status from couch: %d\n", couchResponse.StatusCode)
	defer couchResponse.Body.Close()
	if couchResponse.StatusCode >= 300 {
		return nil, CouchError{
			time.Now(),
			"problem talking to couch",
		}
	}
	bodyCU, errCU := ioutil.ReadAll(couchResponse.Body)
	if errCU != nil {
		log.Println("error reading response from couch")
		log.Println(errR)
		return nil,errR
	}
	return bodyCU,nil
}

type CouchQueryNoBody func(req *http.Request, query string)([]byte,error)

func CouchQueryGet(req *http.Request, query string)([]byte,error) {
	bodyCU, err := couchQuery(req, query, nil, http.MethodGet)
	//log.Println("response from couch: " + (string(bodyCU)))
	return bodyCU,err
}

func CouchQueryDelete(req *http.Request, query string)([]byte,error) {
	bodyCU, err := couchQuery(req, query, nil, http.MethodDelete)
	//log.Println("response from couch: " + (string(bodyCU)))
	return bodyCU,err
}

type CouchQueryWithBody func(req *http.Request, query string, data string)([]byte,error)

func CouchQueryPut(req *http.Request, query string, data string)([]byte,error) {
	bodyCU, err := couchQuery(req, query, []byte(data), http.MethodPut)
	//log.Println("response from couch: " + (string(bodyCU)))
	return bodyCU,err
}

func CouchQueryPost(req *http.Request, query string, data string)([]byte,error) {
	bodyCU, err := couchQuery(req, query, []byte(data), http.MethodPost)
	//log.Println("response from couch: " + (string(bodyCU)))
	return bodyCU,err
}
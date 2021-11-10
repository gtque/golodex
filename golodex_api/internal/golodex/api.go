package golodex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ApiError struct {
	When time.Time
	What string
}

func (c ApiError) Error() string {
	return fmt.Sprintf("%v: %v", c.When, c.What)
}

func apiQuery(req *http.Request, query string, data []byte, method string)([]byte,error) {
	session := SessionInformation(req)
	//setup client to make the call
	client := http.Client{Timeout: 5 * time.Second}
	log.Println("building new api request...")

	var apiQ *http.Request
	var errNR error
	if data == nil {
		apiQ, errNR = http.NewRequest(method, Property("GOLODEX_DATA_HOST", "http://localhost:8095")+"/api/"+query, nil)
	} else {
		apiQ, errNR = http.NewRequest(method, Property("GOLODEX_DATA_HOST", "http://localhost:8095")+"/api/"+query, bytes.NewBuffer(data))
	}
	if errNR != nil {
		log.Println("invalid session.")
		log.Println(errNR)
		return nil,errNR
	}
	apiQ.Header.Set("X-GOLODEX-ID", session.Id)
	apiQ.Header.Set("X-GOLODEX-TOKEN", session.Pw)
	apiResponse, errR := client.Do(apiQ)
	if errR != nil {
		//problem checking user.
		log.Println("error calling api")
		log.Println(errR)
		return nil,errR
	}
	log.Printf("status from api: %d\n", apiResponse.StatusCode)
	defer apiResponse.Body.Close()
	if apiResponse.StatusCode >= 300 {
		return nil, ApiError{
			time.Now(),
			"problem talking to api",
		}
	}
	bodyCU, errCU := ioutil.ReadAll(apiResponse.Body)
	if errCU != nil {
		log.Println("error reading response from api")
		log.Println(errR)
		return nil,errR
	}
	return bodyCU,nil
}

func ApiQueryGet(req *http.Request, query string)([]byte,error) {
	bodyCU, err := apiQuery(req, query, nil, http.MethodGet)
	//log.Println("response from api: " + (string(bodyCU)))
	return bodyCU,err
}

func ApiQueryDelete(req *http.Request, query string)([]byte,error) {
	bodyCU, err := apiQuery(req, query, nil, http.MethodDelete)
	//log.Println("response from api: " + (string(bodyCU)))
	return bodyCU,err
}

func ApiQueryPut(req *http.Request, query string, data string)([]byte,error) {
	bodyCU, err := apiQuery(req, query, []byte(data), http.MethodPut)
	//log.Println("response from api: " + (string(bodyCU)))
	return bodyCU,err
}

func ApiQueryPost(req *http.Request, query string, data string)([]byte,error) {
	bodyCU, err := apiQuery(req, query, []byte(data), http.MethodPost)
	//log.Println("response from api: " + (string(bodyCU)))
	return bodyCU,err
}
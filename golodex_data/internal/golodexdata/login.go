package golodexdata

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Auth struct {
	Name string
	Email  string
	Token string
	Id string
}

type UserCreds struct {
	Name string `json:"name"`
	Password string `json:"password"`
}

type UserCreate struct {
	Name string `json:"name"`
	Password string `json:"password"`
	TheType string `json:"type"`
	Roles []string `json:"roles"`
}

type UserCreateCheck struct {
	Ok bool `json:"ok"`
	Name string `json:"id"`
	Rev string `json:"rev"`
}

type UserCheck struct {
	Ok bool `json:"ok"`
	Name string `json:"name"`
	Roles []string `json:"roles"`
}

func waitForDb(userCreate UserCreate, udb string, count int) (string, error) {
	if count == 0 {
		return "failed", CouchError{
			time.Now(),
			"the db never existed...",
		}
	}
	client := http.Client{Timeout: 5 * time.Second}
	couchDbExists, errUVR := http.NewRequest(http.MethodGet, Property("COUCHDB_HOST", "http://localhost:5984") + udb +"/_all_docs?include_docs=true", nil)
	couchDbExists.SetBasicAuth(userCreate.Name, userCreate.Password)
	//couchUserViewReg.SetBasicAuth(Property("COUCHDB_ADMIN_USER","admin"), Secret("COUCHDB_ADMIN_PW"))
	if errUVR != nil {
		log.Println("failed to build request to check for the DB.")
		log.Println(errUVR)
		return "failed", CouchError{
			time.Now(),
			"problem checking for DB",
		}
	}
	//execute request to create the user
	couchDbExists.Header.Add("Content-Type", "application/json")
	userView, errUCVR := client.Do(couchDbExists)
	if errUCVR != nil {
		log.Println("failed to check db, executing request failed.")
		log.Println(errUCVR)
		return "failed", CouchError{
			time.Now(),
			"problem checking for DB",
		}
	}
	defer userView.Body.Close()
	if userView.StatusCode != 200 {
		bodyV, errVR := ioutil.ReadAll(userView.Body)
		if errVR != nil {
			log.Println("problem creating view: " + errVR.Error())
			return "failed", CouchError{
				time.Now(),
				"problem checking for DB",
			}
		}
		log.Println("got a response from couch for the view but the DB doesn't exist...")
		log.Println(string(bodyV))
		time.Sleep(2 * time.Second)
		return waitForDb(userCreate, udb, count-1)
	}
	log.Println("db exists!!!")
	return "success", nil
}

func Login(w http.ResponseWriter, req *http.Request) {
	log.Println("---------------LOGGING IN---------------")
	// Declare a new Person struct.
	var a Auth
	if req.Method != "POST" {
		http.Error(w, "{\"error\":\"unable to log in.\"}", http.StatusUnauthorized)
		return
	}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(req.Body).Decode(&a)
	if err != nil {
		log.Println("problem with json: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//b, err := json.Marshal(a)
	//log.Println(string(b))

	//setup client to make the call
	client := http.Client{Timeout: 5 * time.Second}

	//check if user exists
	var user UserCreds
	var userCheck UserCheck
	//defining user credentials
	user.Name = a.Id
	user.Password = strings.Split(a.Email, "@")[0] + Secret("COUCHDB_USER_PW")
	//defining body for user check.
	userBody,errJ := json.Marshal(user)
	if errJ != nil {
		//problem checking user.
		log.Println("error marshalling user...")
		log.Println(errJ)
		http.Error(w, "{\"error\": \"failed to check user, try again.\"}", http.StatusInternalServerError)
		return
	}
	//fmt.Printf("%#v\n", user)
	//log.Println(string(userBody))
	//user exists check.
	couchUserReg,errNR := http.NewRequest(http.MethodPost, Property("COUCHDB_HOST","http://localhost:5984")+"/_session", bytes.NewBuffer(userBody))
	if errNR != nil {
		//problem checking user.
		log.Println("error building check user request...")
		log.Println(errNR)
		http.Error(w, "{\"error\": \"failed to check user, try again.\"}", http.StatusInternalServerError)
		return
	}
	couchUserReg.Header.Add("Content-Type", "application/json")
	//set basic auth
	//couchUserReg.SetBasicAuth(properties.Property("COUCHDB_ADMIN_USER","admin"), properties.Secret("COUCHDB_ADMIN_PW"))
	//call user check endpoint
	couchUser,errR := client.Do(couchUserReg)
	log.Println("check user status code: " + strconv.Itoa(couchUser.StatusCode))
	if errR != nil {
		//problem checking user.
		log.Println("error checking if user already exists...")
		log.Println(errR)
		http.Error(w, "{\"error\": \"failed to check user, try again.\"}", http.StatusInternalServerError)
		return
	} else {
		defer couchUser.Body.Close()
		bodyCU,errCU := ioutil.ReadAll(couchUser.Body)
		if errCU != nil {
			log.Println("problem reading body from check user call: " + errCU.Error())
			http.Error(w, "{\"error\": \"invalid request\"}", http.StatusInternalServerError)
			return
		}
		log.Println("user response from couch")
		log.Println(string(bodyCU))
		//err := json.NewDecoder(req.Body).Decode(&a)
		//decode user check
		errUC := json.Unmarshal(bodyCU,&userCheck)
		log.Println(userCheck.Ok)
		log.Println(userCheck.Name)
		if errUC != nil || !userCheck.Ok {
			log.Println("have to create new user: " + user.Name)
			/**
			curl -X PUT http://localhost:5984/_users/org.couchdb.user:jan \
			     -H "Accept: application/json" \
			     -H "Content-Type: application/json" \
			     -d '{"name": "jan", "password": "apple", "roles": [], "type": "user"}'
			*/
			//assemble user create info
			var userCreate UserCreate
			userCreate.Name = user.Name
			userCreate.Password = user.Password
			userCreate.TheType = "user"
			userCreate.Roles = []string{}
			userCreateBody,errUCJ := json.Marshal(userCreate)
			log.Println(string(userCreateBody))
			if errUCJ != nil {
				log.Println("failed to json marshal user create body...")
				log.Println(err)
				http.Error(w, "{\"error\": \"invalid request: unable to define user create request\"}", http.StatusInternalServerError)
				return
			} else {
				//define request for creating a new user
				couchUserCreateReg, errUCR := http.NewRequest(http.MethodPut, Property("COUCHDB_HOST", "http://localhost:5984") + "/_users/org.couchdb.user:" + user.Name, bytes.NewBuffer(userCreateBody))
				couchUserCreateReg.SetBasicAuth(Property("COUCHDB_ADMIN_USER","admin"), Secret("COUCHDB_ADMIN_PW"))
				if errUCR != nil {
					log.Println("failed to build request to create user.")
					log.Println(err)
					http.Error(w, "{\"error\": \"failed to define user create request\"}", http.StatusInternalServerError)
					return
				}
				//execute request to create the user
				couchUserCreateReg.Header.Add("Content-Type", "application/json")
				couchUserCreated, errUCR := client.Do(couchUserCreateReg)
				if errUCR != nil {
					log.Println("failed to create user, executing request failed.")
					log.Println(errUCR)
					http.Error(w, "{\"error\": \"failed to create user\"}", http.StatusInternalServerError)
					return
				}
				defer couchUserCreated.Body.Close()
				bodyCUC,errCUC := ioutil.ReadAll(couchUserCreated.Body)
				if errCUC != nil {
					log.Println("problem creating user: " + errCUC.Error())
					http.Error(w, "{\"error\": \"invalid request\"}", http.StatusInternalServerError)
					return
				}
				log.Println("response from creating user...")
				log.Println(string(bodyCUC))
				var userCreateCheck UserCreateCheck
				errC := json.Unmarshal(bodyCUC, &userCreateCheck)
				if errC != nil {
					log.Println(errC)
					log.Println("error decoding createcheck...")
					http.Error(w, "{\"error\": \"failed to create user, decoding create check\"}", http.StatusInternalServerError)
					return
				}
				if !userCreateCheck.Ok {
					b, errUCCJ := json.Marshal(userCreateCheck)
					if errUCCJ != nil {
						log.Println(errUCCJ)
						http.Error(w, "{\"error\": \"failed to create user, unmarshalling create check\"}", http.StatusInternalServerError)
						return
					}
					log.Println(b)
					log.Println("failed to create user, or at least checking results of create user failed.")
					http.Error(w, "{\"error\": \"failed to create user, results failed\"}", http.StatusInternalServerError)
					return
				}
				b, errUCCJ := json.Marshal(userCreateCheck)
				if errUCCJ != nil {
					log.Println(errUCCJ)
					http.Error(w, "{\"error\": \"failed to create user\"}", http.StatusInternalServerError)
					return
				}
				log.Println(string(b))
				log.Println("wainting for db to exist.")
				_,errWait := waitForDb(userCreate, "/userdb-" + hex.EncodeToString([]byte(userCreate.Name)), 5)
				if errWait != nil {
					log.Println("failed to wait for the db to exist.")
					log.Println(errWait)
					http.Error(w, "{\"error\": \"user does never got created\"}", http.StatusInternalServerError)
					return
				}
				log.Println("adding the view...")
				theView := "{\n    \"views\": {\n        \"sorted_view\": {\n            \"map\": \"function(doc) { if(doc.views) {} else { emit(doc.name.last + ' ' +doc.name.first + ' ' + doc.name.middle, doc); }}\"\n        }\n    }\n}"
				//can't use CouchQueryPut here, must use session.DB = "userdb-" + hex.EncodeToString([]byte(req.Header.Get("X-GOLODEX-ID")))
				//and manually build the request...
				//_, errView := CouchQueryPut(req, "_design/golodex/_view/sorted_view", theView)
				log.Println("/userdb-" + hex.EncodeToString([]byte(userCreate.Name)))
				couchUserViewReg, errUVR := http.NewRequest(http.MethodPut, Property("COUCHDB_HOST", "http://localhost:5984") + "/userdb-" + hex.EncodeToString([]byte(userCreate.Name))+"/_design/golodex", bytes.NewBuffer([]byte(theView)))
				couchUserViewReg.SetBasicAuth(userCreate.Name, userCreate.Password)
				//couchUserViewReg.SetBasicAuth(Property("COUCHDB_ADMIN_USER","admin"), Secret("COUCHDB_ADMIN_PW"))
				if errUVR != nil {
					log.Println("failed to build request to create user view .")
					log.Println(err)
					http.Error(w, "{\"error\": \"failed to define user view request\"}", http.StatusInternalServerError)
					return
				}
				//execute request to create the user
				couchUserViewReg.Header.Add("Content-Type", "application/json")
				userView, errUCVR := client.Do(couchUserViewReg)
				if errUCVR != nil {
					log.Println("failed to create user view, executing request failed.")
					log.Println(errUCVR)
					http.Error(w, "{\"error\": \"failed to create user view\"}", http.StatusInternalServerError)
					return
				}
				defer userView.Body.Close()
				bodyV,errVR := ioutil.ReadAll(userView.Body)
				if errVR != nil {
					log.Println("problem creating view: " + errVR.Error())
					http.Error(w, "{\"error\": \"invalid request\"}", http.StatusInternalServerError)
					return
				}
				log.Println("got a response from couch for the view...")
				log.Println(string(bodyV))
			}

			log.Println("successfully created user")
			fmt.Fprintf(w, "{\"status\": \"logged in\", \"id\": \"" + user.Name + "\",\"token\": \"" + user.Password + "\"}")
			return
		} else {
			if !userCheck.Ok {
				b, errUCC := json.Marshal(userCheck)
				if errUCC != nil {
					log.Println(err)
				}
				log.Println(string(b))
				log.Println("user check failed")
				http.Error(w, "{\"error\":\"invalid login.\"}", http.StatusUnauthorized)
				return
			} else {
				log.Println("user already existed!!!")
				fmt.Fprintf(w, "{\"status\": \"logged in\", \"id\": \"" + user.Name + "\",\"token\": \"" + user.Password + "\"}")
				return
			}
		}
	}

	//define the request
	couch,err := http.NewRequest(http.MethodGet, Property("COUCHDB_HOST","http://localhost:5984"), http.NoBody)
	if err != nil {
		log.Println(err)
	}
	//set basic auth
	couch.SetBasicAuth(Property("COUCHDB_ADMIN_USER","admin"), Secret("COUCHDB_ADMIN_PW"))


	//make the actual call.
	dResponse,err := client.Do(couch)
	//dResponse,err := http.Header{}.Get("http://localhost:5984/")
	if err != nil {
		log.Println("problem talking to couch: " + err.Error())
		http.Error(w, "{\"error\": \"failed to talk to the db\"}", http.StatusBadRequest)
		return
	}
	defer dResponse.Body.Close()
	body,err := ioutil.ReadAll(dResponse.Body)
	if err != nil {
		log.Println("problem reading db response: " + err.Error())
		http.Error(w, "{\"error\": \"invalid request\"}", http.StatusInternalServerError)
		return
	}
	log.Println("got a response from couch")
	log.Println(string(body))
	fmt.Fprintf(w, string(body))
}

package golodex

import (
	"fmt"
	"log"
	"net/http"
	"plugin"
	"strings"
)

type Auth struct {
	Name string
	Email  string
	Token string
	Id string
}

type Pager interface {
	Page(w http.ResponseWriter, req *http.Request)
	Methods() map[string]bool
	Authenticate(req *http.Request) bool
}

func Page(w http.ResponseWriter, req *http.Request) {
	log.Println("-------------------api page!!!-------------------")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var path = req.URL.Path
	reqM := req.Method
	//var param = req.URL.Query().Get("var1")
	//fmt.Fprintf(w, "greetings: " + path +":" + param +"\n")
	fn := strings.Replace(path, "/api", "", 1)
	var mod = "./pages" + path + fn + ".so"

	log.Println("plugin path: " + mod)

	plug, err := plugin.Open(mod)
	if err != nil {
		log.Println("error"+err.Error())
		log.Println("mod not found: " + mod)
		WriteResponse(w, "{\"error\":\"the page was not found\""+err.Error()+"}", http.StatusNotFound)
		return
	}

	log.Println("opened page...")
	// 2. look up a symbol (an exported function or variable)
	page, err := plug.Lookup("Pager")
	if err != nil {
		log.Println("invalid page request, plug not found: " + path)
		log.Println(err)
		WriteResponse(w, "{\"error\":\"page not found\"}", http.StatusNotImplemented)
		return
	}

	log.Println("found page")
	// 3. Assert that loaded symbol is of a desired type
	var pager Pager
	pager, ok := page.(Pager)
	if !ok {
		log.Println("problem with mod")
		log.Println(ok)
		log.Println("check error above.")
		WriteResponse(w, "{\"error\":\"problem loading page\"}", http.StatusInternalServerError)
		return
	}
	methods := pager.Methods()

	if reqM == "OPTIONS" {
		if methods[req.Header.Get("Access-Control-Request-Method")] {
			log.Println("checking preflight request...")
			log.Println("preflight request approved..." + req.Header.Get("Access-Control-Request-Method"))
			w.Header().Set("Access-Control-Allow-Methods", req.Header.Get("Access-Control-Request-Method"))
			w.Header().Set("Access-Control-Allow-Headers", req.Header.Get("Access-Control-Request-Headers"))
			WriteResponse(w, "", http.StatusNoContent)
			return
		} else {
			WriteResponse(w, "", http.StatusMethodNotAllowed)
			return
		}
	}
	log.Println("loaded a page, trying to authenticate...")
	auth := pager.Authenticate(req)
	if auth {
		log.Println("api authenticated")
		// 4. Authenticate
		log.Println("method: "+reqM)
		if methods[reqM] {
			// 5. use the module if authenticated
			pager.Page(w, req)
		} else {
			log.Println("operation not supported...")
			WriteResponse(w, "{\"error\":\"operation not supported\"}", http.StatusMethodNotAllowed)
			return
		}
	} else {
		log.Println("unauthorized")
		WriteResponse(w, "{\"error\":\"unauthorized\"}", http.StatusForbidden)
		return
	}
}

func WriteResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, message)
}
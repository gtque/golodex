package golodexdata

import (
	"fmt"
	"log"
	"net/http"
	"plugin"
	"strings"
)

type DataPager interface {
	Page(w http.ResponseWriter, req *http.Request)
	Methods() map[string]bool
	Authenticate(req *http.Request) bool
}

func DataPage(w http.ResponseWriter, req *http.Request) {
	log.Println("----------------data page!!!----------------")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("setting some headers on response writer.")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	var path = req.URL.Path
	//var param = req.URL.Query().Get("var1")
	//fmt.Fprintf(w, "greetings: " + path +":" + param +"\n")
	fn := strings.Replace(path, "/api", "", 1)
	var mod = "./pages" + path + fn + ".so"

	log.Println("the plugin path found: " + mod)
	reqM := req.Method
	plug, err := plugin.Open(mod)
	if err != nil {
		log.Println("error"+err.Error())
		log.Println("mod not found: " + mod)
		WriteResponse(w, "{\"error\":\"the page was not found\""+err.Error()+"}", http.StatusNotFound)
		return
	}

	// 2. look up a symbol (an exported function or variable)
	page, err := plug.Lookup("DataPager")
	if err != nil {
		log.Println("invalid page request, plug not found: " + path)
		log.Println(err)
		WriteResponse(w, "{\"error\":\"page not found\"}", http.StatusNotImplemented)
		return
	}

	log.Println("found page")
	// 3. Assert that loaded symbol is of a desired type
	var pager DataPager
	pager, ok := page.(DataPager)
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
	auth := pager.Authenticate(req)
	if auth {
	// 4. Authenticate
		log.Println("method: "+reqM)
		if methods[reqM] {
			// 5. use the module if authenticated
			pager.Page(w, req)
		} else {
			WriteResponse(w, "{\"error\":\"operation not supported\"}", http.StatusMethodNotAllowed)
			return
		}
	} else {
		WriteResponse(w, "{\"error\":\"unauthorized\"}", http.StatusForbidden)
		return
	}
}

func WriteResponse(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	fmt.Fprintf(w, message)
}

package main

import (
	"encoding/json"
	"fmt"
	"internal/golodexdata"
	"log"
	"net/http"
	"strings"
)

type pager string

func (p pager) Page(w http.ResponseWriter, req *http.Request) {
	var sort = req.URL.Query().Get("sort")
	var search = req.URL.Query().Get("search")
	var cards golodexdata.GolodexRolodex
	cards.Cards = []golodexdata.GolodexCard{}

	log.Println("sort: " + sort)
	log.Println("search: " + search)
	if sort == "ascending" || sort == "descending" {
		var rows golodexdata.GolodexCouchViewRows
		var data []byte
		if sort == "ascending" {
			_data, errAll := golodexdata.CouchQueryGet(req, "_design/golodex/_view/sorted_view")
			if errAll != nil {
				golodexdata.WriteResponse(w, "{\"error\": \"reading from couch\"}", http.StatusInternalServerError)
			}
			data = _data
		} else {
			_data, errAll := golodexdata.CouchQueryGet(req, "_design/golodex/_view/sorted_view?descending=true")
			if errAll != nil {
				golodexdata.WriteResponse(w, "{\"error\": \"reading from couch\"}", http.StatusInternalServerError)
			}
			data = _data
		}
		errUm := json.Unmarshal(data, &rows)
		if errUm != nil {
			log.Println("problem unmarshalling to couch rows")
			log.Println(errUm)
			golodexdata.WriteResponse(w, "{\"error\": \"parsing from couch\"}", http.StatusInternalServerError)
			return
		}
		for _, card := range rows.Rows {
			cards.Cards = append(cards.Cards, card.Doc)
		}
	} else {
		var rows golodexdata.GolodexCouchRows
		data, errAll := golodexdata.CouchQueryGet(req, "_all_docs?include_docs=true")
		if errAll != nil {
			golodexdata.WriteResponse(w, "{\"error\": \"reading from couch\"}", http.StatusInternalServerError)
		}
		errUm := json.Unmarshal(data, &rows)
		if errUm != nil {
			log.Println("problem unmarshalling to couch rows")
			log.Println(errUm)
			golodexdata.WriteResponse(w, "{\"error\": \"parsing from couch\"}", http.StatusInternalServerError)
			return
		}
		for _, card := range rows.Rows {
			cards.Cards = append(cards.Cards, card.Doc)
		}
	}

	filteredCards := searchCards(cards, search)
	rolodex, errM := json.Marshal(filteredCards)
	if errM != nil {
		log.Println("problem marshalling the rolodex")
		log.Println(errM)
		golodexdata.WriteResponse(w, "{\"error\": \"problem building rolodex\"}", http.StatusInternalServerError)
		return
	}
	//log.Println(string(rolodex))
	golodexdata.WriteResponse(w, string(rolodex), http.StatusOK)
}

func searchCards(cards golodexdata.GolodexRolodex, search string) golodexdata.GolodexRolodex {
	var filteredCards golodexdata.GolodexRolodex
	filteredCards.Cards = []golodexdata.GolodexCard{}
	if search != "" {
		searchParameters := strings.Fields(search)
		for _, card := range cards.Cards {
			cardSearch := fmt.Sprintf("%#v", card)
			cardSearch = cardSearch + "; " + card.Name.First + " "
			if card.Name.Middle != "" {
				cardSearch = cardSearch + card.Name.Middle + " "
			}
			cardSearch = cardSearch + card.Name.Last
			for _, searchParameter := range searchParameters {
				if strings.Contains(cardSearch, searchParameter) {
					filteredCards.Cards = append(filteredCards.Cards, card)
					break
				}
			}
		}
	} else {
		filteredCards = cards
	}
	return filteredCards
}
func (p pager) Methods() map[string]bool {
	return map[string]bool{"GET": true}
}

func (p pager) Authenticate(req *http.Request) bool {
	return golodexdata.CheckAuthorization(req)
}

var DataPager pager

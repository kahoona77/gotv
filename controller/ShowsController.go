package controller

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/service"
	"log"
	"net/http"
	"strings"
)

// ShowsResult response result
type ShowsResult struct {
	Success  bool                         `json:"success"`
	Status   string                       `json:"status"`
	Shows    []domain.Show                `json:"shows,omitempty"`
	Episodes map[string][]*domain.Episode `json:"episodes,omitempty"`
}

// ShowsController handles show requests
type ShowsController struct {
	*service.Context
}

// HandleRequests does what it says
func (sc *ShowsController) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()

	switch {
	case url == "/shows/load":
		sc.load(w, r)
	case url == "/shows/save":
		sc.save(w, r)
	case url == "/shows/delete":
		sc.delete(w, r)
	case url == "/shows/updateEpisodes":
		sc.updateEpisodes(w, r)
	case strings.HasPrefix(url, "/shows/search"):
		sc.search(w, r)
	case strings.HasPrefix(url, "/shows/loadEpisodes"):
		sc.loadEpisodes(w, r)
	}

	return
}

func (sc *ShowsController) load(w http.ResponseWriter, r *http.Request) {
	var results []domain.Show
	sc.ShowsRepo.All(&results)

	data := ShowsResult{true, "ok", results, nil}
	json.NewEncoder(w).Encode(data)
}

func (sc *ShowsController) save(w http.ResponseWriter, r *http.Request) {
	var show domain.Show
	data := ShowsResult{true, "ok", nil, nil}
	if readJson(r, "data", &show) {

		_, err := sc.ShowsRepo.Save(show.Id, &show)

		if err != nil {
			data.Success = false
			data.Status = "error"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (sc *ShowsController) delete(w http.ResponseWriter, r *http.Request) {
	var show domain.Show
	data := ShowsResult{true, "ok", nil, nil}
	if readJson(r, "data", &show) {
		err := sc.ShowsRepo.Remove(show.Id)

		if err != nil {
			log.Printf("%v", err)
			data.Success = false
			data.Status = "error"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (sc *ShowsController) search(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	query := params["query"][0]

	results := sc.ShowService.SearchShow(query)

	data := ShowsResult{true, "ok", results, nil}
	json.NewEncoder(w).Encode(data)
}

func (sc *ShowsController) loadEpisodes(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	showID := params["showId"][0]

	results := sc.ShowService.LoadEpisodes(showID)

	data := ShowsResult{true, "ok", nil, results}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func (sc *ShowsController) updateEpisodes(w http.ResponseWriter, r *http.Request) {
	sc.ShowService.UpdateEpisodes (sc.GetSettings())
	data := ShowsResult{true, "ok", nil, nil}
	json.NewEncoder(w).Encode(data)
}

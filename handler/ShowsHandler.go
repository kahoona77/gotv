package handler

import (
	"encoding/json"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/tvdb"
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

// ShowsHandler handles show requests
type ShowsHandler struct {
	showsRepo *domain.GoTvRepository
	parser    *tvdb.ShowParser
	settings  *domain.XtvSettings
}

// NewShowsHandler creates a new ShowsHandler
func NewShowsHandler(showsRepo *domain.GoTvRepository, parser *tvdb.ShowParser, settings *domain.XtvSettings) *ShowsHandler {
	h := new(ShowsHandler)
	h.showsRepo = showsRepo
	h.parser = parser
	h.settings = settings
	return h
}

// HandleRequests does what it says
func (sh *ShowsHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()

	switch {
	case url == "/shows/load":
		sh.load(w, r)
	case url == "/shows/save":
		sh.save(w, r)
	case url == "/shows/delete":
		sh.delete(w, r)
	case url == "/shows/updateEpisodes":
		sh.updateEpisodes(w, r)
	case strings.HasPrefix(url, "/shows/search"):
		sh.search(w, r)
	case strings.HasPrefix(url, "/shows/loadEpisodes"):
		sh.loadEpisodes(w, r)
	}

	return
}

func (sh *ShowsHandler) load(w http.ResponseWriter, r *http.Request) {
	var results []domain.Show
	sh.showsRepo.All(&results)

	data := ShowsResult{true, "ok", results, nil}
	json.NewEncoder(w).Encode(data)
}

func (sh *ShowsHandler) save(w http.ResponseWriter, r *http.Request) {
	var show domain.Show
	data := ShowsResult{true, "ok", nil, nil}
	if readJson(r, "data", &show) {

		_, err := sh.showsRepo.Save(show.Id, &show)

		if err != nil {
			data.Success = false
			data.Status = "error"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (sh *ShowsHandler) delete(w http.ResponseWriter, r *http.Request) {
	var show domain.Show
	data := ShowsResult{true, "ok", nil, nil}
	if readJson(r, "data", &show) {
		err := sh.showsRepo.Remove(show.Id)

		if err != nil {
			log.Printf("%v", err)
			data.Success = false
			data.Status = "error"
		}
	}

	json.NewEncoder(w).Encode(data)
}

func (sh *ShowsHandler) search(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	query := params["query"][0]

	results := sh.parser.GetTvdbClient().SearchShow(query)

	data := ShowsResult{true, "ok", results, nil}
	json.NewEncoder(w).Encode(data)
}

func (sh *ShowsHandler) loadEpisodes(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	showID := params["showId"][0]

	results := sh.parser.GetTvdbClient().LoadEpisodes(showID)

	data := ShowsResult{true, "ok", nil, results}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("%v", err)
	}
}

func (sh *ShowsHandler) updateEpisodes(w http.ResponseWriter, r *http.Request) {
	sh.parser.UpdateEpisodes (sh.settings)
	data := ShowsResult{true, "ok", nil, nil}
	json.NewEncoder(w).Encode(data)
}

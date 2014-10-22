package handler

import (
  "encoding/json"
  "github.com/kahoona77/gotv/domain"
  "github.com/kahoona77/gotv/tvdb"
  "net/http"
  "strings"
)

type ShowsResult struct {
  Success  bool                `json:"success"`
  Status   string              `json:"status"`
  Shows    []domain.Show       `json:"shows,omitempty"`
}

type ShowsHandler struct {
  showsRepo   *domain.GoTvRepository
  tvdb        *tvdb.Client
}

func NewShowsHandler(showsRepo *domain.GoTvRepository, client *tvdb.Client) *ShowsHandler {
  h := new(ShowsHandler)
  h.showsRepo = showsRepo
  h.tvdb      = client
  return h
}

func (sh *ShowsHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
  url := r.URL.String()

  switch {
  case url == "/shows/load":
    sh.load(w, r)
  case strings.HasPrefix(url,"/shows/search"):
    sh.search(w, r)
  }

  return
}

func (sh *ShowsHandler) load (w http.ResponseWriter, r *http.Request) {
  var results []domain.Show
  sh.showsRepo.All(&results)

  data := ShowsResult{true, "ok", results}
  json.NewEncoder(w).Encode(data)
}

func (sh *ShowsHandler) search (w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  query := params["query"][0]

  results := sh.tvdb.SearchShow (query)

  data := ShowsResult{true, "ok", results}
  json.NewEncoder(w).Encode(data)
}

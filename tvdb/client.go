package tvdb

import (
  "log"
  "github.com/kahoona77/gotv/domain"
  tvdb "github.com/garfunkel/go-tvdb"
)

type Client struct {}

func NewClient() *Client {
  client := new(Client)
  return client
}

func (client *Client) SearchShow (query string) []domain.Show {
  var shows []domain.Show

  results, err := tvdb.GetSeries (query)
  if (err != nil) {
    log.Printf("error while searching show: %v", err)
  }

  shows = make([]domain.Show, len(results.Series), len(results.Series))

  for i := range results.Series {
        shows[i] = showFromSeries(results.Series[i])
  }

  return shows
}

func showFromSeries(series *tvdb.Series) domain.Show {
  show := domain.Show{}
  show.Name = series.SeriesName
  show.TvdbId = series.SeriesID
  show.Banner = series.Banner

  return show
}

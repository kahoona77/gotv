package tvdb

import (
	tvdb "github.com/garfunkel/go-tvdb"
	"github.com/kahoona77/gotv/domain"
	"log"
	"strconv"
)

// Client a TVDB-Client
type Client struct{}

// NewClient creates a new TVDB-Client
func NewClient() *Client {
	client := new(Client)
	return client
}

// SearchShow searches for a show
func (client *Client) SearchShow(query string) []domain.Show {
	var shows []domain.Show

	results, err := tvdb.GetSeries(query)
	if err != nil {
		log.Printf("error while searching show: %v", err)
	}

	shows = make([]domain.Show, len(results.Series), len(results.Series))

	for i := range results.Series {
		shows[i] = showFromSeries(results.Series[i])
	}

	return shows
}

// LoadEpisodes loads episodes for a show
func (client *Client) LoadEpisodes(showID string) map[string][]*domain.Episode {
	seriesID, err := strconv.Atoi(showID)
	series, err := tvdb.GetSeriesByID(uint64(seriesID))
	if err != nil {
		log.Printf("error while searching show: %v", err)
	}

	err = series.GetDetail()
	if err != nil {
		log.Printf("error while getting show-detail: %v", err)
	}

	seasons := make(map[string][]*domain.Episode)
	for season, eps := range series.Seasons {
		episodes := make([]*domain.Episode, len(eps), len(eps))
		for i := range eps {
			episodes[i] = episodeFromSeriesEpisode(eps[i])
		}

		seasons[strconv.Itoa(int(season))] = episodes
	}

	return seasons
}

func showFromSeries(series *tvdb.Series) domain.Show {
	show := domain.Show{}
	show.Name = series.SeriesName
	show.Id = strconv.Itoa(int(series.ID))
	show.Banner = series.Banner
	show.FirstAired = series.FirstAired
	show.Overview = series.Overview

	return show
}

func episodeFromSeriesEpisode(seriesEpisode *tvdb.Episode) *domain.Episode {
	episode := domain.Episode{}
	episode.ID = strconv.Itoa(int(seriesEpisode.ID))
	episode.Name = seriesEpisode.EpisodeName
	episode.FirstAired = seriesEpisode.FirstAired
	episode.Overview = seriesEpisode.Overview
	episode.Filename = seriesEpisode.Filename
	episode.EpisodeNumber = seriesEpisode.EpisodeNumber
	episode.SeasonNumber = seriesEpisode.SeasonNumber

	return &episode
}

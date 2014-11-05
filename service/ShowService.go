package service

import (
	"fmt"
	tvdb "github.com/garfunkel/go-tvdb"
	"github.com/kahoona77/gotv/domain"
	"labix.org/v2/mgo/bson"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ShowInfo Info about a show
type ShowInfo struct {
	Name    string
	Season  int
	Episode int
}

//ShowService parses and moves TV shows
type ShowService struct {
	*Context
	seRegex *regexp.Regexp
	xRegex  *regexp.Regexp
}

// NewShowService creates a new ShowService
func NewShowService(ctx *Context) *ShowService {
	p := new(ShowService)
	p.Context = ctx
	p.seRegex, _ = regexp.Compile(`(.*?)[.\s][sS](\d{2})[eE](\d{2}).*`)
	p.xRegex, _ = regexp.Compile(`(.*?)[.\s](\d{1,2})[xX](\d{2}).*`)
	return p
}

// MoveEpisode moves the epsiode to its season folder
func (showService *ShowService) MoveEpisode(file string, settings *domain.XtvSettings) {
	info := showService.parseShow(file)

	if info != nil {
		show, episode := showService.getShowData(info)

		if show == nil || episode == nil {
			//error
			return
		}

		// create output file
		fileEnding := file[strings.LastIndex(file, "."):]
		destinationFolder := settings.ShowsFolder + "/" + show.Folder + "/Season " + strconv.Itoa(int(episode.SeasonNumber)) + "/"
		fileName := show.Name + " - " + strconv.Itoa(int(episode.SeasonNumber)) + "x" + fmt.Sprintf("%0.2d", episode.EpisodeNumber) + " - " + episode.Name

		//move epsiode to destination
		srcFile := filepath.FromSlash(file)
		destFile := filepath.FromSlash(destinationFolder + fileName + fileEnding)
		err := os.Rename(srcFile, destFile)
		if err != nil {
			log.Printf("Error while moving epsiode to destination: %s", err)
			return
		}

		log.Printf("Moved Episode '%s' to folder '%s'", fileName, destinationFolder)
	}
}

func (showService *ShowService) getShowData(info *ShowInfo) (*domain.Show, *domain.Episode) {
	// find show
	var shows []domain.Show
	query := bson.M{"searchName": info.Name}
	err := showService.ShowsRepo.FindWithQuery(&query, &shows)
	if err != nil || len(shows) <= 0 {
		log.Printf("could not find show: %v", info.Name)
		return nil, nil
	}
	show := shows[0]

	// find Episode
	var episode *domain.Episode
	seasons := showService.LoadEpisodes(show.Id)
	episodes := seasons[strconv.Itoa(int(info.Season))]
	for i := range episodes {
		if int(episodes[i].EpisodeNumber) == info.Episode {
			episode = episodes[i]
		}
	}

	if episode == nil {
		log.Printf("could not find episode: %v for show %v", info.Episode, info.Name)
		return nil, nil
	}

	return &show, episode
}

func sanitizeFilename(filename string) string {
	// Remove all strange characters
	seps, err := regexp.Compile(`[&_=+:]`)
	if err == nil {
		filename = seps.ReplaceAllString(filename, "")
	}

	return filename
}

func (showService *ShowService) parseShow(absoluteFile string) *ShowInfo {
	info := new(ShowInfo)
	// cut off the path
	file := absoluteFile[strings.LastIndex(absoluteFile, "/")+1:]

	// Replace all _ with dots
	file = strings.Replace(file, "_", ".", -1)

	result := showService.seRegex.FindStringSubmatch(file)
	if result != nil {
		info.Name = strings.Replace(result[1], ".", " ", -1)
		info.Season, _ = strconv.Atoi(result[2])
		info.Episode, _ = strconv.Atoi(result[3])
	} else {
		// try othe rpattern
		result = showService.xRegex.FindStringSubmatch(file)
		if result != nil {
			info.Name = strings.Replace(result[1], ".", " ", -1)
			info.Season, _ = strconv.Atoi(result[2])
			info.Episode, _ = strconv.Atoi(result[3])
		} else {
			return nil
		}
	}

	return info
}

func (showService *ShowService) UpdateEpisodes(settings *domain.XtvSettings) {
	// iterate over files in downlod-Dir
	err := filepath.Walk(settings.DownloadDir, func(file string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			showService.MoveEpisode(file, settings)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error while updating episodes: %v", err)
	}

	UpdateKodi(settings)
}

func UpdateKodi(settings *domain.XtvSettings) {
	//connect
	conn, err := net.Dial("tcp", settings.KodiAddress)
	if err != nil {
		log.Printf("Error while connecting to Kodi: %v", err)
		return
	}
	defer conn.Close()

	msg := `{"id":1,"method":"VideoLibrary.Scan","params":[],"jsonrpc":"2.0"}`
	// json.NewEncoder(conn).Encode(data)
	if _, err = conn.Write([]byte(msg)); err != nil {
		log.Printf("Error while sending update command to Kodi: %s", err)
	}
}

// SearchShow searches for a show
func (showService *ShowService) SearchShow(query string) []domain.Show {
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
func (showService *ShowService) LoadEpisodes(showID string) map[string][]*domain.Episode {
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

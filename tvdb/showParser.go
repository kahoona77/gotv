package tvdb

import (
	"github.com/kahoona77/gotv/domain"
	"labix.org/v2/mgo/bson"
	"log"
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

//ShowParser parses and moves TV shows
type ShowParser struct {
	showsRepo *domain.GoTvRepository
	tvdb      *Client
	seRegex   *regexp.Regexp
	xRegex    *regexp.Regexp
}

// NewShowParser creates a new ShowParser
func NewShowParser(showsRepo *domain.GoTvRepository, client *Client) *ShowParser {
	p := new(ShowParser)
	p.showsRepo = showsRepo
	p.tvdb = client
	p.seRegex, _ = regexp.Compile(`(.*?)[.\s][sS](\d{2})[eE](\d{2}).*`)
	p.xRegex, _ = regexp.Compile(`(.*?)[.\s](\d{1,2})[xX](\d{2}).*`)
	return p
}

// MoveEpisode moves the epsiode to its season folder
func (parser *ShowParser) MoveEpisode(file string, settings *domain.XtvSettings) {
	info := parser.parseShow(file)

	show, episode := parser.getShowData(info)

	// create output file
	fileEnding := file[strings.LastIndex(file, "."):]
	destinationFolder := settings.ShowsFolder + "/" + show.Folder + "/Season " + strconv.Itoa(int(episode.SeasonNumber)) + "/"
	fileName := show.Name + " - " + strconv.Itoa(int(episode.SeasonNumber)) + "x" + strconv.Itoa(int(episode.EpisodeNumber)) + " - " + episode.Name 

	//move epsiode to destination
	srcFile := filepath.FromSlash(file)
	destFile := filepath.FromSlash(destinationFolder + fileName + fileEnding)
	err := os.Rename(srcFile, destFile)
	if err != nil {
		log.Printf("Error while moving epsiode to destination: %s", err)
	}
}

func (parser *ShowParser) getShowData(info *ShowInfo) (*domain.Show, *domain.Episode) {
	// find show
	var shows []domain.Show
	query := bson.M{"name": info.Name}
	err := parser.showsRepo.FindWithQuery(&query, &shows)
	if err != nil {
		log.Printf("could not find show: %v", info.Name)
		return nil, nil
	}
	show := shows[0]

	// find Episode
	var episode *domain.Episode
	seasons := parser.tvdb.LoadEpisodes(show.Id)
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

func (parser *ShowParser) parseShow(absoluteFile string) *ShowInfo {
	info := new(ShowInfo)

	// cut off the path
	file := absoluteFile[strings.LastIndex(absoluteFile, "/")+1:]

	result := parser.seRegex.FindStringSubmatch(file)
	if result != nil {
		info.Name = strings.Replace(result[1], ".", " ", -1)
		info.Season, _ = strconv.Atoi(result[2])
		info.Episode, _ = strconv.Atoi(result[3])
	} else {
		// try othe rpattern
		result = parser.xRegex.FindStringSubmatch(file)
		if result != nil {
			info.Name = strings.Replace(result[1], ".", " ", -1)
			info.Season, _ = strconv.Atoi(result[2])
			info.Episode, _ = strconv.Atoi(result[3])
		}
	}

	return info
}

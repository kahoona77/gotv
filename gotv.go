package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/handler"
	"github.com/kahoona77/gotv/irc"
	"github.com/kahoona77/gotv/tvdb"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// command line flags
	port := flag.Int("port", 8080, "port to serve on")
	logFile := flag.String("log", "xtv.log", "log-file")
	flag.Parse()

	// setup log
	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	//creating db
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)

	fs := http.Dir("web")
	fileHandler := http.FileServer(fs)
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", fileHandler))

	//Repositories
	serverRepo := domain.NewRepository(session, "servers")
	settingsRepo := domain.NewRepository(session, "settings")
	packetsRepo := domain.NewRepository(session, "packets")
	showsRepo := domain.NewRepository(session, "shows")

	//load Settings
	var settings domain.XtvSettings
	settingsRepo.FindFirst(&settings)
	settings.LogFile = *logFile

	//TVDB-Client
	tvdbClient := tvdb.NewClient()

	//Parser
	parser := tvdb.NewShowParser(showsRepo, tvdbClient)

	//IrcClient
	ircClient := irc.NewClient(packetsRepo, serverRepo, &settings)

	//DccServie
	dccService := irc.NewDccService(ircClient, parser)
	dccService.UpdateSettings(&settings)
	ircClient.DccService = dccService

	//Handlers
	dataHandler := handler.NewDataHandler(serverRepo, settingsRepo, dccService)
	r.PathPrefix("/data/").HandlerFunc(dataHandler.HandleRequests)

	packetsHandler := handler.NewPacketsHandler(packetsRepo)
	r.PathPrefix("/packets/").HandlerFunc(packetsHandler.HandleRequests)

	showssHandler := handler.NewShowsHandler(showsRepo, parser, &settings)
	r.PathPrefix("/shows/").HandlerFunc(showssHandler.HandleRequests)

	ircHandler := handler.NewIrcHandler(ircClient)
	r.PathPrefix("/irc/").HandlerFunc(ircHandler.HandleRequests)

	downloadsHandler := handler.NewDownloadsHandler(dccService)
	r.PathPrefix("/downloads/").HandlerFunc(downloadsHandler.HandleRequests)

	log.Printf("XTV (Go) started port %d\n", *port)
	fmt.Printf("XTV (Go) started port %d\n", *port)
	addr := fmt.Sprintf(":%d", *port)
	// this call blocks -- the progam runs here forever
	http.Handle("/", r)
	err = http.ListenAndServe(addr, nil)
	log.Println(err.Error())
}

func notFound(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if p == "/" || strings.HasPrefix(p, "/home") || strings.HasPrefix(p, "/search") || strings.HasPrefix(p, "/downloads") || strings.HasPrefix(p, "/logFile") || strings.HasPrefix(p, "/shows") {
		body, _ := ioutil.ReadFile("./web/index.html")
		fmt.Fprintf(w, string(body))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 Not found.")
}

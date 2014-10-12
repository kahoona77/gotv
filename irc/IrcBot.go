package irc

import (
	irc "github.com/fluffle/goirc/client"
	"github.com/kahoona77/gotv/domain"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type IrcBot struct {
	PacketsRepo *domain.GoTvRepository
	DccService  *DccService
	Settings    *domain.XtvSettings
	Server      *domain.Server
	Conn        *irc.Conn
	ConsoleLog  []string
	Regex       *regexp.Regexp
	LogCount    int
}

func NewIrcBot(packetsRepo *domain.GoTvRepository, dccService *DccService, settings *domain.XtvSettings, server *domain.Server) *IrcBot {
	bot := new(IrcBot)
	bot.PacketsRepo = packetsRepo
	bot.DccService = dccService
	bot.Settings = settings
	bot.Server = server
	bot.ConsoleLog = make([]string, 0)
	bot.Regex, _ = regexp.Compile(`(#[0-9]+).*\[\s*([0-9|\.]+[BbGgiKMs]+)\]\s+(.+).*`)
	return bot
}

func (this *IrcBot) IsConnected() bool {
	if this.Conn == nil {
		return false
	}
	return this.Conn.Connected()
}

func (this *IrcBot) Connect() {
	// create a config and fiddle with it first:
	cfg := irc.NewConfig(this.Settings.Nick)
	cfg.Server = this.Server.Name + ":" + strconv.Itoa(this.Server.Port)
	this.Conn = irc.Client(cfg)

	// Join channels
	this.Conn.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) {
			this.logToConsole("connected to " + this.Server.Name + ":" + strconv.Itoa(this.Server.Port))

			for _, channel := range this.Server.Channels {
				this.logToConsole("joining channel " + channel.Name)
				conn.Join(channel.Name)
			}
		})

	// Parse Messages
	this.Conn.HandleFunc("PRIVMSG", this.parseMessage)

	this.Conn.HandleFunc("372", this.log372)

	this.Conn.HandleFunc("CTCP", this.DccService.handleDCC)

	// Tell client to connect.
	log.Print("connecting")
	if err := this.Conn.Connect(); err != nil {
		log.Printf("Connection error: %v\n", err)
		this.logToConsole("Connection error: " + err.Error())
	}
}

func (this *IrcBot) Disconnect() {
	// TODO
	//this.Conn.shutdown()
}

func (this *IrcBot) log372(conn *irc.Conn, line *irc.Line) {
	this.logToConsole(line.Text())
}

func (this *IrcBot) parseMessage(conn *irc.Conn, line *irc.Line) {
	packet := this.parsePacket(conn, line)
	if packet == nil {
		// log message
		this.logToConsole(line.Text())
	}
}

func (this *IrcBot) parsePacket(conn *irc.Conn, line *irc.Line) *domain.Packet {
	result := this.Regex.FindStringSubmatch(line.Text())
	if result == nil {
		return nil
	}

	fileName := cleanFileName (result[3])
	packet := domain.NewPacket(result[1], result[2], fileName, line.Nick, line.Target(), this.Server.Name, line.Time)

	//save packet
	this.PacketsRepo.Save(packet.Id, packet)

	return packet
}

func cleanFileName (filename string) string {
	return strings.Trim(filename, "\u263B\u263C\u0002\u000f ")
}

func (this *IrcBot) logToConsole(msg string) {
	if this.LogCount > 500 {
		this.ConsoleLog = make([]string, 0)
		this.LogCount = 0
	}
	this.ConsoleLog = append(this.ConsoleLog, msg)
	this.LogCount++
}

func (this *IrcBot) StartDownload(download *Download) {
	this.logToConsole("Starting Download: " + download.File)
	msg := "xdcc send " + getCleanPacketId(download)
	this.Conn.Privmsg(download.Bot, msg)
}

func (this *IrcBot) StopDownload(download *Download) {
	this.logToConsole("Stopping Download: " + download.File)
	msg := "xdcc cancel"
	this.Conn.Privmsg(download.Bot, msg)
}

func getCleanPacketId(download *Download) string {
	return strings.Replace(download.PacketId, "#", "", -1)
}

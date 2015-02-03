package service

import (
	irc "github.com/fluffle/goirc/client"
	"github.com/kahoona77/gotv/domain"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type IrcBot struct {
	*Context
	Server     *domain.Server
	Conn       *irc.Conn
	ConsoleLog []string
	Regex      *regexp.Regexp
	LogCount   int
}

func NewIrcBot(ctx *Context, server *domain.Server) *IrcBot {
	bot := new(IrcBot)
	bot.Context = ctx
	bot.Server = server
	bot.ConsoleLog = make([]string, 0)
	bot.Regex, _ = regexp.Compile(`(#[0-9]+).*\[\s*([0-9|\.]+[BbGgiKMs]+)\]\s+(.+).*`)
	return bot
}

func (ib *IrcBot) IsConnected() bool {
	if ib.Conn == nil {
		return false
	}
	return ib.Conn.Connected()
}

func (ib *IrcBot) Connect() {
	//reset log
	ib.ConsoleLog = make([]string, 0)

	// create a config and fiddle with it first:
	cfg := irc.NewConfig(ib.GetSettings().Nick)
	cfg.Server = ib.Server.Name + ":" + strconv.Itoa(ib.Server.Port)
	ib.Conn = irc.Client(cfg)

	// Join channels
	ib.Conn.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) {
			ib.logToConsole("connected to " + ib.Server.Name + ":" + strconv.Itoa(ib.Server.Port))

			for _, channel := range ib.Server.Channels {
				ib.logToConsole("joining channel " + channel.Name)
				conn.Join(channel.Name)
			}
		})

	// Parse Messages
	ib.Conn.HandleFunc("PRIVMSG", ib.parseMessage)

	ib.Conn.HandleFunc("372", ib.log372)

	ib.Conn.HandleFunc("DISCONNECTED", ib.reconnect)

	ib.Conn.HandleFunc("CTCP", ib.DccService.handleDCC)

	// Tell client to connect.
	log.Printf("Connecting to '%v'", ib.Server.Name)
	if err := ib.Conn.Connect(); err != nil {
		log.Printf("Connection error: %v\n", err)
		ib.logToConsole("Connection error: " + err.Error())
	}
}

func (ib *IrcBot) Disconnect() {
	// TODO
	//ib.Conn.shutdown()
}

func (ib *IrcBot) reconnect(conn *irc.Conn, line *irc.Line) {
	log.Printf("Discconected from '%v'. Reconnecting now ...", ib.Server.Name)
	ib.Connect()
}

func (ib *IrcBot) log372(conn *irc.Conn, line *irc.Line) {
	ib.logToConsole(line.Text())
}

func (ib *IrcBot) parseMessage(conn *irc.Conn, line *irc.Line) {
	ib.parsePacket(conn, line)
}

func (ib *IrcBot) parsePacket(conn *irc.Conn, line *irc.Line) *domain.Packet {
	result := ib.Regex.FindStringSubmatch(line.Text())
	if result == nil {
		return nil
	}

	fileName := cleanFileName(result[3])
	packet := domain.NewPacket(result[1], result[2], fileName, line.Nick, line.Target(), ib.Server.Name, line.Time)

	//save packet
	if packet != nil {
		ib.PacketsRepo.Save(packet.Id, packet)
	}

	return packet
}

func cleanFileName(filename string) string {
	return strings.Trim(filename, "\u263B\u263C\u0002\u000f ")
}

func (ib *IrcBot) logToConsole(msg string) {
	if ib.LogCount > 500 {
		ib.ConsoleLog = make([]string, 0)
		ib.LogCount = 0
	}
	ib.ConsoleLog = append(ib.ConsoleLog, msg)
	ib.LogCount++
}

func (ib *IrcBot) StartDownload(download *Download) {
	ib.logToConsole("Starting Download: " + download.File)
	msg := "xdcc send " + getCleanPacketId(download)
	ib.Conn.Privmsg(download.Bot, msg)
}

func (ib *IrcBot) StopDownload(download *Download) {
	ib.logToConsole("Stopping Download: " + download.File)
	msg := "xdcc cancel"
	ib.Conn.Privmsg(download.Bot, msg)
}

func getCleanPacketId(download *Download) string {
	return strings.Replace(download.PacketId, "#", "", -1)
}

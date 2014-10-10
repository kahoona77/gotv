package irc

import (
	"bufio"
	"encoding/binary"
	irc "github.com/fluffle/goirc/client"
	"github.com/kahoona77/gotv/domain"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type IrcBot struct {
	PacketsRepo *domain.GoTvRepository
	Settings    *domain.XtvSettings
	Server      *domain.Server
	Conn        *irc.Conn
	ConsoleLog  []string
	Regex       *regexp.Regexp
	LogCount    int
}

func NewIrcBot(packetsRepo *domain.GoTvRepository, settings *domain.XtvSettings, server *domain.Server) *IrcBot {
	bot := new(IrcBot)
	bot.PacketsRepo = packetsRepo
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

	this.Conn.HandleFunc("CTCP", this.handleDCC)

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

func (this *IrcBot) handleDCC(conn *irc.Conn, line *irc.Line) {
	log.Printf("CTCP-> cmd: %v, args: %v, src: %v\n", line.Cmd, line.Args, line.Src)
	request := strings.Split(line.Args[2], " ")
	ctcpType := line.Args[0]

	log.Printf("req1: %v, req2: %v, req3: %v\n", request[0], request[1], request[2])

	if ctcpType == "DCC" {
		cmd := request[0]
		if cmd == "SEND" {
			log.Printf("DCC SEND\n")
			fileName := request[1]
			addrInt, _ := strconv.ParseInt(request[2], 0, 64)
			address := inet_ntoa(addrInt)
			port := request[3]
			size := request[4]

			log.Printf("file: %v, addr: %v, port: %v, size:%v\n", fileName, address.String(), port, size)

			//write to file
			// open output file
			fo, err := os.Create(fileName)
			if err != nil {
				log.Printf("File create error: %s", err)
			}
			// make a write buffer
			w := bufio.NewWriter(fo)

			thepath, err := filepath.Abs(filepath.Dir(fo.Name()))

			log.Println("file created: " + thepath + fo.Name())

			// close fo on exit and check for its returned error
			defer func() {
				if err := fo.Close(); err != nil {
					panic(err)
				}
			}()

			conn, err := net.Dial("tcp", address.String()+":"+port)
			if err != nil {
				log.Printf("Connect error: %v", err)
			}

			var inBuf = make([]byte, 1024)
			totalBytes := 0
			log.Println("reading...")
			for {
				n, err := conn.Read(inBuf)
				totalBytes += n
				// Was there an error in reading ?
				if err != nil {
					log.Printf("Read error: %s", err)
					break
				}

				log.Printf("transferred: %d", totalBytes)

				// write a chunk
				if _, err := w.Write(inBuf[:n]); err != nil {
					panic(err)
				}

				log.Println("file written...")

				outBuf := makeOutBuffer(totalBytes)

				log.Printf("out buffer: %v", outBuf)

				n, err = conn.Write(outBuf)
				if err != nil {
					log.Printf("Write error: %s", err)
					break
				}

				log.Println("out written...")

				if err = w.Flush(); err != nil {
					panic(err)
				}

				log.Println("flushed")

			}
			log.Printf("%d bytes read", totalBytes)
			conn.Close()

		}
	}

}

func makeOutBuffer(totalBytes int) []byte {
	var bytes = make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(totalBytes))
	return bytes
}

// Convert uint to net.IP
func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
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

	packet := domain.NewPacket(result[1], result[2], result[3], line.Nick, line.Target(), this.Server.Name, line.Time)

	//save packet
	this.PacketsRepo.Save(packet.Id, packet)

	return packet
}

func (this *IrcBot) logToConsole(msg string) {
	if this.LogCount > 500 {
		this.ConsoleLog = make([]string, 0)
		this.LogCount = 0
	}
	this.ConsoleLog = append(this.ConsoleLog, msg)
	this.LogCount++
}

func (this *IrcBot) DownloadPacket(packet *domain.Packet) {
	msg := "xdcc send " + getCleanPacketId(packet)
	this.Conn.Privmsg(packet.Bot, msg)
}

func getCleanPacketId(packet *domain.Packet) string {
	return strings.Replace(packet.PacketId, "#", "", -1)
}

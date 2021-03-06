package service

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	irc "github.com/fluffle/goirc/client"
	"github.com/juju/ratelimit"
	"github.com/kahoona77/gotv/domain"
)

// DccFileEvent -
type DccFileEvent struct {
	Type     string
	FileName string
	IP       net.IP
	Port     string
	Size     int64
}

// DccUpdate -
type DccUpdate struct {
	File       string
	TotalBytes int64
	Size       int64
}

// DccService -
type DccService struct {
	*Context
	downloads  map[string]*Download
	resumes    map[string]*DccFileEvent
	updateChan chan DccUpdate
	bucket     *ratelimit.Bucket
}

// NewDccService creates a new DccService
func NewDccService(ctx *Context) *DccService {
	dcc := new(DccService)
	dcc.Context = ctx
	dcc.downloads = make(map[string]*Download)
	dcc.resumes = make(map[string]*DccFileEvent)
	dcc.updateChan = make(chan DccUpdate)
	dcc.SetDownloadLimit(0)

	//start download update
	go dcc.updateDownloads()
	return dcc
}

func (dcc *DccService) handleDCC(conn *irc.Conn, line *irc.Line) {
	request := strings.Split(line.Args[2], " ")
	ctcpType := line.Args[0]
	settings := dcc.GetSettings()
	if ctcpType == "DCC" {
		cmd := request[0]
		if cmd == "SEND" {
			dcc.handleSend(request, conn, line, settings)
		} else if cmd == "ACCEPT" {
			dcc.handleAccept(request, settings)
		} else {
			log.Printf("received unmatched DCC command: %v", cmd)
		}
	}
}

func (dcc *DccService) handleSend(request []string, conn *irc.Conn, line *irc.Line, settings *domain.XtvSettings) {
	fileName := request[1]
	addrInt, _ := strconv.ParseInt(request[2], 0, 64)
	address := inetNtoa(addrInt)
	port := request[3]
	size, _ := strconv.ParseInt(request[4], 0, 64)

	log.Printf("received SEND - file: %v, addr: %v, port: %v, size:%v\n", fileName, address.String(), port, size)
	fileEvent := DccFileEvent{"SEND", fileName, address, port, size}

	resume, startPos := dcc.fileExists(&fileEvent, settings)

	if resume {
		// file already exists -> send resume request
		msg := fileName + " " + port + " " + strconv.FormatInt(startPos, 10)
		log.Printf("sending resume [%v]", msg)
		conn.Ctcp(line.Nick, "DCC RESUME", msg)
		//add to resumes
		dcc.resumes[fileEvent.FileName] = &fileEvent
	} else {
		// This is a new file start from beginning
		go dcc.startDownload(&fileEvent, startPos, settings)
	}
}

func (dcc *DccService) handleAccept(request []string, settings *domain.XtvSettings) {
	log.Printf("received ACCEPT")

	fileName := request[1]
	//port := request[2]
	position, err := strconv.ParseInt(request[3], 10, 64)

	if err != nil {
		log.Printf("error while parsing position %v", err)
		return
	}

	//find resume
	fileEvent := dcc.resumes[fileName]
	delete(dcc.resumes, fileName)
	if fileEvent == nil {
		log.Printf("can not find resume for %v", fileName)
		return
	}

	go dcc.startDownload(fileEvent, position, settings)
}

func (dcc *DccService) startDownload(fileEvent *DccFileEvent, startPos int64, settings *domain.XtvSettings) {
	file := dcc.getTempFile(fileEvent, settings)

	// set start position
	var totalBytes int64
	totalBytes = startPos
	file.Seek(startPos, 0)

	// make a write buffer
	w := bufio.NewWriter(file)

	// close file on exit and check for its returned error
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	//connect
	conn, err := net.Dial("tcp", fileEvent.IP.String()+":"+fileEvent.Port)
	if err != nil {
		log.Printf("Connect error: %v", err)
		return
	}

	//set rateLimit
	var reader = ratelimit.Reader(conn, dcc.bucket)

	var complete bool
	var inBuf = make([]byte, 1024)
	counter := 0

	//read-loop
	for {
		//read a chunk
		n, err := reader.Read(inBuf)
		// fmt.Printf("read: %d\n", n)
		if err != nil {
			if err == io.EOF {
				complete = true
			} else {
				log.Printf("Read error: %s", err)
			}
			break
		}
		totalBytes += int64(n)

		// write to File
		if _, err := w.Write(inBuf[:n]); err != nil {
			log.Printf("Write to file error: %s", err)
			break
		}

		//Send back an acknowledgement of how many bytes we have got so far.
		//Convert bytesTransfered to an "unsigned, 4 byte integer in network byte order", per DCC specification
		outBuf := makeOutBuffer(totalBytes)
		if _, err = conn.Write(outBuf); err != nil {
			log.Printf("Write error: %s", err)
			break
		}

		if err = w.Flush(); err != nil {
			log.Printf("Flush error: %s", err)
			break
		}

		counter++
		//every 100kb update
		if counter == 100 {
			dcc.updateChan <- DccUpdate{fileEvent.FileName, totalBytes, fileEvent.Size}

			//reset rateLimit
			reader = ratelimit.Reader(conn, dcc.bucket)

			counter = 0
		}
	}
	conn.Close()

	if complete {
		dcc.completeDownload(fileEvent.FileName)
	} else {
		dcc.failDownload(fileEvent.FileName)
	}

}

func (dcc *DccService) getTempFile(fileEvent *DccFileEvent, settings *domain.XtvSettings) *os.File {
	filename := filepath.FromSlash(settings.TempDir + "/" + fileEvent.FileName)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fo, err := os.Create(filename)
		if err != nil {
			log.Printf("File create error: %s", err)
		}
		return fo
	}

	fo, err := os.OpenFile(filename, os.O_WRONLY, 0777)
	if err != nil {
		log.Printf("File open error: %s", err)
	}
	return fo

}

func (dcc *DccService) fileExists(fileEvent *DccFileEvent, settings *domain.XtvSettings) (bool, int64) {
	filename := filepath.FromSlash(settings.TempDir + "/" + fileEvent.FileName)
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, 0
	}
	return true, info.Size()
}

func makeOutBuffer(totalBytes int64) []byte {
	var bytes = make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(totalBytes))
	return bytes
}

// Convert uint to net.IP
func inetNtoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// SetDownloadLimit - Sets the downloadlimit in KiloByte / Second
func (dcc *DccService) SetDownloadLimit(kBs int) {
	if kBs <= 0 {
		dcc.bucket = ratelimit.NewBucketWithQuantum(100*time.Millisecond, 1024*1000000, 100000000000)
		// dcc.connPool.SetBandwidth(iothrottler.Unlimited)
		log.Printf("download unlimited")
	} else {
		// dcc.connPool.SetBandwidth(iothrottler.Kbps * iothrottler.Bandwidth(kBs*8))
		dcc.bucket = ratelimit.NewBucketWithQuantum(10*time.Millisecond, int64(1000*kBs), int64(1000*(kBs/100)))
		log.Printf("currentDownloadLimit: %v kB/s", kBs)
	}
}

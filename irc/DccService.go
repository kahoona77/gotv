package irc

import (
  "bufio"
  "encoding/binary"
  "github.com/kahoona77/gotv/domain"
  irc "github.com/fluffle/goirc/client"
  "log"
  "net"
  "os"
  "io"
  "path/filepath"
  "strconv"
  "strings"
)

type DccFileEvent struct{
  Type      string
  FileName  string
  IP        net.IP
  Port      string
  Size      int64
}

type DccUpdate struct{
  File       string
  TotalBytes int64
  Size       int64
}

type DccService struct {
  client       *IrcClient
  settings     *domain.XtvSettings
  downloads    map[string]*Download
  updateChan   chan DccUpdate
}

func NewDccService(client *IrcClient) *DccService {
  dcc := new(DccService)
  dcc.client       = client
  dcc.settings     = client.Settings
  dcc.downloads    = make(map[string]*Download)
  dcc.updateChan   = make(chan DccUpdate)

  //start download update
  go dcc.updateDownloads()
  return dcc
}


func (dcc *DccService) handleDCC(conn *irc.Conn, line *irc.Line) {
  request := strings.Split(line.Args[2], " ")
  ctcpType := line.Args[0]

  if ctcpType == "DCC" {
    cmd := request[0]
    if cmd == "SEND" {
      fileName := request[1]
      addrInt, _ := strconv.ParseInt(request[2], 0, 64)
      address := inet_ntoa(addrInt)
      port := request[3]
      size,_ := strconv.ParseInt(request[4],0,64)

      log.Printf("file: %v, addr: %v, port: %v, size:%v\n", fileName, address.String(), port, size)
      fileEvent := DccFileEvent {"SEND", fileName, address, port, size}

      resume, startPos := dcc.fileExists (&fileEvent)

      if (resume) {
        // file already exists -> send resume request
        msg:=" :\u0001" + "DCC RESUME " + fileName + " " + port + " " + strconv.FormatInt(startPos, 10) + "\u0001"
        conn.Privmsg(line.Nick, msg)
      } else {
        // This is a new file start from beginning
        go dcc.startDownload (&fileEvent, startPos)
      }
    } else if cmd == "RESUME" {
      log.Printf("received RESUME")
    } else {
      log.Printf("received unmatched DCC command: %v", cmd)
    }
  }
}

func (dcc *DccService) startDownload (fileEvent *DccFileEvent, startPos int64) {
  file := dcc.getTempFile (fileEvent)

  // set start position
  var totalBytes int64
  totalBytes = startPos
  file.Seek (startPos, 0)

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

  var complete bool
  var inBuf = make([]byte, 1024)
  counter := 0

  //read-loop
  for {
    //read a chunk
    n, err := conn.Read(inBuf)
    if err != nil {
      if (err == io.EOF) {
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
    if (counter == 500) {
      dcc.updateChan <- DccUpdate{fileEvent.FileName, totalBytes, fileEvent.Size}
      counter = 0;
    }
  }
  conn.Close()

  if (complete) {
    dcc.completeDownload (fileEvent.FileName)
  }

}

func (dcc *DccService) getTempFile (fileEvent *DccFileEvent) *os.File{
  filename := filepath.FromSlash (dcc.settings.TempDir + "/" + fileEvent.FileName)

  fo, err := os.Create(filename)
  if err != nil {
    log.Printf("File create error: %s", err)
  }

  return fo
}

func (dcc *DccService) fileExists (fileEvent *DccFileEvent) (bool, int64) {
  filename := filepath.FromSlash (dcc.settings.TempDir + "/" + fileEvent.FileName)
  info, err := os.Stat(filename)
  if (os.IsNotExist(err)) {
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
func inet_ntoa(ipnr int64) net.IP {
  var bytes [4]byte
  bytes[0] = byte(ipnr & 0xFF)
  bytes[1] = byte((ipnr >> 8) & 0xFF)
  bytes[2] = byte((ipnr >> 16) & 0xFF)
  bytes[3] = byte((ipnr >> 24) & 0xFF)

  return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}
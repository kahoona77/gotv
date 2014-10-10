package irc

import (
  "bufio"
  "encoding/binary"
  "github.com/kahoona77/gotv/domain"
  irc "github.com/fluffle/goirc/client"
  "log"
  "net"
  "os"
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

type DccService struct {
  client       *IrcClient
  settings     *domain.XtvSettings
  downloads    map[string]*Download
}

func NewDccService(client *IrcClient) *DccService {
  dcc := new(DccService)
  dcc.client       = client
  dcc.settings     = client.Settings
  dcc.downloads    = make(map[string]*Download)
  return dcc
}


func (dcc *DccService) handleDCC(conn *irc.Conn, line *irc.Line) {
  request := strings.Split(line.Args[2], " ")
  ctcpType := line.Args[0]

  if ctcpType == "DCC" {
    cmd := request[0]
    if cmd == "SEND" {
      log.Printf("DCC SEND\n")
      fileName := request[1]
      addrInt, _ := strconv.ParseInt(request[2], 0, 64)
      address := inet_ntoa(addrInt)
      port := request[3]
      size,_ := strconv.ParseInt(request[4],0,64)

      log.Printf("file: %v, addr: %v, port: %v, size:%v\n", fileName, address.String(), port, size)
      fileEvent := DccFileEvent {"SEND", fileName, address, port, size}

      go dcc.startDownload (&fileEvent)
    }
  }
}

func (dcc *DccService) startDownload (fileEvent *DccFileEvent) {
  file := dcc.getTempFile (fileEvent)
  // make a write buffer
  w := bufio.NewWriter(file)

  // close fo on exit and check for its returned error
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

  var inBuf = make([]byte, 1024)
  var totalBytes int64
  totalBytes = 0
  counter := 0

  //read-loop
  for {
    //read a chunk
    n, err := conn.Read(inBuf)
    if err != nil {
      log.Printf("Read error: %s", err)
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

    if (counter == 500) {
      dcc.updateDownload(fileEvent.FileName, totalBytes);
      counter = 0;
    }
  }
  conn.Close()
}

func (dcc *DccService) getTempFile (fileEvent *DccFileEvent) *os.File{
  filename := filepath.FromSlash (dcc.settings.TempDir + "/" + fileEvent.FileName)

  fo, err := os.Create(filename)
  if err != nil {
    log.Printf("File create error: %s", err)
  }

  return fo
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

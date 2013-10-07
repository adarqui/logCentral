// http://golang.org/pkg/crypto/tls/

package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io"
    "log"
	"bufio"
	"github.com/postfix/goconf"
	"strings"
	"strconv"
	"os"
	"net"
)

type BindInfo struct {
	host string
	port int
}

func (s * BindInfo) SetHost(host string) {
	s.host = host
}

func (s * BindInfo) GetHost() string {
	return s.host
}

func (s * BindInfo) SetPort(port int) {
	s.port = port
}

func (s * BindInfo) GetPort() int {
	return s.port
}

func (s * BindInfo) ConnectString() string {
	x := []string{s.host,strconv.Itoa(s.port)}
	return strings.Join(x,":")
}


func reply_monitor(conn net.Conn) {
	for {
		reply := make([]byte, 256)
		n, err := conn.Read(reply)
		if err != nil {
			fmt.Println(err)
		}
		if n == 0 {
			fmt.Println("Disconnected")
			return
		}
		log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	}
}


func main() {

	conf, err := goconf.ReadConfigFile("go.conf")
	if err != nil {
		log.Fatalf("goconf: file not found: %s", err)
	}

	_bindInfo := BindInfo{}
	host, err := conf.GetString("default", "remote_host")
	port, err := conf.GetInt("default", "port")

	_bindInfo.SetHost(host)
	_bindInfo.SetPort(port)

    cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
    if err != nil {
        log.Fatalf("server: loadkeys: %s", err)
    }

    config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
    conn, err := tls.Dial("tcp", _bindInfo.ConnectString(), &config)
    if err != nil {
        log.Fatalf("client: dial: %s", err)
    }
    defer conn.Close()

    log.Println("client: connected to: ", conn.RemoteAddr())

    state := conn.ConnectionState()
    for _, v := range state.PeerCertificates {
        fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
        fmt.Println(v.Subject)
    }
    log.Println("client: handshake: ", state.HandshakeComplete)
    log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

	go reply_monitor(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.WriteString(conn, line)
		if err != nil {
			log.Fatal(err)
		}
	}
}

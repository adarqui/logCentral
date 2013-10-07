// http://golang.org/pkg/crypto/tls/

package main

import (
    "crypto/rand"
    "crypto/tls"
    "log"
    "net"
    "crypto/x509"
	"github.com/postfix/goconf"
	"strconv"
	"strings"
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

func (s * BindInfo) ListenString() string {
    x := []string{s.host,strconv.Itoa(s.port)}
    return strings.Join(x,":")
}


func main() {

    conf, err := goconf.ReadConfigFile("go.conf")
    if err != nil {
        log.Fatalf("goconf: file not found: %s", err)
    }

    _bindInfo := BindInfo{}
    host, err := conf.GetString("default", "local_host")
    port, err := conf.GetInt("default", "port")

    _bindInfo.SetHost(host)
    _bindInfo.SetPort(port)

    cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
    if err != nil {
        log.Fatalf("server: loadkeys: %s", err)
    }
    config := tls.Config{Certificates: []tls.Certificate{cert}}
    config.Rand = rand.Reader
    service := _bindInfo.ListenString()
    listener, err := tls.Listen("tcp", service, &config)
    if err != nil {
        log.Fatalf("server: listen: %s", err)
    }
    log.Print("server: listening")
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("server: accept: %s", err)
            break
        }
        defer conn.Close()
        log.Printf("server: accepted from %s", conn.RemoteAddr())
        tlscon, ok := conn.(*tls.Conn)
        if ok {
            log.Print("ok=true")
            state := tlscon.ConnectionState()
            for _, v := range state.PeerCertificates {
                log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
            }
        }
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 512)
    for {
        log.Print("server: conn: waiting")
        n, err := conn.Read(buf)
        if err != nil {
			log.Printf("server: conn: read: %s", err)
            break
        }
        log.Printf("server: conn: echo %q\n", string(buf[:n]))
        n, err = conn.Write(buf[:n])
    }
    log.Println("server: conn: closed")
}

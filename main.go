package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

func main() {
	arg1 := flag.String("ip", "127.0.0.1", "which `IP` address to listen")
	arg2 := flag.Int("port", 9527, "which `port` to listen")
	arg3 := flag.String("mode", "tcp", "which `protocol` to use, possible value: tcp, http")
	flag.Parse()
	ip := net.ParseIP(*arg1)
	if ip == nil {
		log.Fatal("could not parse IP")
	}
	port := *arg2
	if port < 0 || port > 65535 {
		log.Fatal("port number is out of range")
	}
	addr := net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	switch mode := *arg3; mode {
	case "tcp":
		tcpEcho(addr)
	case "http":
		httpEcho(addr)
	default:
		log.Fatal("unknown mode")
	}
}

func tcpEcho(addr net.TCPAddr) {
	log.Printf("tcp-echo is runnning on tcp mode, listening %s\n", addr.String())
	l, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			io.Copy(c, c)
			c.Close()
		}(conn)
	}
}

func httpEcho(addr net.TCPAddr) {
	addrStr := addr.String()
	log.Printf("tcp-echo is runnning on http mode, listening %s\n", addrStr)
	echoHandler := func(w http.ResponseWriter, req *http.Request) {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		io.WriteString(w, string(dump))
	}
	http.HandleFunc("/", echoHandler)
	log.Fatal(http.ListenAndServe(addrStr, nil))
}

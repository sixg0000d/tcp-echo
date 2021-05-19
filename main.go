package main

import (
	"flag"
	"io"
	"log"
	"net"
)

func main() {
	arg1 := flag.String("ip", "127.0.0.1", "`IP address` to listen")
	arg2 := flag.Int("port", 9527, "`Port` to listen")
	flag.Parse()
	ip := net.ParseIP(*arg1)
	if ip == nil {
		log.Fatal("could not parse IP")
	}
	port := *arg2
	if port < 0 || port > 65535 {
		log.Fatal("port number is out of range")
	}

	tcpAddr := net.TCPAddr{
		IP:   ip,
		Port: port,
	}

	l, err := net.ListenTCP("tcp", &tcpAddr)
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

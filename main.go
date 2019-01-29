package main

import (
	"io"
	"log"
	"net"
	"os"
)

func link(a, b io.ReadWriteCloser) {
	go func() {
		io.Copy(b, a)
		a.Close()
		b.Close()
	}()
	io.Copy(a, b)
	b.Close()
	a.Close()
}

func move(locale, remote string) {
	ln, err := net.Listen("tcp", locale)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("[%s] Create new forwarding from %s\n", remote, locale)
	for {
		connl, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		connr, err := net.Dial("tcp", remote)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("[%s] Income new connection from %s\n", remote, connl.RemoteAddr())
		go link(connl, connr)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		return
	}
	for i := 0; i < len(args); i += 2 {
		locale := args[i]
		remote := args[i+1]
		go move(locale, remote)
	}
	select {}
}

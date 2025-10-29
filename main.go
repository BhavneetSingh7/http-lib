package main

import (
	"fmt"
	"net"
	"os"
	"io"
	"os/signal"
	"syscall"
	"time"
	"strings"
)


func handleSignals(ln *net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	for sig := range c {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			(*ln).Close()
			fmt.Println("closing server")
			// fmt.Println(*ln)
			syscall.Exit(0)
		}
	}
}

func handleConnection(conn *net.Conn) {
	x := make([]byte, 8192)
	n, err := (*conn).Read(x)
	if err == nil {
		http_req := string(x)[:n]
		// fmt.Println(http_req)
		res_body := "Request sent by client is:\n" + http_req
		length := len(res_body) + 2
		res := fmt.Sprintf("HTTP/1.1 200 \r\nHost:localhost\r\nContent-Length:%d\r\n\r\n" + res_body, length)
		(*conn).Write([]byte(res+"\r\n"))
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	(*conn).Close()
}


func SendReq() {
	r := "GET /hey HTTP/1.1 \r\n Host: localhost \r\n"
	// fmt.Println(r)
	d, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("failed to connect")
		return
	}
	x := make([]byte, 8192)
	d.Write([]byte(r))

	n, _ := d.Read(x)
	// fmt.Println(n)
	fmt.Println(string(x[:n]))
	d.Close()
}


func Server() {
	sock, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("error occured while creating listener")
	}
	go handleSignals(&sock)
	for {
		conn, err := sock.Accept()
		if err != nil {
			time.Sleep(1 * time.Second)
			fmt.Println("error occured while accepting conn")
			continue
		}
		go handleConnection(&conn)
	}
}

func GetLinesChannel(f io.ReadCloser) <-chan string {
	d := make([]byte, 8)
	line := ""
	ch := make(chan string)
	go func() {
		for {
			// Read
			n, err := f.Read(d)
			if err != nil {
				if err != io.EOF {fmt.Println("error occured while reading file: ", err)}
				close(ch)
				break
			}
			
			// Parse
			part := strings.Split(string(d[:n]), "\n")
			line += part[0]
			if len(part)==1 {
				if n!=8 {ch <- line}
				continue
			}

			ch <- line
			for i:=1; i<len(part)-1; i++ {
				line = part[i]
				ch <- line
			}
			line = part[len(part)-1]
		}
	}()
	return ch
}


func main() {
	// ReadFile("messages.txt")
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("error occured while listening: ", err)
		return
	}
	defer l.Close()
	go handleSignals(&l)

	for {
		conn, err := l.Accept()
		if err!=nil {
			fmt.Println("error occured while accepting connection: ", err)
			continue
		}
		fmt.Println("connection has been accepted")
		ch := GetLinesChannel(conn)
		for i:= range ch {
			fmt.Printf("%s\n", i)
		}
		conn.Close()
		fmt.Println("connection has been closed")
	}
}

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleSignals(ln *net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	for sig := range c {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			(*ln).Close()
			fmt.Println("closing server")
			fmt.Println(*ln)
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

func main() {
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
		// time.Sleep(2*time.Second)
	}

}

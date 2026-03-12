package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	var host, port string
	var timeout time.Duration
	flag.StringVar(&host, "host", "", "")
	flag.StringVar(&port, "port", "", "")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "")
	flag.Parse()
	if host == "" || port == "" {
		fmt.Fprintln(os.Stderr, "invalid host/port")
		os.Exit(1)
	}
	addr := net.JoinHostPort(host, port)
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial error:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("connected")
	errCh := make(chan error, 2)
	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if err == nil {
			if tcp, ok := conn.(*net.TCPConn); ok {
				_ = tcp.CloseWrite()
			} else {
				_ = conn.Close()
			}
			errCh <- io.EOF
			return
		}
		errCh <- err
	}()

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err == nil {
			errCh <- io.EOF
			return
		}
		errCh <- err
	}()

	err = <-errCh
	if err == io.EOF {
		return
	}

	fmt.Fprintln(os.Stderr, "io error:", err)
	os.Exit(1)
}

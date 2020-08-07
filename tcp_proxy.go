package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type Proxy interface {
	ListenAndServe() error
	Stop() error
}

type TcpProxy struct {
	*Config
	listener *net.TCPListener
}

var bufferSize = 16 * 1024

func (s *TcpProxy) handleConn(conn *net.TCPConn) {
	defer conn.Close()

	fmt.Printf("handleConn: %v\n", s)

	peerAddr, err := ResolveTCPAddr(s.remoteIP, s.remotePort)
	if err != nil {
		fmt.Println(err)
		return
	}

	peer, err := net.DialTCP("tcp", nil, peerAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer peer.Close()

	stop := make(chan bool)

	Copy := func(w *net.TCPConn, r *net.TCPConn, needstop bool) {
		data := make([]byte, bufferSize)
		for {
			nr, er := r.Read(data)
			if nr > 0 {
				nw, ew := w.Write(data[0:nr])

				if nw != nr || ew != nil {
					r.CloseRead()
					w.CloseWrite()
					break
				}
			}

			if er != nil {
				r.CloseRead()
				w.CloseWrite()
				break
			}
		}

		if needstop {
			stop <- true
		}
	}

	go Copy(conn, peer, true)

	Copy(peer, conn, false)

	<-stop
}

func ResolveTCPAddr(ip string, port uint16) (*net.TCPAddr, error) {
	var s string
	if strings.Contains(ip, ":") {
		s = fmt.Sprintf("[%s]:%d", ip, port)
	} else {
		s = fmt.Sprintf("%s:%d", ip, port)
	}
	addr, err := net.ResolveTCPAddr("tcp", s)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func (s *TcpProxy) ListenAndServe() error {
	fmt.Printf("listen: %s:%d remote: %s:%d, bufferSize: %d\n",
		s.localIP, s.localPort, s.remoteIP, s.remotePort, bufferSize)

	listenAddr, err := ResolveTCPAddr(s.localIP, s.localPort)

	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		return err
	}

	s.listener = listener

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println(err)
				break
			}

			go s.handleConn(conn)
		}
	}()
	return nil
}

func (s *TcpProxy) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return errors.New("server not started")
}

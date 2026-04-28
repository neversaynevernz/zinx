package core

import (
    "net"
)

type Server struct {
    Addr string
}

func NewServer(addr string) *Server {
    return &Server{Addr: addr}
}

func (s *Server) Start() error {
    listener, err := net.Listen("tcp", s.Addr)
    if err != nil {
        return err
    }
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go s.handleConnection(conn)
    }
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 1024)
    for {
        n, err := conn.Read(buf)
        if err != nil {
            return
        }
        _, err = conn.Write(buf[:n])
        if err != nil {
            return
        }
    }
}

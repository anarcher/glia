package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
)

type Sender struct {
	ctx       context.Context
	conn      net.Conn
	network   string
	addr      string
	connected bool
}

func NewSender(ctx context.Context,
	network, addr string,
	metricCh chan []byte) *Sender {
	s := &Sender{
		ctx:     ctx,
		network: network,
		addr:    addr,
	}

	go s.looper(metricCh)

	return s
}

func (s *Sender) Connect() error {
	c, err := net.Dial(s.network, s.addr)
	if err != nil {
		return err
	}

	s.conn = c
	s.connected = true
	return nil
}

func (s *Sender) Disconnect() error {
	if s.connected == true && s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return err
		}
		s.connected = false
	}
	return nil
}

func (s *Sender) ConnectIfNot() {
	if s.connected == false || s.conn == nil {
		if err := s.Connect(); err != nil {
			//TODO: LOG
			log.Printf("Fetch connection error: %v", err)
		}
	}
}

func (s *Sender) looper(metricCh chan []byte) {
L:
	for {
		select {
		case <-s.ctx.Done():
			break L
		case metrics := <-metricCh:
			s.ConnectIfNot()
			if _, err := s.conn.Write(metrics); err != nil {
				log.Printf("Sender: write err: %v", err)
			}
		}
	}

	if err := s.Disconnect(); err != nil {
		log.Printf("Sender: disconnect err: %v", err)
	}
}

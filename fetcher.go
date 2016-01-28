package main

import (
	"golang.org/x/net/context"

	"bytes"
	"log"
	"net"
)

type Fetcher struct {
	ctx       context.Context
	network   string
	addr      string
	conn      net.Conn
	connected bool
	flushCnt  int
}

func NewFetcher(ctx context.Context,
	network, addr string,
	fetchSignal chan struct{}, metricCh chan []byte,
	flushCnt int) *Fetcher {

	fetcher := &Fetcher{
		ctx:      ctx,
		network:  network,
		addr:     addr,
		flushCnt: flushCnt,
	}

	go fetcher.looper(fetchSignal, metricCh)

	return fetcher
}

func (f *Fetcher) Connect() error {
	c, err := net.Dial(f.network, f.addr)
	if err != nil {
		return err
	}

	f.conn = c
	f.connected = true
	return nil
}

func (f *Fetcher) Disconnect() error {
	if f.connected == true && f.conn != nil {
		if err := f.conn.Close(); err != nil {
			return err
		}
		f.connected = false
	}
	return nil
}

func (f *Fetcher) ConnectIfNot() {
	if f.connected == false || f.conn == nil {
		if err := f.Connect(); err != nil {
			//TODO: LOG
			log.Printf("Fetch connection error: %v", err)
		}
	}
}

func (f *Fetcher) looper(fetchSignal chan struct{}, metricCh chan []byte) {
	if f.connected == false {
		f.Connect()
	}

	var (
		mb      bytes.Buffer
		metrics bytes.Buffer
	)

L:
	for {
		select {
		case <-f.ctx.Done():
			break L

		case <-fetchSignal:
			log.Println("Got fetchSignal")
			f.ConnectIfNot()
			if err := f.fetch(metricCh, &metrics, &mb); err != nil {
				log.Printf("Fetch error: %v", err)
			}
		}
	}

	if err := f.Disconnect(); err != nil {
		log.Printf("Fetcher disconnect: %v", err)
	}
}

// And fetch.go

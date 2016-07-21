package glia

import (
	gliametrics "github.com/anarcher/glia/lib/metrics"

	"golang.org/x/net/context"

	"net"
	"time"
)

type Sender struct {
	ctx     context.Context
	conn    net.Conn
	network string
	addr    string
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

func (s *Sender) Connect() (net.Conn, error) {
	c, err := net.Dial(s.network, s.addr)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Sender) Disconnect(c net.Conn) error {
	if c == nil {
		return nil
	}

	if err := c.Close(); err != nil {
		Logger.Log("sender", "disconnect", "err", err)
		return err
	}

	return nil
}

func (s *Sender) ConnectIfNot(c net.Conn, err error) (net.Conn, error) {
	if c != nil && err == nil {
		return c, err
	}

	if c != nil {
		if err := s.Disconnect(c); err != nil {
			return nil, err
		}
	}
	var cErr error
	if c, cErr = s.Connect(); cErr != nil {
		Logger.Log("sender", "connection", "err", err)
		return nil, err
	}

	return c, nil
}

func (s *Sender) looper(metricCh chan []byte) {
	var (
		conn net.Conn
		err  error
	)
L:
	for {
		select {
		case <-s.ctx.Done():
			break L
		case metrics := <-metricCh:
			gliametrics.Sending.Add(1)
			st := time.Now()

			var start, c int

			for {
				if conn, err = s.ConnectIfNot(conn, err); err == nil {
					if c, err = conn.Write(metrics[start:]); err != nil {
						gliametrics.SendErrorCount.Add(1)
						Logger.Log("sender", "write", "err", err)
					}
					start += c
					if c == 0 || start == len(metrics) {
						break
					}
				} else {
					gliametrics.SendErrorCount.Add(1)
					Logger.Log("sender", "write", "err", err)
					break
				}
			}

			gliametrics.SendLatency.Observe(time.Since(st))
			gliametrics.Sending.Add(-1)
		}
	}

	s.Disconnect(conn)
	WaitGroup.Done()
	Logger.Log("sender", "done")
}

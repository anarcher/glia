package glia

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/pusher/buddha/tcptest"
	"golang.org/x/net/context"
)

func Test_Sender(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	metrics := []byte("name.space 1.0 132131\n")
	ts := tcptest.NewServer(func(conn net.Conn) {
		defer conn.Close()

		bs, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(bs, metrics) {
			t.Error(fmt.Errorf("send and receive not match"))
		}
		wg.Done()
	})

	metricCh := make(chan []byte, 1)
	ctx := context.Background()

	NewSender(ctx, "tcp", ts.Addr.String(), metricCh)

	metricCh <- metrics

	wg.Wait()
}

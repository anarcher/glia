package glia

import (
	"bytes"
	"io/ioutil"
	"net"
	"sync"
	"testing"

	"github.com/pusher/buddha/tcptest"
	"golang.org/x/net/context"
)

func testFetch(t *testing.T, fixture []byte) {
	ts := tcptest.NewServer(func(conn net.Conn) {
		defer conn.Close()
		conn.Write(fixture)
	})
	defer ts.Close()

	ctx := context.Background()
	fetchCnt := 30000

	fetchSignal := make(chan struct{}, 1)
	metricCh := make(chan []byte, 0)

	NewFetcher(ctx, "tcp", ts.Addr.String(),
		fetchSignal,
		metricCh,
		fetchCnt,
		"test",
		false,
		FetchBufSize)

	fetchSignal <- struct{}{}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		metrics := <-metricCh
		metricLines := bytes.Split(metrics, []byte("\n"))
		for _, m := range metricLines {
			if len(m) <= 0 {
				continue
			}
			parts := bytes.Split(m, []byte(" "))
			if len(parts) != 3 {
				t.Errorf("metric format is wrong m:%v", string(m))
				t.Logf("%v", string(m))
			}

			//			t.Logf("%v", string(m))
		}

		wg.Done()
	}()

	wg.Wait()

}

func Test_Fetcher(t *testing.T) {
	samples := []string{"sample1.txt", "sample2.txt", "sample3.txt"}
	for _, filename := range samples {
		sample, _ := ioutil.ReadFile("../test/" + filename)
		t.Logf("Testing %v", filename)
		testFetch(t, sample)
	}
}

package glia

import (
	gliametrics "github.com/anarcher/glia/lib/metrics"

	"golang.org/x/net/context"

	"bytes"
	"fmt"
	"net"
	"time"
)

type Fetcher struct {
	ctx                  context.Context
	network              string
	addr                 string
	conn                 net.Conn
	connected            bool
	flushCnt             int
	graphitePrefix       string
	ignoreMetricOverTmax bool
}

func NewFetcher(ctx context.Context,
	network, addr string,
	fetchSignal chan struct{}, metricCh chan []byte,
	flushCnt int,
	graphitePrefix string,
	ignoreMetricOverTmax bool) *Fetcher {

	fetcher := &Fetcher{
		ctx:                  ctx,
		network:              network,
		addr:                 addr,
		flushCnt:             flushCnt,
		graphitePrefix:       graphitePrefix,
		ignoreMetricOverTmax: ignoreMetricOverTmax,
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
			Logger.Log("fetcher", "disconnect", "err", err)
			return err
		}
		f.connected = false
	}
	return nil
}

func (f *Fetcher) ConnectIfNot() error {
	if f.connected == false || f.conn == nil {
		if err := f.Connect(); err != nil {
			Logger.Log("fetcher", "connect", "err", err)
			return err
		}
	}
	return nil
}

func (f *Fetcher) looper(fetchSignal chan struct{}, metricCh chan []byte) {
	var (
		metrics bytes.Buffer
		mb      bytes.Buffer
	)

L:
	for {
		select {
		case <-f.ctx.Done():
			break L

		case <-fetchSignal:
			Logger.Log("fetch", "start")
			gliametrics.Fetching.Add(1)
			st := time.Now()
			if err := f.ConnectIfNot(); err == nil {
				if err := f.fetch(metricCh, &metrics, &mb); err != nil {
					gliametrics.FetchErrorCount.Add(1)
					Logger.Log("fetch", "err", "err", err)
				}
				Logger.Log("fetch", "done", "elapsed", fmt.Sprintf("%s", time.Since(st)))
			}
			f.Disconnect()
			gliametrics.FetchLatency.Observe(time.Since(st))
			gliametrics.Fetching.Add(-1)
		}
	}

	f.Disconnect()
	WaitGroup.Done()
	Logger.Log("fetcher", "done")
}

// And fetch.go

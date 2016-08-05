package glia

import (
	gliametrics "github.com/anarcher/glia/lib/metrics"

	"golang.org/x/net/context"

	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

const FetchBufSize = 2048

type Fetcher struct {
	ctx                  context.Context
	network              string
	addr                 string
	flushCnt             int
	graphitePrefix       string
	ignoreMetricOverTmax bool
	fetchBufSize         int
	fetchInterval        time.Duration
	clusterName          string
}

func NewFetcher(ctx context.Context,
	network, addr string,
	fetchSignal chan struct{}, metricCh chan []byte,
	flushCnt int,
	graphitePrefix string,
	ignoreMetricOverTmax bool,
	fetchBufSize int,
	fetchInterval time.Duration,
	clusterName string) *Fetcher {

	fetcher := &Fetcher{
		ctx:                  ctx,
		network:              network,
		addr:                 addr,
		flushCnt:             flushCnt,
		graphitePrefix:       graphitePrefix,
		ignoreMetricOverTmax: ignoreMetricOverTmax,
		fetchBufSize:         fetchBufSize,
		fetchInterval:        fetchInterval,
		clusterName:          clusterName,
	}
	if fetchBufSize <= 0 {
		fetcher.fetchBufSize = FetchBufSize
	}

	go fetcher.looper(fetchSignal, metricCh)

	return fetcher
}

func (f *Fetcher) Connect() (net.Conn, error) {
	c, err := net.Dial(f.network, f.addr)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (f *Fetcher) Disconnect(c net.Conn) error {
	if c == nil {
		return nil
	}

	if err := c.Close(); err != nil {
		Logger.Log("fetcher", "disconnect", "err", err)
		return err
	}

	return nil
}

func (f *Fetcher) ConnectIfNot(c net.Conn, err error) (net.Conn, error) {
	if c != nil && err == nil {
		return c, err
	}

	if c != nil {
		if err := f.Disconnect(c); err != nil {
			return nil, err
		}
	}
	var cErr error
	if c, cErr = f.Connect(); cErr != nil {
		Logger.Log("sender", "connection", "err", err)
		return nil, cErr
	}

	return c, nil
}

func (f *Fetcher) looper(fetchSignal chan struct{}, metricCh chan []byte) {
	var (
		metrics bytes.Buffer
		mb      bytes.Buffer
		conn    net.Conn
		connErr error
	)

L:
	for {
		select {
		case <-f.ctx.Done():
			break L

		case <-fetchSignal:
			gliametrics.Fetching.Add(1)
			st := time.Now()

			if conn, connErr = f.ConnectIfNot(conn, connErr); connErr == nil {
				if err := f.fetch(conn, metricCh, &metrics, &mb); err != nil {
					gliametrics.FetchErrorCount.Add(1)
					if err == io.EOF {
						connErr = err
					} else {
						Logger.Log("fetch", "err", "err", err)
					}
				}
				Logger.Log("fetch", "done", "elapsed", fmt.Sprintf("%s", time.Since(st)))
			}
			gliametrics.FetchLatency.Observe(time.Since(st))
			gliametrics.Fetching.Add(-1)
		}
	}

	f.Disconnect(conn)
	WaitGroup.Done()
	Logger.Log("fetcher", "done")
}

// And fetch.go

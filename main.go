package main

import (
	//"github.com/pusher/buddha/tcptest"

	"github.com/anarcher/glia/lib"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"gopkg.in/urfave/cli.v1"

	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

func MainAction(c *cli.Context) error {
	glia.Logger.Log("glia", "start", "version", Version, "gitcommit", GitCommit)
	ctx, cancelFunc := context.WithCancel(context.Background())
	glia.Shutdown(cancelFunc)

	//metric handler
	http.Handle("/metrics", stdprometheus.Handler())
	//pprof handler
	http.Handle("/debug/pprof", http.HandlerFunc(pprof.Index))

	go func() {
		glia.Logger.Log("err", http.ListenAndServe(c.String("metric_addr"), nil))
	}()

	fetchSignal := make(chan struct{})
	metricCh := make(chan []byte)

	fetchInterval, err := time.ParseDuration(c.String("fetch_interval"))
	if err != nil {
		glia.Logger.Log("err", err)
		cancelFunc()
		return err
	}

	var fetchers []*glia.Fetcher
	{
		network := c.String("gmond_network")
		addr := c.String("gmond")
		bufSize := c.Int("buffer_size")
		graphitePrefix := c.String("graphite_prefix")
		ignoreMetricOverTmax := c.Bool("ignore_metric_over_tmax")
		fetch_buf_size := c.Int("fetch_buf_size")
		for i := 0; i < c.Int("fetcher"); i++ {
			f := glia.NewFetcher(ctx, network, addr, fetchSignal, metricCh, bufSize, graphitePrefix, ignoreMetricOverTmax, fetch_buf_size, fetchInterval)
			fetchers = append(fetchers, f)
			glia.WaitGroup.Add(1)
		}
	}

	var senders []*glia.Sender
	{
		network := c.String("graphite_network")
		addr := c.String("graphite")
		for i := 0; i < c.Int("sender"); i++ {
			s := glia.NewSender(ctx, network, addr, metricCh)
			senders = append(senders, s)
			glia.WaitGroup.Add(1)
		}
	}

	{
		go func() {
			tick := time.Tick(fetchInterval)
		L:
			for {
				select {
				case <-ctx.Done():
					break L
				case <-tick:
					fetchSignal <- struct{}{}
				}
			}
		}()

		fetchSignal <- struct{}{} // Start on
	}

	glia.WaitGroup.Wait()
	glia.Logger.Log("glia", "end")

	return nil
}

func main() {

	app := cli.NewApp()
	app.Name = "glia"
	app.Usage = "It comes between Gmond and Graphite"
	app.Action = MainAction
	app.Version = fmt.Sprintf("%s (%s)", Version, GitCommit)
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "fetcher,f",
			Value:  10,
			Usage:  "The number of fetchers",
			EnvVar: "FETCHER",
		},
		cli.IntFlag{
			Name:   "sender,s",
			Value:  10,
			Usage:  "The number of senders",
			EnvVar: "SENDER",
		},
		cli.StringFlag{
			Name:   "fetch_interval,i",
			Value:  "10s",
			Usage:  "The duration of to fetch  from gmond",
			EnvVar: "FETCH_INTERVAL",
		},
		cli.IntFlag{
			Name:   "buffer_size,b",
			Value:  1000,
			Usage:  "The buffer size of sending (the number of metric line)",
			EnvVar: "BUFFER_SIZE",
		},
		cli.StringFlag{
			Name:   "gmond,g",
			Value:  "localhost:8649",
			Usage:  "The address of gmond",
			EnvVar: "GMOND_ADDR",
		},
		cli.StringFlag{
			Name:   "gmond_network",
			Value:  "tcp",
			Usage:  "The network of gmond",
			EnvVar: "GMOND_NETWORK",
		},
		cli.StringFlag{
			Name:   "graphite,c",
			Value:  "localhost:2013",
			Usage:  "The address of graphtie carbon",
			EnvVar: "GRAPHITE_ADDR",
		},
		cli.StringFlag{
			Name:   "graphite_network",
			Value:  "udp",
			Usage:  "The network of graphite carbon",
			EnvVar: "GRAPHITE_NETWORK",
		},
		cli.StringFlag{
			Name:   "graphite_prefix,p",
			Value:  "ganglia",
			Usage:  "The prefix to prepend to the metric names exported",
			EnvVar: "GRAPHITE_PREFIX",
		},
		cli.BoolTFlag{
			Name:   "ignore_metric_over_tmax",
			Usage:  "enable or disable to ignore metric over tmax",
			EnvVar: "IGNORE_METRIC_OVER_TMAX",
		},
		cli.IntFlag{
			Name:   "fetch_buf_size",
			Usage:  "flusing fetch buffer size (sending packet size)",
			Value:  1024,
			EnvVar: "FETCH_BUF_SIZE",
		},
		cli.StringFlag{
			Name:   "metric_addr",
			Value:  ":8002",
			Usage:  "The Prometheus metrics export addr",
			EnvVar: "METRIC_ADDR",
		},
	}
	app.Run(os.Args)

}

package main

import (
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"

	"fmt"
	"os"
	"time"
)

func MainAction(c *cli.Context) {
	Logger.Log("glia", "start", "version", Version)
	ctx, cancelFunc := context.WithCancel(context.Background())
	Shutdown(cancelFunc)

	fetchSignal := make(chan struct{})
	metricCh := make(chan []byte)

	var fetchers []*Fetcher
	{
		network := c.String("gmond_network")
		addr := c.String("gmond")
		bufSize := c.Int("buffer_size")
		for i := 1; i < c.Int("fetcher"); i++ {
			f := NewFetcher(ctx, network, addr, fetchSignal, metricCh, bufSize)
			fetchers = append(fetchers, f)
			WaitGroup.Add(1)
		}
	}

	var senders []*Sender
	{
		network := c.String("graphtie_network")
		addr := c.String("graphite")
		for i := 1; i < c.Int("fetcher"); i++ {
			s := NewSender(ctx, network, addr, metricCh)
			senders = append(senders, s)
			WaitGroup.Add(1)
		}
	}

	{
		interval, err := time.ParseDuration(c.String("fetch_interval"))
		if err != nil {
			Logger.Log("err", err)
			cancelFunc()
			return
		}

		go func() {
			tick := time.Tick(interval)
		L:
			for {
				select {
				case <-ctx.Done():
					break L
				case <-tick:
					Logger.Log("fetch", "event", "fire", true)
					fetchSignal <- struct{}{}
				}
			}
		}()
	}

	WaitGroup.Wait()
	Logger.Log("glia", "end")
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
			Value:  "1m",
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
			Name:   "graphtie_network",
			Value:  "tcp",
			Usage:  "The network of graphite carbon",
			EnvVar: "GRAPHITE_NETWORK",
		},
	}
	app.Run(os.Args)

}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var (
	n         = flag.Int("n", 1, "number of workers")
	i         = flag.Int("i", 10, "number of query executions")
	t         = flag.Int("t", 10000, "client timeout (in ms)")
	f         = flag.String("f", "obs-queries.csv", "queries will be read from this file")
	serverURL = flag.String("url", "http://localhost:9201", "Promscale server URL")
	cpus      = flag.Int("cpus", runtime.NumCPU(), "number of CPUs to use for go runtime")
)

func main() {
	ctx := context.Background()
	flag.Parse()

	fmt.Println("psb - benchmarking Promscale")

	queries, err := PrepareQueries(*f)
	if err != nil {
		log.Printf("could not build queries. err: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d queries in input file: %s\n", len(queries), *f)

	w := Work{
		numWorkers:    *n,
		numIterations: *i,
		cpus:          *cpus,
		serverURL:     *serverURL,
		clientTimeout: *t,
	}
	w.PrintConf()

	if res := w.Health(ctx); res.isError {
		log.Println("could not confirm health of given Promscale server, is it running?")
		log.Println("quitting")
		os.Exit(1)
	}
	log.Println("Promscale health OK, moving to benchmark")

	// handle CtrlC
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		log.Println("Ctrl-C handler registerd")
		<-c
		log.Println("stopping now")
		w.Stop()
	}()

	time.Sleep(1 * time.Second)
	log.Println("ready to benchmark")
	for i, q := range queries {
		if !w.stopRequested {
			fmt.Println("#", i, "query:", q.queryStr)
			w.Init(q)
			w.Run(ctx)
			w.PrintReport()
		}
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

type Work struct {
	// config
	serverURL     string
	numWorkers    int
	numIterations int
	cpus          int
	clientTimeout int // in milliseconds
	query         query

	// orchestration
	stopCh        chan struct{}
	stopRequested bool
	reportDoneCh  chan struct{}

	// reporting
	resultsCh chan result
	report    report
}

func (w *Work) PrintConf() {
	fmt.Println()
	fmt.Println("benchmark config")
	fmt.Println("server url:", w.serverURL)
	fmt.Println("# concurrent workers/clients:", w.numWorkers)
	fmt.Println("# iterations of each query:", w.numIterations)
	fmt.Println("# CPUs to use:", w.cpus)
	fmt.Println("# client timeout:", time.Duration(w.clientTimeout*int(time.Millisecond)))
	fmt.Println()
}

func (w *Work) Init(q query) {
	w.query = q
	w.stopCh = make(chan struct{}, w.numWorkers)
	w.reportDoneCh = make(chan struct{}, 1)

	w.resultsCh = make(chan result)
	w.report = newReport(w.query)

	runtime.GOMAXPROCS(w.cpus)
}

func (w *Work) Run(ctx context.Context) {
	// start looking for results
	go w.ResultListener()

	var wg sync.WaitGroup
	wg.Add(w.numWorkers)

	for i := 0; i < w.numWorkers; i++ {
		go func(i int, q query) {
			defer wg.Done()
			for iter := 0; iter < w.numIterations; iter++ {
				select {
				case <-w.stopCh:
					log.Println("stop signal received to", i)
					return
				default:
					// execute query
					w.resultsCh <- w.Query(ctx)
				}
			}
		}(i, w.query)
	}

	wg.Wait()
	log.Println("wait over, work is done, no more pending results")
	close(w.resultsCh)
}

func (w *Work) Stop() {
	w.stopRequested = true
	log.Println("signalling all workers to stop")
	for i := 0; i < w.numWorkers; i++ {
		w.stopCh <- struct{}{}
	}
}

func (w *Work) ResultListener() {
	log.Println("waiting for all results")

	for v := range w.resultsCh {
		w.report.results = append(w.report.results, v)
	}

	w.reportDoneCh <- struct{}{}
}

func (w *Work) PrintReport() {
	<-w.reportDoneCh
	log.Println("report is ready:", len(w.report.results))
	w.report.Print()
}

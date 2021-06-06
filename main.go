package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Work struct {
	numWorkers    int
	numIterations int
	query         string

	startTime time.Time
	endTime   time.Time

	resultsCh     chan int
	stopCh        chan struct{}
	stopRequested bool
	reportDoneCh  chan struct{}

	stats []int
}

func (w *Work) Init(query string) {
	w.query = query
	w.stopCh = make(chan struct{}, w.numWorkers)
	w.resultsCh = make(chan int)
	w.reportDoneCh = make(chan struct{}, 1)
	w.stats = []int{}
}

func (w *Work) Run() {
	w.startTime = time.Now()
	defer func() {
		w.endTime = time.Now()
	}()

	// start looking for results
	go w.ResultListener()

	var wg sync.WaitGroup
	wg.Add(w.numWorkers)

	for i := 0; i < w.numWorkers; i++ {
		go func(i int, query string) {
			defer wg.Done()
			for iter := 0; iter < w.numIterations; iter++ {
				select {
				case <-w.stopCh:
					log.Println("stop signal received to", i)
					return
				default:
					log.Println(i, "making query:", query, "iteration:", iter)
					time.Sleep(1 * time.Second)
					w.resultsCh <- int(time.Now().Unix())
				}
			}
		}(i, w.query)
	}

	wg.Wait()
	log.Println("wait over, work is done, no more results")
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
		w.stats = append(w.stats, v)
	}

	w.reportDoneCh <- struct{}{}
}

func (w *Work) PrintReport() {
	<-w.reportDoneCh
	log.Println("time spent:", w.endTime.Sub(w.startTime))
	log.Println("report is ready:", len(w.stats))
}

func main() {
	fmt.Println("psb")

	w := Work{
		numWorkers:    2,
		numIterations: 3,
	}

	// handle CtrlC
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		log.Println("Ctrl-C handler registerd")
		<-c
		log.Println("stopping now")
		w.Stop()
	}()

	queries := []string{"aa", "bb"}
	for _, q := range queries {
		log.Println("processing query:", q)
		if !w.stopRequested {
			w.Init(q)
			w.Run()
			w.PrintReport()
		}
	}
}

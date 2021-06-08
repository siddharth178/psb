package main

import (
	"fmt"
	"sort"
	"time"
)

type report struct {
	query query

	// captured results
	results []result

	// derived numbers
	numberOfQueries int
	failedQueries   int
	totalTime       time.Duration
	totalExecTime   time.Duration

	maxTime time.Duration
	minTime time.Duration
	avgTime time.Duration
	medTime time.Duration
}

func newReport(q query) report {
	r := report{}
	r.results = []result{}
	r.query = q

	return r
}

type result struct {
	st time.Time
	et time.Time

	statusCode int
	err        error
	isError    bool
}

func (r *report) Prepare() {
	validResults := []result{}
	r.numberOfQueries = len(r.results)
	for _, res := range r.results {
		if res.isError {
			r.failedQueries++
		} else {
			validResults = append(validResults, res)
		}
	}
	goodQueries := r.numberOfQueries - r.failedQueries
	if goodQueries <= 0 {
		return
	}

	minSt := validResults[0].st
	maxEt := validResults[0].et

	execTime := validResults[0].et.Sub(validResults[0].st)
	r.minTime = execTime
	r.maxTime = execTime
	r.totalExecTime = execTime

	// for median
	execTimes := []time.Duration{}
	execTimes = append(execTimes, execTime)

	for i := 1; i < len(validResults); i++ {
		if validResults[i].st.Before(minSt) {
			minSt = validResults[i].st
		}
		if validResults[i].et.After(maxEt) {
			maxEt = validResults[i].et
		}

		execTime := validResults[i].et.Sub(validResults[i].st)
		if r.minTime > execTime {
			r.minTime = execTime
		}
		if r.maxTime < execTime {
			r.maxTime = execTime
		}
		r.totalExecTime += execTime             // for avg
		execTimes = append(execTimes, execTime) // for median
	}
	r.totalTime = maxEt.Sub(minSt)
	r.avgTime = time.Duration(int64(r.totalExecTime) / int64(goodQueries))

	sort.Slice(execTimes, func(i, j int) bool {
		return int64(execTimes[i]) < int64(execTimes[j])
	})
	if len(execTimes)%2 == 0 {
		r.medTime = (execTimes[(len(execTimes)/2)-1] + execTimes[len(execTimes)/2]) / 2
	} else {
		r.medTime = execTimes[len(execTimes)/2]
	}
}

func (r *report) Print() {
	r.Prepare()

	fmt.Println("--------")
	fmt.Println("query str:", r.query.queryStr)
	fmt.Println("#queries:", r.numberOfQueries)
	fmt.Println("#failed queries:", r.failedQueries)
	fmt.Println("total time:", r.totalTime)
	fmt.Println("min time:", r.minTime)
	fmt.Println("max time:", r.maxTime)
	fmt.Println("avg time:", r.avgTime)
	fmt.Println("med time:", r.medTime)
	fmt.Println("--------")
	fmt.Println()
}

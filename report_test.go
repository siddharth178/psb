package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	goodResults = []result{
		{
			// 2 sec
			st: time.Date(2021, 6, 7, 0, 0, 0, 0, time.UTC),
			et: time.Date(2021, 6, 7, 0, 0, 2, 0, time.UTC),
		},
		{
			// 4 sec
			st: time.Date(2021, 6, 7, 0, 0, 1, 0, time.UTC),
			et: time.Date(2021, 6, 7, 0, 0, 5, 0, time.UTC),
		},
		{
			// 5 sec
			st: time.Date(2021, 6, 7, 0, 0, 1, 0, time.UTC),
			et: time.Date(2021, 6, 7, 0, 0, 6, 0, time.UTC),
		},
		{
			// 1 sec
			st: time.Date(2021, 6, 7, 0, 0, 0, 0, time.UTC),
			et: time.Date(2021, 6, 7, 0, 0, 1, 0, time.UTC),
		},
	}

	badResults = []result{
		{
			// 10 sec failure
			isError: true,
			st:      time.Date(2021, 6, 7, 0, 0, 0, 0, time.UTC),
			et:      time.Date(2021, 6, 7, 0, 0, 10, 0, time.UTC),
		},
	}
)

func TestPrepare_Allgood(t *testing.T) {
	r := report{
		results: goodResults,
	}
	r.Prepare()
	assert.Equal(t, 4, r.numberOfQueries)
	assert.Equal(t, 0, r.failedQueries)
	assert.Equal(t, float64(6), r.totalTime.Seconds())
	assert.Equal(t, float64(12), r.totalExecTime.Seconds())
	assert.Equal(t, float64(1), r.minTime.Seconds())
	assert.Equal(t, float64(5), r.maxTime.Seconds())
	assert.Equal(t, float64(3), r.avgTime.Seconds())
	assert.Equal(t, float64(3), r.medTime.Seconds())
}

func TestPrepare_CoupleFailures(t *testing.T) {
	r := report{
		results: append(goodResults, badResults...),
	}
	r.Prepare()
	assert.Equal(t, 5, r.numberOfQueries)
	assert.Equal(t, 1, r.failedQueries)
	assert.Equal(t, float64(6), r.totalTime.Seconds())
	assert.Equal(t, float64(12), r.totalExecTime.Seconds())
	assert.Equal(t, float64(1), r.minTime.Seconds())
	assert.Equal(t, float64(5), r.maxTime.Seconds())
	assert.Equal(t, float64(3), r.avgTime.Seconds())
	assert.Equal(t, float64(3), r.medTime.Seconds())
}

func TestPrepare_AllFailures(t *testing.T) {
	r := report{
		results: badResults,
	}
	r.Prepare()
	assert.Equal(t, 1, r.numberOfQueries)
	assert.Equal(t, 1, r.failedQueries)
	assert.Equal(t, float64(0), r.totalTime.Seconds())
	assert.Equal(t, float64(0), r.totalExecTime.Seconds())
	assert.Equal(t, float64(0), r.minTime.Seconds())
	assert.Equal(t, float64(0), r.maxTime.Seconds())
	assert.Equal(t, float64(0), r.avgTime.Seconds())
	assert.Equal(t, float64(0), r.medTime.Seconds())
}

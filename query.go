package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type query struct {
	queryStr string
	params   map[string]string
}

func newQuery(queryStr string) (*query, error) {
	toks := strings.Split(queryStr, "|")
	if len(toks) == 4 {
		// query range

		st := toks[1]
		if len(st) > 10 {
			st = st[:10] + "." + st[10:]
		}
		et := toks[2]
		if len(et) > 10 {
			et = et[:10] + "." + et[10:]
		}

		m := map[string]string{
			"endpoint": "query_range",
			"query":    toks[0],
			"start":    st,
			"end":      et,
			"step":     toks[3],
		}
		return &query{
			queryStr: queryStr,
			params:   m,
		}, nil
	} else {
		return nil, errors.New("unknown query format")
	}
}

func PrepareQueries(fileName string) ([]query, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queries := []query{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		queryStr := strings.TrimSpace(s.Text())
		q, err := newQuery(queryStr)
		if err != nil {
			log.Println("skipping query with incorrect format:", queryStr)
			continue
		}
		queries = append(queries, *q)
	}

	return queries, nil
}

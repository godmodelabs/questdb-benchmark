package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	qdb "github.com/questdb/go-questdb-client/v3"
)

// const QUESTDB_API_URL = "http://questdb-benchmark:9000/exec"
const QUESTDB_API_URL = "http://localhost:9000/exec"
const QUESTDB_HOST = "localhost:9000"
const TABLE_NAME = "benchmark"

type QueryResponse struct {
	Query   string `json:"query"`
	Columns []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"columns"`
	Timestamp int     `json:"timestamp"`
	Dataset   [][]any `json:"dataset"`
	Count     int     `json:"count"`
}

func Query(query string) (r string, err error) {

	base, err := url.Parse(QUESTDB_API_URL)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("query", query)
	base.RawQuery = params.Encode()

	qdbClient := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, base.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := qdbClient.Do(req)
	if err != nil {
		return "", err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func main() {

	_, err := Query("DROP TABLE " + TABLE_NAME)
	if err != nil {
		fmt.Printf("error executing DROP TABLE query %v\n", err)
		os.Exit(1)
	}

	_, err = Query("CREATE TABLE " + TABLE_NAME + " (timestamp TIMESTAMP, tick INT, instrumentID SYMBOL, price DOUBLE, volume INT) TIMESTAMP(timestamp) PARTITION BY DAY WAL;")
	if err != nil {
		fmt.Printf("error executing CREATE TABLE query %v\n", err)
		os.Exit(1)
	}

	// setup qdb client
	ctx := context.TODO()
	qdbClient, err := qdb.NewLineSender(ctx,
		qdb.WithHttp(),
		qdb.WithAddress(QUESTDB_HOST),
		qdb.WithBasicAuth("admin", "quest"),
		//qdb.WithAutoFlushRows(1000),
		qdb.WithRequestTimeout(20*time.Second),
		qdb.WithRetryTimeout(20*time.Second))
	if err != nil {
		fmt.Printf("error when setting up qdb client %v\n", err)
		os.Exit(1)
	}
	defer qdbClient.Close(ctx)

	g := NewGenerator(20000000, 5000, 50000, time.Now(), time.Hour)
	g.Start()

	statsTicker := *time.NewTicker(1 * time.Second)
	// run statistics output
	go func() {
		var prevTicks int
		for range statsTicker.C {
			r, err := Query("wal_tables()")
			if err != nil {
				fmt.Printf("error executing wal_tables() query %v\n", err)
				os.Exit(1)
			}

			dec := json.NewDecoder(strings.NewReader(r))
			w := QueryResponse{}
			err = dec.Decode(&w)
			if err != nil {
				fmt.Printf("error decoding query json%v\n", err)
				os.Exit(1)
			}

			writerTxn, sequencerTxn := w.Dataset[0][2].(float64), w.Dataset[0][4].(float64)

			r, err = Query("SELECT count() FROM " + TABLE_NAME)
			if err != nil {
				fmt.Printf("error executing select count query %v\n", err)
				os.Exit(1)
			}
			dec = json.NewDecoder(strings.NewReader(r))
			err = dec.Decode(&w)
			if err != nil {
				fmt.Printf("error decoding query json%v\n", err)
				os.Exit(1)
			}
			count := int(w.Dataset[0][0].(float64))

			generatedTicks := g.GeneratedTicks()
			ticksPerSec := generatedTicks - prevTicks
			tickDelta := generatedTicks - count
			prevTicks = generatedTicks
			fmt.Printf("generatedTicks %v - ticks/second %v - questdb count %v - writerTxn %v readerTxn %v w/r delta %v ticksDelta %v\n",
				generatedTicks, ticksPerSec, count, writerTxn, sequencerTxn, sequencerTxn-writerTxn, tickDelta)
		}
	}()

	for tick := range g.TickChannel {
		err = qdbClient.Table(TABLE_NAME).
			Symbol("instrumentID", strconv.Itoa(tick.id)).
			Int64Column("tick", int64(tick.tick)).
			Float64Column("price", tick.price).
			Int64Column("volume", tick.volume).
			At(ctx, tick.ts)
		if err != nil {
			fmt.Printf("error when adding ILP %v\n", err)
			os.Exit(1)
		}
	}
	qdbClient.Flush(ctx)
}

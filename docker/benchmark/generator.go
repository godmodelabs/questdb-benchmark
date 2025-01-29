package main

import (
	"math/rand/v2"
	"sync/atomic"
	"time"
)

type Generator struct {
	generatedTicks atomic.Uint64
	totalTicks     int
	totalSymbols   int
	ticksPerSecond int
	TickChannel    chan Tick
	doneChannel    chan bool
	ticker         time.Ticker
	tickMap        map[int]int
	delay          time.Duration
}

type Tick struct {
	ts     time.Time
	id     int
	tick   int
	price  float64
	volume int64
}

func NewGenerator(totalTicks int, totalSymbols int, ticksPerSecond int, start time.Time, delay time.Duration) *Generator {
	return &Generator{
		totalTicks:     totalTicks,
		totalSymbols:   totalSymbols,
		ticksPerSecond: ticksPerSecond,
		TickChannel:    make(chan Tick),
		tickMap:        make(map[int]int, totalSymbols),
		ticker:         *time.NewTicker(1 * time.Second),
		doneChannel:    make(chan bool),
		delay:          delay,
	}
}

func (g *Generator) Start() {

	go func() {
		defer close(g.TickChannel)

		for {
			select {
			case ts := <-g.ticker.C:
				if int(g.generatedTicks.Load()) >= g.totalTicks {
					return
				}

				if g.delay > 0 {
					ts = ts.Add(-g.delay)
				}

				for i := 0; i < g.ticksPerSecond && int(g.generatedTicks.Load()) < g.totalTicks; i++ {
					id := rand.IntN(g.totalSymbols) + 1

					_, exists := g.tickMap[id]
					if exists {
						g.tickMap[id]++
					} else {
						g.tickMap[id] = 1
					}

					g.TickChannel <- Tick{
						ts:     ts,
						id:     id,
						tick:   g.tickMap[id],
						price:  rand.Float64() * 1000,
						volume: int64(rand.IntN(10000)),
					}

					g.generatedTicks.Add(1)
				}
			case <-g.doneChannel:
				return
			}
		}
	}()
}

func (g *Generator) GeneratedTicks() int {
	return int(g.generatedTicks.Load())
}

func (g *Generator) Stop() {
	g.doneChannel <- true
}

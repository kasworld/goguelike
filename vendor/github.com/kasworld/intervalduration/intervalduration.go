// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)

// Package intervalduration manage actpersec, dur jitter, act jitter
package intervalduration

import (
	"fmt"
	"time"

	"github.com/kasworld/durjitter"
)

type IntervalDuration struct {
	interval  *durjitter.DurJitter
	duration  *durjitter.DurJitter
	lastEnded *Act
}

func New(name string) *IntervalDuration {
	indu := &IntervalDuration{
		interval: durjitter.New(name + " interval"),
		duration: durjitter.New(name + " duration"),
	}
	return indu
}

func (indu *IntervalDuration) String() string {
	return fmt.Sprintf("%d%v%v",
		indu.GetCount(), indu.interval, indu.duration,
	)
}

func (indu *IntervalDuration) GetInterval() *durjitter.DurJitter {
	return indu.interval
}
func (indu *IntervalDuration) GetDuration() *durjitter.DurJitter {
	return indu.duration
}

func (indu *IntervalDuration) GetCount() int {
	return indu.interval.GetCount()
}

func (indu *IntervalDuration) BeginAct() *Act {
	act := &Act{
		startTime: time.Now(),
		endFn:     indu.endAct,
	}
	return act
}

func (indu *IntervalDuration) endAct(act *Act) {
	dur := time.Now().Sub(act.startTime)
	indu.duration.Add(dur)

	if indu.lastEnded != nil {
		interval := act.startTime.Sub(indu.lastEnded.startTime)
		indu.interval.Add(interval)
	}

	indu.lastEnded = act
}

type Act struct {
	endFn     func(*Act)
	startTime time.Time `prettystring:"simple"`
}

func (act *Act) End() {
	act.endFn(act)
}

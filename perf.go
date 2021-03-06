// Copyright © 2013 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Quick and Easy Performance Analyzer
// Useful for comparing A/B against different drafts of functions or different functions
// Loosely inspired by the go benchmark package
//
// Example:
//	import "github.com/spf13/nitro"
//	timer := nitro.Start()
//	prepTemplates()
//	timer.Step("initialize & template prep")
//	CreatePages()
//	timer.Stop("import pages")
package nitro

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

// Used for every benchmark for measuring memory.
var memStats runtime.MemStats

var AnalysisOn = false

var condition func() bool
var writer io.Writer

func init() {
	flag.BoolVar(&AnalysisOn, "stepAnalysis", false, "display memory and timing of different steps of the program")
	condition = func() bool { return true }
	writer = os.Stdout
}

type B struct {
	initial  time.Time
	start    time.Time // Time step started
	duration time.Duration
	timerOn  bool
	// The initial states of memStats.Mallocs and memStats.TotalAlloc.
	startAllocs uint64
	startBytes  uint64
	// The net total of this test after being run.
	netAllocs uint64
	netBytes  uint64
}

func (b *B) startTimer() {
	if b == nil {
		fmt.Println("ERROR: can't call startTimer on a nil value")
		os.Exit(-1)
	}
	if !b.timerOn {
		runtime.ReadMemStats(&memStats)
		b.startAllocs = memStats.Mallocs
		b.startBytes = memStats.TotalAlloc
		b.start = time.Now()
		b.timerOn = true
	}
}

func (b *B) stopTimer() {
	if b == nil {
		fmt.Println("ERROR: can't call stopTimer on a nil value")
		os.Exit(-1)
	}
	if b.timerOn {
		b.duration += time.Since(b.start)
		runtime.ReadMemStats(&memStats)
		b.netAllocs += memStats.Mallocs - b.startAllocs
		b.netBytes += memStats.TotalAlloc - b.startBytes
		b.timerOn = false
	}
}

// ResetTimer sets the elapsed benchmark time to zero.
// It does not affect whether the timer is running.
func (b *B) resetTimer() {
	if b.timerOn {
		runtime.ReadMemStats(&memStats)
		b.startAllocs = memStats.Mallocs
		b.startBytes = memStats.TotalAlloc
		b.start = time.Now()
	}
	b.duration = 0
	b.netAllocs = 0
	b.netBytes = 0
}

func SetCondition(c func() bool) {
	condition = c
}

func SetWriter(w io.Writer) {
	writer = w
}

func Start(title string) *B {
	if !AnalysisOn {
		return nil
	}
	if !condition() {
		return nil
	}
	b := &B{}
	b.initial = time.Now()
	b.writeHeader(title)
	b.resetTimer()
	b.startTimer()
	return b
}

// Call perf.Step("step name") at each step in your
// application you want to benchmark
// Measures time spent since last Step call.
func (b *B) Step(str string) {
	if !AnalysisOn {
		return
	}
	if b == nil {
		return
	}

	b.stopTimer()
	b.write(str)

	b.resetTimer()
	b.startTimer()
}

func (b *B) Stop(str string) {
	if !AnalysisOn {
		return
	}
	if b == nil {
		return
	}

	b.stopTimer()
	b.write(str)
}

func (b *B) writeHeader(str string) {
	fmt.Fprintf(writer, "%9s\t%9s\t%9s\t%9s\t%s\n", "during", "total", "memBytes", "memAllocs", str)
}

func (b *B) write(str string) {
	r := b.results()
	fmt.Fprintf(writer, "%9d\t%9d\t%9d\t%9v\t%s\n", r.T*time.Nanosecond, r.C*time.Nanosecond, r.MemBytes, r.MemAllocs, str)
}

func (b *B) results() R {
	return R{time.Since(b.initial), b.duration, b.netAllocs, b.netBytes}
}

type R struct {
	C         time.Duration // Cumulative time taken
	T         time.Duration // The total time taken.
	MemAllocs uint64        // The total number of memory allocations.
	MemBytes  uint64        // The total number of bytes allocated.
}

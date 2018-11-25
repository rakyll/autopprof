// Copyright 2018 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package autopprof provides a development-time
// library to collect pprof profiles from Go programs.
//
// This package is experimental and APIs may change.
package autopprof

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

// Profile represents a pprof profile.
type Profile interface {
	Capture() (profile string, err error)
}

// CPUProfile captures the CPU profile.
type CPUProfile struct {
	Duration time.Duration // 30 seconds by default
}

// TODO(jbd): Add CPU profiling rate as an option.

func (p CPUProfile) Capture() (string, error) {
	dur := p.Duration
	if dur == 0 {
		dur = 30 * time.Second
	}

	f := newTemp()
	if err := pprof.StartCPUProfile(f); err != nil {
		return "", nil
	}
	time.Sleep(dur)
	pprof.StopCPUProfile()
	if err := f.Close(); err != nil {
		return "", nil
	}
	return f.Name(), nil
}

// HeapProfile captures the heap profile.
type HeapProfile struct{}

func (p HeapProfile) Capture() (string, error) {
	return captureProfile("heap")
}

// MutexProfile captures stack traces of holders of contended mutexes.
type MutexProfile struct{}

func (p MutexProfile) Capture() (string, error) {
	return captureProfile("mutex")
}

// BlockProfile captures stack traces that led to blocking on synchronization primitives.
type BlockProfile struct {
	// Rate is the fraction of goroutine blocking events that
	// are reported in the blocking profile. The profiler aims to
	// sample an average of one blocking event per rate nanoseconds spent blocked.
	//
	// If zero value is provided, it will include every blocking event
	// in the profile.
	Rate int
}

func (p BlockProfile) Capture() (string, error) {
	if p.Rate > 0 {
		runtime.SetBlockProfileRate(p.Rate)
	}
	return captureProfile("block")
}

// GoroutineProfile captures stack traces of all current goroutines.
type GoroutineProfile struct{}

func (p GoroutineProfile) Capture() (string, error) {
	return captureProfile("goroutine")
}

// Threadcreate profile captures the stack traces that led to the creation of new OS threads.
type ThreadcreateProfile struct{}

func (p ThreadcreateProfile) Capture() (string, error) {
	return captureProfile("threadcreate")
}

func captureProfile(name string) (string, error) {
	f := newTemp()
	if err := pprof.Lookup(name).WriteTo(f, 2); err != nil {
		return "", nil
	}
	if err := f.Close(); err != nil {
		return "", nil
	}
	return f.Name(), nil
}

// TODO(jbd): Add support for custom.

// Capture captures the given profiles at SIGINT
// and opens a browser with the collected profiles.
//
// Capture should be used in development-time
// and shouldn't be in production binaries.
func Capture(p Profile) {
	// TODO(jbd): As a library, we shouldn't be in the
	// business of signal handling. Provide a better way
	// trigger the capture.
	go capture(p)
}

func capture(p Profile) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT) // TODO(jbd): Add windows support.

	fmt.Println("Send SIGQUIT (CTRL+\\) to the process to capture...")

	for {
		<-c
		log.Println("Starting to capture.")

		profile, err := p.Capture()
		if err != nil {
			log.Printf("Cannot capture profile: %v", err)
		}

		// Open profile with pprof.
		log.Printf("Starting go tool pprof %v", profile)
		cmd := exec.Command("go", "tool", "pprof", "-http=:", profile)
		if err := cmd.Run(); err != nil {
			log.Printf("Cannot start pprof UI: %v", err)
		}
	}
}

func newTemp() (f *os.File) {
	f, err := ioutil.TempFile("", "profile-")
	if err != nil {
		log.Fatalf("Cannot create new temp profile file: %v", err)
	}
	return f
}

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
	f := newTemp()
	if err := pprof.WriteHeapProfile(f); err != nil {
		return "", nil
	}
	if err := f.Close(); err != nil {
		return "", nil
	}
	return f.Name(), nil
}

// TODO(jbd): Add all supported profiles.

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

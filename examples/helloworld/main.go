package main

import (
	"time"

	"github.com/rakyll/autopprof"
)

func main() {
	autopprof.Capture(autopprof.HeapProfile{})

	time.Sleep(time.Hour)
}

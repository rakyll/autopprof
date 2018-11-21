package main

import (
	"strconv"
	"time"

	"github.com/rakyll/autopprof"
)

func main() {
	autopprof.Capture(autopprof.CPUProfile{
		Duration: 15 * time.Second,
	})

	for i := 0; i < 5000; i++ {
		generateID(5, 1000)
	}
}

func generateID(duration int, usage int) {
	for j := 0; j < duration; j++ {
		go func() {
			for i := 0; i < usage*80000; i++ {
				str := "str" + strconv.Itoa(i)
				str = str + "a"
			}
		}()
		time.Sleep(1 * time.Second)
	}
}

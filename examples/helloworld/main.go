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

// Package main contains a simple hello world example for autopprof.
package main

import (
	"strconv"
	"time"

	"github.com/rakyll/autopprof"
)

func main() {
	autopprof.Capture(autopprof.CPUProfile{
		Duration: 5 * time.Second,
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

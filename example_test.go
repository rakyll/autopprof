// Copyright 2014 Google Inc. All Rights Reserved.
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

package autopprof_test

import (
	"time"

	"github.com/rakyll/autopprof"
)

func Example() {
	// Add the following to your main, then
	// use CTRL+\ to intercept and capture.
	// Pprof UI will start in 15 seconds once
	// the profile is captured.
	autopprof.Capture(autopprof.CPUProfile{
		Duration: 30 * time.Second,
	})
}

# autopprof

[![GoDoc](https://godoc.org/github.com/rakyll/autopprof?status.svg)](https://godoc.org/github.com/rakyll/autopprof)

Pprof made easy at development time.

## Guide

Add autopprof.Capture to your main function.

```go
import "github.com/rakyll/autopprof"

autopprof.Capture(autopprof.CPUProfile{
    Duration: 15 * time.Second,
})
```

Run your program and send SIGQUIT to the process
(or CTRL+\\ on Mac).

Profile capturing will start. Pprof UI will be started
once capture is completed.

See [godoc](https://godoc.org/github.com/rakyll/autopprof) for other profile types.

## Why autopprof?

autopprof is a easy-to-setup pprof profile data collection library
for development time.
It highly depends on the standard library packages such as
[runtime/pprof](https://golang.org/pkg/runtime/pprof/) and the existing
tools such as `go tool pprof`.

Collecting and visualizing profiling data from Go programs is a
multi-step process. First, you need to collect and write the collected
data to a file. Then you should use the `go tool pprof` tool to analyze
and visualize.

autopprof makes it easier to collect and start the pprof UI with a
one-line configuration. It collects profiles once the process is triggered
with a SIGQUIT and starts the pprof UI with the collected data. Since it
does signal handling and starting the browser, it is only recommended
at development-time.

For production cases, please see the
[runtime/pprof](https://golang.org/pkg/runtime/pprof/)
and [net/http/pprof](https://golang.org/pkg/net/http/pprof/) packages.

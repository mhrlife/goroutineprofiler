package main

import "github.com/mhrlife/goroutineprofiler/profiler"

func main() {
	p := profiler.NewProfiler(profiler.DefaultConfig())
	p.Run()
}

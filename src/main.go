package main

import (
	"time"

	"github.com/ShuvraneelMitra/hungry-daemons/gui"
	"github.com/ShuvraneelMitra/hungry-daemons/profiler"
)

func main(){
	sm := profiler.NewGoRoutineSampler()
	sm.SetSamplingFrequency(5)
	stop := make(chan any)
	time.AfterFunc(10 * time.Second, func(){close(stop)})
	channel := sm.Sample(stop)
	profiler.Parse(stop, channel)

	gui.Run()
}
package main

import (
	"github.com/hink/go-blink1"
	"time"
)

var OffState = blink1.State{Duration: time.Duration(10) * time.Millisecond}

func main() {

	device, err := blink1.OpenNextDevice()
	defer device.Close()

	if err != nil {
		panic(err)
	}

	red := blink1.State{
		Red:      255,
		LED:      1,
		FadeTime: time.Duration(2) * time.Second,
		Duration: time.Duration(3) * time.Second,
	}

	device.SetState(red)

	green := blink1.State{
		Green:    255,
		LED:      2,
		FadeTime: time.Duration(2) * time.Second,
		Duration: time.Duration(3) * time.Second,
	}

	device.SetState(green)

	time.Sleep(time.Duration(7) * time.Second)

	device.SetState(OffState)
}

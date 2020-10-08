package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	log.Println("Started Log Shipping Sidecar")

	go func() {
		var gracefulStop = make(chan os.Signal, 1)
		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)

		sig := <-gracefulStop
		log.Println("Received stop signal:", sig)
		os.Exit(0)
	}()

	tick := time.Tick(ScanInterval)
	for range tick {
		scanForLogfiles()
	}
}

func scanForLogfiles() {
	log.Printf("Scanning for files %v in %v", FilePattern, Directory)
	files, err := ioutil.ReadDir(Directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if !f.IsDir() {
			log.Printf("Processing %v/%v %v %v", Directory, f.Name(), f.Size(), f.ModTime())
		}
	}
}

package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Directory to watch for changes
var Directory = logDirectory()

// FilePattern to match filenames
var FilePattern = logFilePattern()

// ScanInterval in seconds
var ScanInterval = logScanInterval()

func logScanInterval() time.Duration {
	intervalEnv, ok := os.LookupEnv("LOG_SCAN_INTERVAL")
	if !ok || len(intervalEnv) == 0 {
		log.Fatal("LOG_SCAN_INTERVAL environment variable not set.")
	}
	interval, err := strconv.ParseInt(intervalEnv, 10, 64)
	if err != nil {
		log.Fatal("Invalid value for LOG_SCAN_INTERVAL")
	}
	return time.Duration(interval) * time.Second
}

func logDirectory() string {
	name, ok := os.LookupEnv("LOG_DIRECTORY")
	if !ok || len(name) == 0 {
		log.Fatal("LOG_DIRECTORY environment variable not set.")
	}

	// Create the directory if it does not exist
	if _, err := os.Stat(name); os.IsNotExist(err) {
		os.Mkdir(name, os.ModePerm)
	}

	return name
}

func logFilePattern() *regexp.Regexp {
	pattern, ok := os.LookupEnv("LOG_FILE_PATTERN")
	if !ok || len(pattern) == 0 {
		log.Fatal("LOG_FILE_PATTERN environment variable not set.")
	}

	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal("LOG_FILE_PATTERN not a valid regular expression.", err)
	}

	return r
}

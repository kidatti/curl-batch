package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var curlFile = flag.String("curl", "", "Curl template file (required)")
	var csvFile = flag.String("csv", "", "CSV data file (required)")
	var outputFile = flag.String("output", "", "Output file (required)")
	var sleepMsec = flag.Int("sleep", 0, "Sleep duration in milliseconds between requests")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -curl <file> -csv <file> -output <file> [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s -curl curl.txt -csv users.csv -output results.txt -sleep 1000\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *curlFile == "" || *csvFile == "" || *outputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: All required flags must be specified\n\n")
		flag.Usage()
		os.Exit(1)
	}

	batch, err := NewCurlBatch(*curlFile, *csvFile, *outputFile, *sleepMsec)
	if err != nil {
		log.Fatalf("Failed to initialize curl batch: %v", err)
	}

	fmt.Printf("Starting batch execution with %d requests", len(batch.CSVData))
	if *sleepMsec > 0 {
		fmt.Printf(" (sleep: %dms between requests)", *sleepMsec)
	}
	fmt.Println("...")

	if err := batch.Run(); err != nil {
		log.Fatalf("Failed to run batch: %v", err)
	}

	fmt.Printf("Batch execution completed. Results saved to %s\n", *outputFile)
}

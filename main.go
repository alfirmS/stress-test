// main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourusername/yourproject/query"
	"github.com/yourusername/yourproject/stats"
)

func main() {
	var host, username, password, database, customQuery string
	var testInterval time.Duration
	var concurrency, iteration int

	flag.StringVar(&host, "host", "localhost:3306", "Database host")
	flag.StringVar(&username, "user", "root", "Database username")
	flag.StringVar(&password, "password", "", "Database password")
	flag.StringVar(&database, "database", "your_database_name", "Database name")
	flag.StringVar(&customQuery, "query", "", "Custom database query")
	flag.DurationVar(&testInterval, "interval", 10*time.Minute, "Interval between queries")
	flag.IntVar(&concurrency, "concurrency", 10, "Number of concurrent users")
	flag.IntVar(&iteration, "iteration", 5, "Number of iterations")

	flag.Parse()

	if customQuery == "" {
		fmt.Println("Please provide a custom query using the -query flag.")
		os.Exit(1)
	}

	// Initialize QueryStats with empty map and start time
	stats := &stats.QueryStats{
		QueriesPerUser: make(map[int]stats.UserStats),
		StartTime:      time.Now(),
	}

	results := make(chan time.Duration, concurrency*iteration)

	// Start stress test in a goroutine
	go query.StressTest(host, username, password, database, customQuery, testInterval, concurrency, iteration, results, stats)

	// Collect results from the channel
	var resultDurations []time.Duration
	for r := range results {
		resultDurations = append(resultDurations, r)
	}

	// Get the end time
	endTime := time.Now()

	// Print the final results
	stats.PrintResults(stats, stats.StartTime, endTime)
}

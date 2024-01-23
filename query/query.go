// query/query.go
package query

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/alfirmS/stress-test/stats"
	_ "github.com/go-sql-driver/mysql"
)

// runQuery collects statistics for each user's query
func runQuery(db *sql.DB, query string, interval time.Duration, concurrency, iteration int, wg *sync.WaitGroup, results chan<- time.Duration, stats *stats.QueryStats) {
	defer wg.Done()

	for i := 0; i < iteration; i++ {
		for user := 1; user <= concurrency; user++ {
			startTime := time.Now()

			_, err := db.Exec(query)
			if err != nil {
				log.Println("Error executing query:", err)
				stats.FailedQueryCount++
			} else {
				stats.SuccessfulQueries++
			}

			elapsedTime := time.Since(startTime)

			stats.TotalQueries++
			// Update user-specific statistics
			userStats, ok := stats.QueriesPerUser[user]
			if !ok {
				userStats = stats.UserStats{}
			}
			userStats.Update(elapsedTime, query)
			stats.QueriesPerUser[user] = userStats
		}
		time.Sleep(interval)
	}

	elapsedTime := time.Since(stats.StartTime)
	results <- elapsedTime

	stats.UpdateQueryTimeStats(elapsedTime)
}

// StressTest simulates stress testing on the database
func StressTest(host, username, password, database, query string, interval time.Duration, concurrency, iteration int, results chan<- time.Duration, stats *stats.QueryStats) {
	// Establish database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, database))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var wg sync.WaitGroup

	// Run queries concurrently
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go runQuery(db, query, interval, concurrency, iteration, &wg, results, stats)
	}

	// Wait for all queries to finish
	wg.Wait()
	close(results)
}

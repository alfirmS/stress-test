// stats/stats.go
package stats

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// UserStats stores specific statistics for each user
type UserStats struct {
	TotalQueries  int
	TotalTime     time.Duration
	LongestTime   time.Duration
	ShortestTime  time.Duration
	LongestQuery  string
	ShortestQuery string
}

// QueryStats stores global statistics including statistics for each user
type QueryStats struct {
	TotalQueries      int
	QueriesPerUser    map[int]UserStats
	AverageQueryTime  time.Duration
	MaxQueryTime      time.Duration
	MinQueryTime      time.Duration
	FailedQueryCount  int
	SuccessfulQueries int
	StartTime         time.Time
	mu                sync.Mutex
}

// UpdateQueryTimeStats updates global execution time statistics
func (stats *QueryStats) UpdateQueryTimeStats(elapsedTime time.Duration) {
	stats.mu.Lock()
	defer stats.mu.Unlock()

	if elapsedTime > stats.MaxQueryTime {
		stats.MaxQueryTime = elapsedTime
	}

	if elapsedTime < stats.MinQueryTime || stats.MinQueryTime == 0 {
		stats.MinQueryTime = elapsedTime
	}

	stats.AverageQueryTime = time.Duration(int64(stats.AverageQueryTime)*int64(stats.TotalQueries-1)/int64(stats.TotalQueries) + int64(elapsedTime)/int64(stats.TotalQueries))
}

// Update updates the user-specific statistics
func (userStats *UserStats) Update(elapsedTime time.Duration, query string) {
	userStats.TotalQueries++
	userStats.TotalTime += elapsedTime

	if elapsedTime > userStats.LongestTime {
		userStats.LongestTime = elapsedTime
		userStats.LongestQuery = query
	}

	if elapsedTime < userStats.ShortestTime || userStats.ShortestTime == 0 {
		userStats.ShortestTime = elapsedTime
		userStats.ShortestQuery = query
	}
}

// PrintResults displays results and statistics
func (stats *QueryStats) PrintResults(startTime, endTime time.Time) {
	stats.mu.Lock()
	defer stats.mu.Unlock()

	fmt.Println("Results:")
	fmt.Printf("1. Total queries executed: %d\n", stats.TotalQueries)
	stats.printQueriesPerUser()
	fmt.Printf("3. Average query completion time: %s\n", stats.AverageQueryTime)
	fmt.Printf("4. Number of unsuccessful queries (in percentage): %.2f%%\n", stats.calculatePercentage(stats.FailedQueryCount, stats.TotalQueries))
	fmt.Printf("5. Number of successful queries (in percentage): %.2f%%\n", stats.calculatePercentage(stats.SuccessfulQueries, stats.TotalQueries))
	fmt.Printf("6. CLI start time: %s\n", startTime.Format(time.RFC3339))
	fmt.Printf("7. CLI end time: %s\n", endTime.Format(time.RFC3339))
}

// printQueriesPerUser displays queries executed per user
func (stats *QueryStats) printQueriesPerUser() {
	fmt.Println("2. Queries per user:")
	fmt.Println("   User   |   Total Queries   |   Average Time   |   Fastest Time   |   Slowest Time")
	for user, userStats := range stats.QueriesPerUser {
		averageTime := userStats.TotalTime / time.Duration(userStats.TotalQueries)
		fmt.Printf("   %d     |   %d              |   %s           |   %s           |   %s\n", user, userStats.TotalQueries, averageTime, userStats.ShortestTime, userStats.LongestTime)
	}
}

// calculatePercentage calculates the percentage of successful or unsuccessful queries
func (stats *QueryStats) calculatePercentage(count, total int) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(count) / float64(total) * 100
}

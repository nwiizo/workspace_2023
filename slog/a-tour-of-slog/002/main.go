package main

import (
	"database/sql"
	"log/slog"
	"os"
)

type CategorySummary struct {
	Title    string
	Tasks    int
	diff     int
	AvgValue float64
}

func createTables(db *sql.DB) {
	db.Exec("CREATE TABLE tasks (id INTEGER PRIMARY KEY, title TEXT, value INTEGER, category TEXT)")
}

func createCategorySummaries(db *sql.DB) []CategorySummary {
	rows, _ := db.Query("SELECT category, COUNT(*), AVG(value) FROM tasks GROUP BY category")
	defer rows.Close()

	summaries := []CategorySummary{}
	for rows.Next() {
		var summary CategorySummary
		rows.Scan(&summary.Title, &summary.Tasks, &summary.AvgValue)
		summaries = append(summaries, summary)
	}
	return summaries
}

func main() {
	// Default logger
	logger := slog.Default()
	logger.Info("hello, world", "user", os.Getenv("USER"))
}

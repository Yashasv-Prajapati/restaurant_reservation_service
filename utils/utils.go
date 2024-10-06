package utils

import (
	"fmt"
	"time"
)

// Function to compare times ignoring the date
func IsWithinOpenHours(bookingTime, openingTime, closingTime time.Time) bool {
	// Create new time.Time instances with the same hour and minute
	// but set the date to a fixed value (e.g., January 1, 2000)
	bookingTimeOnly := time.Date(2000, 1, 1, bookingTime.Hour(), bookingTime.Minute(), 0, 0, time.UTC)
	openingTimeOnly := time.Date(2000, 1, 1, openingTime.Hour(), openingTime.Minute(), 0, 0, time.UTC)
	closingTimeOnly := time.Date(2000, 1, 1, closingTime.Hour(), closingTime.Minute(), 0, 0, time.UTC)

	// Compare the times
	return bookingTimeOnly.After(openingTimeOnly) && bookingTimeOnly.Before(closingTimeOnly)
}

func PrintConsole( args ...interface{}) {
	fmt.Print("-> ")
	fmt.Println(args...)
}
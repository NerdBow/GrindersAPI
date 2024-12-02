package logs

import (
	"time"
)

// The struct and date format for each log in the system
type Log struct {
	Id       int
	Date     time.Time
	Duration time.Duration
	Name     string
	Category string
	Tags     []string
	UserId   int
}

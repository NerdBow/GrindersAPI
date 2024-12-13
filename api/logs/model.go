package logs

import (
	"fmt"
	"time"
)

// The struct and date format for each log in the system
type Log struct {
	Id       int           `json:"id"`
	Date     int64         `json:"date"`
	Duration time.Duration `json:"duration"`
	Name     string        `json:"name"`
	Category string        `json:"category"`
	Tags     []string      `json:"tags"`
	UserId   int           `json:"userId"`
}

func (log Log) Validate() bool {
	if log.Date == 0 {
		return false
	}

	if log.Duration == 0 {
		return false
	}

	if log.Name == "" {
		return false
	}

	if log.Category == "" {
		return false
	}

	if log.UserId == 0 {
		return false
	}
	return true
}

func (log Log) String() string {
	return fmt.Sprintf("Id: %d\nDate: %s\nDuration: %d\nName: %s\nCategory: %s\n", log.Id, time.Unix(log.Date, 0), log.Duration, log.Name, log.Category)
}

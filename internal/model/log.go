package model

import "errors"

var (
	InvalidIdErr       = errors.New("Log id must be greater than 0")
	InvalidDateErr     = errors.New("Log date must be greater than 0")
	InvalidDurationErr = errors.New("Log duration must be greater than 0")
	InvalidNameErr     = errors.New("Log name must not be blank")
	InvalidCategoryErr = errors.New("Log category must not be blank")
	InvalidGoalErr     = errors.New("Log goal must not be blank")
	InvalidUserIdErr   = errors.New("Log user id must be greater than 0")
)

// A struct of the information contatined in each log.
// Date and Duration are int64 of seconds.
type Log struct {
	Id       int64  `json:"id"`
	Date     int64  `json:"date"`
	Duration int64  `json:"duration"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Goal     string `json:"goal"`
	UserId   int    `json:"userId"`
}

// Validate checks for any problems with the structure.
// Returns map[string]string of fields as keys and the problems with the field if there are any.
func (log Log) Validate() map[string]error {
	problems := make(map[string]error)

	if log.Id <= 0 {
		problems["id"] = InvalidIdErr
	}
	if log.Date <= 0 {
		problems["date"] = InvalidDateErr
	}
	if log.Duration <= 0 {
		problems["duration"] = InvalidDurationErr
	}
	if log.Name == "" {
		problems["name"] = InvalidNameErr
	}
	if log.Category == "" {
		problems["category"] = InvalidCategoryErr
	}
	if log.Goal == "" {
		problems["goal"] = InvalidGoalErr
	}
	if log.UserId <= 0 {
		problems["userId"] = InvalidUserIdErr
	}

	return problems
}

// Takes the existing log and fill in all the unfilled fields from that of the otherLog
func (log *Log) Merge(otherLog Log) {
	if log.Date == 0 {
		log.Date = otherLog.Date
	}

	if log.Duration == 0 {
		log.Duration = otherLog.Duration
	}

	if log.Name == "" {
		log.Name = otherLog.Name
	}

	if log.Category == "" {
		log.Category = otherLog.Category
	}

	if log.Goal == "" {
		log.Goal = otherLog.Goal
	}
}

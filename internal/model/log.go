package model

// A struct of the information contatined in each log.
// Date and Duration are int64 of seconds.
type Log struct {
	Id       int    `json:"id"`
	Date     int64  `json:"date"`
	Duration int64  `json:"duration"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Goal     string `json:"goal"`
	UserId   int    `json:"userId"`
}

// Validate checks for any problems with the structure.
// Returns map[string]string of fields as keys and the problems with the field if there are any.
func (log Log) Validate() map[string]string {
	problems := make(map[string]string)

	if log.Id == 0 {
		problems["id"] = "No id"
	}
	if log.Date == 0 {
		problems["date"] = "No date"
	}
	if log.Duration == 0 {
		problems["duration"] = "No duration"
	}
	if log.Name == "" {
		problems["Name"] = "No name"
	}
	if log.Category == "" {
		problems["category"] = "No category"
	}
	if log.Goal == "" {
		problems["goal"] = "No goal"
	}
	if log.UserId == 0 {
		problems["userId"] = "No userId"
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
		log.Category = otherLog.Category
	}
}

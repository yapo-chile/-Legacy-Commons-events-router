package domain

import "time"

// Event defines an action triggered by yapo's operation
type Event struct {
	Type    string
	Date    time.Time
	Content interface{}
}

package test_records

import "time"

type TestRecord struct {
	ID      string    `json:"id"`
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
	Stored  time.Time `json:"stored"`
}

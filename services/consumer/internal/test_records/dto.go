package test_records

import "time"

type CreateTestRecordDTO struct {
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
	Stored  time.Time `json:"stored"`
}

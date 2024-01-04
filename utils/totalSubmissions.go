package utils

import (
	"encoding/json"
)

func GetTotalSubmission(submissionCalendar string) (int, error) {

	// Unmarshal the JSON string to a map
	var calendar map[string]int
	err := json.Unmarshal([]byte(submissionCalendar), &calendar)
	if err != nil {
		return 0, err
	}

	// Calculate the total submissions
	totalSubmissions := 0
	for _, count := range calendar {
		totalSubmissions += count
	}

	return totalSubmissions, nil
}
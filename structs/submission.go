package structs

type Lang struct {
	Name        string `json:"name"`
	VerboseName string `json:"verboseName"`
}

type Question struct {
	QuestionId    string `json:"questionId"`
	QuestionTitle string `json:"questionTitle"`
}

type SubmissionDetails struct {
	Code       string   `json:"code"`
	Timestamp  int      `json:"timestamp"`
	StatusCode int      `json:"statusCode"`
	Lang       Lang     `json:"lang"`
	Question   Question `json:"question"`
	TopicTags  []string `json:"topicTags"`
}

type SubmissionDetailsResponseData struct {
	SubmissionDetails SubmissionDetails `json:"submissionDetails"`
}

type SubmissionDetailsResponse struct {
	Data SubmissionDetailsResponseData `json:"data"`
}

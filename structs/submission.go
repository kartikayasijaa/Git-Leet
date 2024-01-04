package structs

type SubmissionDetails struct {
	Code       string `json:"code"`
	Timestamp  int    `json:"timestamp"`
	StatusCode int    `json:"statusCode"`
	Lang       struct {
		Name        string `json:"name"`
		VerboseName string `json:"verboseName"`
	} `json:"lang"`
	Question struct {
		QuestionId string `json:"questionId"`
	} `json:"question"`
	TopicTags []string `json:"topicTags"`
}
type SubmissionDetailsResponse struct {
	Data struct {
		SubmissionDetails SubmissionDetails `json:"submissionDetails"`
	} `json:"data"`
}

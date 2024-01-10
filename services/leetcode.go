package services

import (
	"errors"
	"gitleet/structs"
	"gitleet/utils"
	"os"
	"github.com/dustyRAIN/leetcode-api-go/leetcodeapi"
)

type LeetcodeService struct{}

func (h *LeetcodeService) GetRecentSubmission(LeetcodeUsername string) ([]leetcodeapi.AcSubmission, error) {
	totalSubmissions, err := h.GetTotalSubmission(LeetcodeUsername)
	if err != nil {
		return nil, err
	}
	prevSubmissions := 450
	newSubmission := int(totalSubmissions) - prevSubmissions
	if newSubmission <= 0 {
		return nil, errors.New("No new submissions found")
	}

	recentSubmissions, err := leetcodeapi.GetUserRecentAcSubmissions(LeetcodeUsername, newSubmission)

	if err != nil {
		return nil, err
	}

	return recentSubmissions, nil
}

func (h *LeetcodeService) GetCode(submissionId string) (structs.SubmissionDetails, error) {
	var responseBody structs.SubmissionDetailsResponse
	payload := `{
		"query": "query submissionDetails($submissionId: Int!) { submissionDetails(submissionId: $submissionId) {  code timestamp statusCode lang { name verboseName } question { questionId questionTitle } topicTags { tagId slug name } } }",
		"variables": {
			"submissionId": "` + submissionId + `"
		}
	}`

	session := os.Getenv("SESSION")
	if len(session) == 0 {
		return structs.SubmissionDetails{}, errors.New("No Session Key")
	}

	leetcodeapi.SetCredentials(session, "f")
	err := (&leetcodeapi.Util{}).MakeGraphQLRequest(payload, &responseBody)
	leetcodeapi.RemoveCredentials()
	if err != nil {
		return structs.SubmissionDetails{}, err
	}
	code := responseBody.Data.SubmissionDetails.Code

	if len(code) == 0 {
		return structs.SubmissionDetails{}, errors.New("No cookie")
	}
	return responseBody.Data.SubmissionDetails, nil
}

func (h *LeetcodeService) GetTotalSubmission(LeetcodeUsername string) (int32, error) {
	submissions, err := leetcodeapi.GetUserProfileCalendar(LeetcodeUsername)
	if err != nil {
		return 0, err
	}
	totalSubmissions, err := utils.GetTotalSubmission(submissions.SubmissionCalendar)
	if err != nil {
		return 0, err
	}

	return int32(totalSubmissions), nil
}
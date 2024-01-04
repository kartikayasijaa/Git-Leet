package services

import (
	"errors"
	"gitleet/structs"
	"gitleet/utils"
	"os"

	"github.com/dustyRAIN/leetcode-api-go/leetcodeapi"
)

type services struct{}

func GetRecentSubmission(userId string) ([]leetcodeapi.AcSubmission, error) {
	submissions, err := leetcodeapi.GetUserProfileCalendar(userId)
	if err != nil {
		return nil, err
	}
	totalSubmissions, err := utils.GetTotalSubmission(submissions.SubmissionCalendar)
	if err != nil {
		return nil, err
	}

	prevSubmissions := 500
	newSubmission := totalSubmissions - prevSubmissions
	if newSubmission == 0 {
		return nil, errors.New("No new submissions found")
	}

	recentSubmissions, err := leetcodeapi.GetUserRecentAcSubmissions(userId, newSubmission)

	if err != nil {
		return nil, err
	}

	return recentSubmissions, nil
}

func GetCode(submissionId string) (string, error) {
	var responseBody structs.SubmissionDetailsResponse
	payload := `{
		"query": "query submissionDetails($submissionId: Int!) { submissionDetails(submissionId: $submissionId) {  code timestamp statusCode lang { name verboseName } question { questionId titleSlug hasFrontendPreview } topicTags { tagId slug name } } }",
		"variables": {
			"submissionId": "` + submissionId + `"
		}
	}`

	session := os.Getenv("SESSION")
	if len(session) == 0 {
		return "", errors.New("No Session Key")
	}

	leetcodeapi.SetCredentials(session, "f")
	err := (&leetcodeapi.Util{}).MakeGraphQLRequest(payload, &responseBody)
	leetcodeapi.RemoveCredentials()
	if err != nil {
		return "", err
	}
	code := responseBody.Data.SubmissionDetails.Code

	if len(code) == 0 {
		return "", errors.New("No cookie")
	}
	return responseBody.Data.SubmissionDetails.Code, nil
}

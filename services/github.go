package services

import (
	"context"
	"fmt"
	"gitleet/structs"
	"gitleet/utils"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GetGitHubClient(accessToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func PushToGithub(accessToken string, owner string, repo string, branch string, submissionDetail structs.SubmissionDetails) error {
	client := GetGitHubClient(accessToken)

	// Repository information
	repoOwner := owner
	repoName := repo

	questionTitle := utils.SanitizeFileName(submissionDetail.Question.QuestionTitle)
	// File information
	fileName := fmt.Sprintf("%s-%s.cpp", submissionDetail.Question.QuestionId, questionTitle)
	// C++ code as a string
	cppCode := submissionDetail.Code

	// Get the current commit SHA
	branchName := branch // or the name of your default branch
	opts := &github.RepositoryContentGetOptions{
		Ref: branchName,
	}
	fileInfo, _, _, err := client.Repositories.GetContents(context.Background(), repoOwner, repoName, fileName, opts)
	// Check if the file exists
	if err != nil {
		// File doesn't exist, create a new file
		commitMessage := fileName
		createFileOptions := &github.RepositoryContentFileOptions{
			Message: &commitMessage,
			Content: []byte(cppCode),
		}
		_, _, err = client.Repositories.CreateFile(context.Background(), repoOwner, repoName, fileName, createFileOptions)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return err
		}
	} else {
		// File exists, get the SHA and update the file
		commitMessage := fileName
		updateFileOptions := &github.RepositoryContentFileOptions{
			Message: &commitMessage,
			Content: []byte(cppCode),
			SHA:     fileInfo.SHA,
		}
		_, _, err = client.Repositories.UpdateFile(context.Background(), repoOwner, repoName, fileName, updateFileOptions)
		if err != nil {
			fmt.Printf("Error updating file: %v\n", err)
			return err
		}
	}

	return nil
}

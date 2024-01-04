package services

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func createGitHubClient(accessToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func PushToGithub(code string) error {
	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	client := createGitHubClient(accessToken)

	// Repository information
	repoOwner := "kartikayasijaa"
	repoName := "DSA-rev"

	// File information
	fileName := "new_file.cpp"

	// C++ code as a string
	cppCode := code

	// Get the current commit SHA
	branchName := "main" // or the name of your default branch
	opts := &github.RepositoryContentGetOptions{
		Ref: branchName,
	}
	fileInfo, _, _, err := client.Repositories.GetContents(context.Background(), repoOwner, repoName, fileName, opts)

	// Check if the file exists
	if err != nil {
		// File doesn't exist, create a new file

		commitMessage := "Add new file"
		_, _, err = client.Repositories.CreateFile(context.Background(), repoOwner, repoName, fileName, &github.RepositoryContentFileOptions{
			Message: &commitMessage,
			Content: []byte(cppCode),
		})
		if err != nil {
			return err
		}
	}

	// File exists, get the SHA and update the file
	commitMessage := "Update file"
	_, _, err = client.Repositories.CreateFile(context.Background(), repoOwner, repoName, fileName, &github.RepositoryContentFileOptions{
		Message: &commitMessage,
		Content: []byte(cppCode),
		SHA:     fileInfo.SHA,
	})
	if err != nil {
		return err
	}

	return nil
}

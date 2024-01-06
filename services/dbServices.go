package services

import (
	"context"
	"gitleet/structs"

	"github.com/google/go-github/github"
	"gorm.io/gorm"
)

type DBServices struct {
	DB      *gorm.DB
	Context context.Context
}

func DBServicesHandler(db *gorm.DB, ctx context.Context) *DBServices {
	return &DBServices{
		DB:      db,
		Context: ctx,
	}
}

func (h *DBServices) CreateOrUpdateUser(githubUser *github.User) (*structs.User, error) {
	user := &structs.User{
		Email:                  githubUser.GetEmail(),
		Name:                   githubUser.GetName(),
		LeetcodeUsername:       "",
		LeetcodePrevSubmission: 0,
		GithubUsername:         githubUser.GetLogin(),
		GithubRepoName:         "",
		GithubAvatarURL:        githubUser.GetAvatarURL(),
	}

	// If found, update the existing user with the data from the Assign method.
	// If not found, create a new user with the provided data.
	result := h.DB.Where(&structs.User{GithubUsername: user.GithubUsername}).Assign(user).FirstOrCreate(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

package services

import (
	"context"
	"errors"
	"gitleet/structs"

	"github.com/google/go-github/github"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DBService struct {
	DB      *gorm.DB
	Context context.Context
}

func DBServicesHandler(db *gorm.DB, ctx context.Context) *DBService {
	return &DBService{
		DB:      db,
		Context: ctx,
	}
}

func (h *DBService) CreateOrUpdateUser(githubUser *github.User) (*structs.User, error) {
	// Check if the user already exists
	existingUser := &structs.User{}
	result := h.DB.Where(&structs.User{GithubUsername: githubUser.GetLogin()}).First(existingUser)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	var user *structs.User
	// If the user doesn't exist, generate a new UUID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user = &structs.User{
			UserID:                 uuid.New().String(),
			Email:                  githubUser.GetEmail(),
			Name:                   githubUser.GetName(),
			LeetcodeUsername:       "",
			LeetcodePrevSubmission: 0,
			GithubUsername:         githubUser.GetLogin(),
			GithubRepoName:         "",
			GithubAvatarURL:        githubUser.GetAvatarURL(),
		}
	} else {
		// If the user already exists, update the existing user with the data from the Assign method
		user = existingUser
		user.Email = githubUser.GetEmail()
		user.Name = githubUser.GetName()
		user.GithubAvatarURL = githubUser.GetAvatarURL()
	}
	// Save the user to the database
	result = h.DB.Where(&structs.User{GithubUsername: user.GithubUsername}).Assign(user).FirstOrCreate(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (h *DBService) UpdateLeetcode(leetcodeID string, userID string, leetcodePrev int32) error {
	user := &structs.User{
		LeetcodeUsername:       leetcodeID,
		LeetcodePrevSubmission: leetcodePrev,
	}
	result := h.DB.Where(&structs.User{UserID: userID}).Updates(user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no user found with userID")
	}

	return nil
}

func (h *DBService) UpdateGithub(repoName string, branch string, userID string) error {
	user := &structs.User{
		GithubRepoName:   repoName,
		GithubRepoBranch: branch,
	}
	result := h.DB.Where(&structs.User{UserID: userID}).Updates(user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no user found with userID")
	}

	return nil
}

func (h *DBService) GetUser(userId string) (*structs.User, error) {
	user := new(structs.User)
	res := h.DB.Where(&structs.User{UserID: userId}).First(&user)

	if res.Error != nil {
		return &structs.User{}, res.Error
	}
	return user, nil
}

func (h *DBService) UpdatePrevSubmission(userId string, leetcodePrev int32) error {
	user := &structs.User{
		LeetcodePrevSubmission: leetcodePrev,
	}
	result := h.DB.Where(&structs.User{UserID: userId}).Updates(user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no user found with userID")
	}

	return nil
}

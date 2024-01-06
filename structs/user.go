package structs

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email                  string `json:"email"`
	Name                   string `json:"name"`
	LeetcodeUsername       string `json:"leetcode_username"`
	LeetcodePrevSubmission int32  `json:"leetcode_prev_submission"`
	GithubUsername         string `json:"github_username"`
	GithubRepoName         string `json:"github_repo"`
	GithubAvatarURL        string `json:"github_avatar_url"`
}

// MigrateUser performs database migration for the User model.
func MigrateUser(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}

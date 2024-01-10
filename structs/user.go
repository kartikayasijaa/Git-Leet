package structs

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID                 string `gorm:"primaryKey" json:"userId"`
	Email                  string `json:"email"`
	Name                   string `json:"name"`
	LeetcodeUsername       string `json:"leetcode_username"`
	LeetcodePrevSubmission int32  `json:"leetcode_prev_submission"`
	GithubUsername         string `json:"github_username"`
	GithubRepoName         string `json:"github_repo"`
	GithubRepoBranch       string `json:"github_repo_branch"`
	GithubAvatarURL        string `json:"github_avatar_url"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// MigrateUser performs database migration for the User model.
func MigrateUser(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}

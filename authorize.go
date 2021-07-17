package storage

import "github.com/gislik/gorm"
import "time"

// Authorize data model
type Authorize struct {
	ClientID            string    // Client information
	Code                string    `gorm:"primary_key"` // Authorization code
	ExpiresIn           int32     // Token expiration in seconds
	Scope               string    // Requested scope
	RedirectUri         string    // Redirect Uri from request
	State               string    // State data from request
	CreatedAt           time.Time // Date created
	UserData            string    // Data to be passed to storage. Not used by the library.
	CodeChallenge       string    // Optional code_challenge as described in rfc7636
	CodeChallengeMethod string    // Optional code_challenge_method as described in rfc7636
}

// TableName is used by `gorm`
func (Authorize) TableName(db *gorm.DB) string {
	return gorm.DefaultTableNameHandler(db, "oauth_authorize")
}

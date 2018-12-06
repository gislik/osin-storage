package storage

import (
	"github.com/gislik/gorm"

	"time"
)

//Access data model
type Access struct {
	ClientID     string    // Client information
	Authorize    string    // Authorize data, for authorization code
	PrvAccess    string    // Previous access data, for refresh token
	AccessToken  string    `gorm:"primary_key"` // Access token
	RefreshToken string    // Refresh Token. Can be blank
	ExpiresIn    int32     // Token expiration in seconds
	Scope        string    // Requested scope
	RedirectUri  string    // Redirect Uri from request
	CreatedAt    time.Time // Date created
	UserData     string    // Data to be passed to storage. Not used by the library.
}

// TableName is used by `gorm`
func (Access) TableName(db *gorm.DB) string {
	return gorm.DefaultTableNameHandler(db, "oauth_access")
}

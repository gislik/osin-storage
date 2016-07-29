package storage

import "time"

// Authorize data model
type Authorize struct {
	ClientID    string    // Client information
	Code        string    `gorm:"primary_key"` // Authorization code
	ExpiresIn   int32     // Token expiration in seconds
	Scope       string    // Requested scope
	RedirectUri string    // Redirect Uri from request
	State       string    // State data from request
	CreatedAt   time.Time // Date created
	UserData    string    // Data to be passed to storage. Not used by the library.
}

//Tablename of Authorize data model
func (Authorize) TableName() string {
	return "oauth_authorize"
}

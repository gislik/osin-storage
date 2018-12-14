package storage

import (
	"github.com/gislik/gorm"
)

//Client model
type Client struct {
	ID          string `gorm:"primary_key"`
	Secret      string
	RedirectUri string
	UserData    string
}

// TableName is used by `gorm`
func (Client) TableName(db *gorm.DB) string {
	return gorm.DefaultTableNameHandler(db, "oauth_client")
}

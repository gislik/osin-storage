package storage

//Client model
type Client struct {
	ID          string `gorm:"primary_key"`
	Secret      string
	RedirectUri string
	UserData    string
}

//TableName of Client model
func (Client) TableName() string {
	return "oauth_client"
}

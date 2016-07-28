# osin-storage
OSIN GORM storage

This package implements the storage for [OSIN](https://github.com/RangelReale/osin) using [GORM](https://github.com/jinzhu/gorm).

## Installation

go get github.com/collinsss/osin-storage

## Usage

```go
import (
	"github.com/RangelReale/osin"
	"github.com/jinzhu/gorm"
	"github.com/collinsss/osin-storage"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//InitDB initial gorm db
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "host=localhost user=user dbname=dbname sslmode=disable password=password")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&storage.Access{},
		&storage.Authorize{},
		&storage.Client{},
	)
	return db, nil
}

func main() {

	db, err := InitDB()
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	defer db.Close()

	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}

	//Add or delete the AllowedAccessType as you want
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.AUTHORIZATION_CODE}
	sconfig.AllowGetAccessRequest = true
	sconfig.AllowClientSecretInParams = true
	server := osin.NewServer(sconfig, oauthstorage.NewStorage(db))

	// For further details how to use osin server check osin documentation

```
package storage

import (
	"encoding/json"

	"github.com/gislik/gorm"
	"github.com/openshift/osin"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db}
}

// Clone the storage if needed. For example, using mgo, you can clone the session with session.Clone
// to avoid concurrent access problems.
// This is to avoid cloning the connection at each method access.
// Can return itself if not a problem.
func (s *Storage) Clone() osin.Storage {
	return s
}

// Close the resources the Storage potentially holds (using Clone for example)
func (s *Storage) Close() {
}

var err error

// GetClient loads the client by id (client_id)
func (s *Storage) GetClient(id string) (osin.Client, error) {
	var c Client

	if err = s.db.Where("id = ?", id).First(&c).Error; err == nil {
		oc := osin.DefaultClient{
			Id:          c.ID,
			Secret:      c.Secret,
			RedirectUri: c.RedirectUri,
			UserData:    c.UserData,
		}
		return &oc, nil
	}
	return nil, err
}

// SaveClient saves client
func (s *Storage) SaveClient(c osin.Client) error {
	client := Client{
		ID:          c.GetId(),
		Secret:      c.GetSecret(),
		RedirectUri: c.GetRedirectUri(),
	}
	if c.GetUserData() != nil {
		v, err := userDataToString(c.GetUserData())
		if err != nil {
			return err
		}
		client.UserData = v
	}
	return s.db.Create(&client).Error

}

// RemoveClient removes the client with matching id
func (s *Storage) RemoveClient(id string) error {
	var a Client
	if err = s.db.Where("id = ?", id).First(&a).Error; err == nil {
		if err = s.db.Delete(&a).Error; err == nil {
			return nil
		}
		return err
	}
	return err
}

// SaveAuthorize saves authorize data.
func (s *Storage) SaveAuthorize(data *osin.AuthorizeData) error {
	authorize := Authorize{
		ClientID:    data.Client.GetId(),
		Code:        data.Code,
		ExpiresIn:   data.ExpiresIn,
		RedirectUri: data.RedirectUri,
		Scope:       data.Scope,
		State:       data.State,
		CreatedAt:   data.CreatedAt,
	}
	if data.UserData != nil {
		v, err := userDataToString(data.UserData)
		if err != nil {
			return err
		}
		authorize.UserData = v
	}
	return s.db.Create(&authorize).Error
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (s *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	var authorize Authorize
	if err = s.db.Where("code = ?", code).First(&authorize).Error; err == nil {
		client, err := s.GetClient(authorize.ClientID)
		if err != nil {
			return nil, err
		}
		oa := &osin.AuthorizeData{
			Client:      client,
			Code:        authorize.Code,
			ExpiresIn:   authorize.ExpiresIn,
			Scope:       authorize.Scope,
			RedirectUri: authorize.RedirectUri,
			State:       authorize.State,
			CreatedAt:   authorize.CreatedAt,
		}

		if authorize.UserData != "" {
			oa.UserData = authorize.UserData
		}
		return oa, nil
	}
	return nil, err
}

// RemoveAuthorize revokes or deletes the authorization code.
func (s *Storage) RemoveAuthorize(code string) error {
	var a Authorize
	if err = s.db.Where("code = ?", code).First(&a).Error; err == nil {
		if err = s.db.Delete(&a).Error; err == nil {
			return nil
		}
		return err
	}
	return err
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (s *Storage) SaveAccess(data *osin.AccessData) error {
	access := Access{
		ClientID:     data.Client.GetId(),
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresIn:    data.ExpiresIn,
		Scope:        data.Scope,
		RedirectUri:  data.RedirectUri,
		CreatedAt:    data.CreatedAt,
	}

	if data.UserData != nil {
		v, err := userDataToString(data.UserData)
		if err != nil {
			return err
		}
		access.UserData = v
	}

	if data.AccessData != nil {
		access.PrvAccess = data.AccessData.AccessToken
	}

	if data.AuthorizeData != nil {
		access.Authorize = data.AuthorizeData.Code
	}

	return s.db.Create(&access).Error
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadAccess(code string) (*osin.AccessData, error) {
	var a Access
	if err = s.db.Where("access_token = ?", code).First(&a).Error; err == nil {
		client, err := s.GetClient(a.ClientID)
		if err != nil {
			return nil, err
		}
		authorize, err := s.LoadAuthorize(a.Authorize)

		oa := &osin.AccessData{
			Client:        client,
			AuthorizeData: authorize,
			AccessToken:   a.AccessToken,
			RefreshToken:  a.RefreshToken,
			ExpiresIn:     a.ExpiresIn,
			Scope:         a.Scope,
			RedirectUri:   a.RedirectUri,
			CreatedAt:     a.CreatedAt,
		}
		if a.UserData != "" {
			oa.UserData = a.UserData
		}
		return oa, nil
	}
	return nil, err
}

// RemoveAccess revokes or deletes an AccessData.
func (s *Storage) RemoveAccess(code string) error {
	var a Access
	if err = s.db.Where("access_token = ?", code).First(&a).Error; err == nil {
		if err = s.db.Delete(&a).Error; err == nil {
			return nil
		}
		return err
	}
	return err
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadRefresh(code string) (*osin.AccessData, error) {
	var (
		a Access
	)
	if err = s.db.Where("refresh_token = ?", code).First(&a).Error; err == nil {
		return s.LoadAccess(a.AccessToken)
	}
	return nil, err
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (s *Storage) RemoveRefresh(code string) error {
	var a Access
	if err = s.db.Where("refresh_token = ?", code).First(&a).Error; err == nil {
		if err = s.db.Delete(&a).Error; err == nil {
			return nil
		}
		return err
	}
	return err
}

func userDataToString(userData interface{}) (string, error) {
	if userData == nil {
		return "", nil
	}
	v, err := json.Marshal(userData)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

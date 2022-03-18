package session

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const SessionId = "id"

var Manager *ConfigSession

type UserInfo struct {
	Name          string `json:"name" form:"name"`
	Email         string `json:"email" form:"email"`
	IsActive      int    `json:"is_active" form:"is_active"`
	IsMultiBranch int    `json:"is_multi_branch" form:"is_multi_branch"`
	TypeUser      int    `json:"type_user" form:"type_user"`
	Branch        string `json:"branch" form:"branch"`
	Type          int    `json:"type"`
	UserID        string `json:"user_id" form:"user_id"`
	UserRoleID    string `json:"user_role_id" form:"user_role_id"`
}

type FlashMessage struct {
	Type    string
	Message string
	Data    interface{}
}

type ConfigSession struct {
	store    *sessions.CookieStore
	valueKey string
}

func NewSessionManager(store *sessions.CookieStore) *ConfigSession {
	s := new(ConfigSession)
	s.valueKey = "data"
	s.store = store

	return s
}

func (s *ConfigSession) Get(c echo.Context, name string) (interface{}, error) {
	session, err := s.store.Get(c.Request(), name)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}
	if val, ok := session.Values[s.valueKey]; ok {
		return val, nil
	} else {
		return nil, nil
	}
}

func (s *ConfigSession) Set(c echo.Context, name string, value interface{}) error {
	session, _ := s.store.Get(c.Request(), name)
	session.Values[s.valueKey] = value

	err := session.Save(c.Request(), c.Response())
	return err
}

func (s *ConfigSession) Delete(c echo.Context, name string) error {
	session, err := s.store.Get(c.Request(), name)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(c.Request(), c.Response())
}

func (s *ConfigSession) GetWithKeyValues(c echo.Context, name string, keyValue string) (interface{}, error) {
	session, err := s.store.Get(c.Request(), name)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}
	if val, ok := session.Values[keyValue]; ok {
		return val, nil
	} else {
		return nil, nil
	}
}

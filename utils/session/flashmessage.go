package session

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func SetFlashMessage(c echo.Context, message string, key string, data interface{}) {
	session, err := Manager.store.Get(c.Request(), "flash-message")
	if err != nil {
		panic(err)
	}
	mapMessage := FlashMessage{
		Type:    key,
		Message: message,
		Data:    data,
	}
	session.AddFlash(mapMessage)
	err = session.Save(c.Request(), c.Response())
	if err != nil {
		panic(err)
	}
}

func GetFlashMessage(c echo.Context) FlashMessage {
	session, err := Manager.store.Get(c.Request(), "flash-message")
	if err != nil {
		return FlashMessage{}
	}
	fm := session.Flashes()
	var flash FlashMessage
	if len(fm) > 0 {
		log.Info("FLASH MESSAGE ", fm[0])
		flash = fm[0].(FlashMessage)
	}
	if err := session.Save(c.Request(), c.Response()); err != nil {
		log.Fatal("ERROR GET FLASH MESSAGE ", err.Error())
	}
	return flash
}

func GetUserInfo(c echo.Context) UserInfo {
	session, err := Manager.store.Get(c.Request(), SessionId)
	if err != nil {
		return UserInfo{}
	}
	fm := session.Flashes()
	var data UserInfo
	if len(fm) > 0 {
		log.Info("USER INFO ", fm[0])
		data = fm[0].(UserInfo)
	}
	if err := session.Save(c.Request(), c.Response()); err != nil {
		log.Fatal("ERROR GET USER INFO ", err.Error())
	}
	return data
}

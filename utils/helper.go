package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GenerateTimeNow() string {
	//fetching current time
	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime := time.Now().In(loc).Format(time.RFC3339)
	//differnce between pastdate and current date
	return currentTime
}

func GetUserInfoFromContext(ctx echo.Context, db *gorm.DB) (userModel models.User, err error) {
	result, err := session.Manager.Get(ctx, session.SessionId)
	if err != nil {
		panic(err)
	}
	userInfo := result.(session.UserInfo)
	err = db.Model(models.User{}).Preload(clause.Associations).Where("user_id =?", userInfo.UserID).First(&userModel).Error
	return
}

func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

func FormatRupiah(prefix bool, amount float64) string {
	humanizeValue := humanize.CommafWithDigits(amount, 0)
	stringValue := strings.Replace(humanizeValue, ",", ".", -1)
	if prefix {
		return "Rp. " + stringValue
	}
	return stringValue
}

func GetAmountFromRupiah(rupiah string) int64 {
	re := regexp.MustCompile("[0-9]+")
	getNumber := re.FindAllString(rupiah, -1)
	amountString := strings.Join(getNumber, "")
	amount, err := strconv.ParseInt(amountString, 10, 64)
	if err != nil {
		return 0
	}
	return amount
}

func UploadImageV2(src string) (string, error) {
	if src == "" {
		return "", errors.New("image base64 is null")
	}
	decode, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	ext := http.DetectContentType(decode)
	if ext == "text/plain; charset=utf-8" {
		return "", errors.New("cannot get extension image")
	}
	res := strings.Split(ext, "image/")

	fName := uuid.New().String() + "." + res[1]
	r := bytes.NewReader(decode)
	dir, _ := os.Getwd()
	fileLocation := filepath.Join(dir, "assets/upload/avatars", fName)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, r); err != nil {
		return "", err
	}
	return "/assets/upload/avatars/" + fName, err
}


func ConfigTime(tm int)(date string) {
	t := time.Now()
	convert := t.AddDate(0,0,-tm)
	date = convert.Format("2006-01-02")
	return	date
}

func StrToUint64(str string) (uint64, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return uint64(i), err
}
package score_generator

import (
	"bufio"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"gopkg.in/resty.v1"
	"mime/multipart"
	"os"
	"strconv"
	"time"
)

type HttpClient struct {
	client *resty.Client
}

func NewUploadHTTPClient() *HttpClient {
	client := resty.New()
	//if os.Getenv("APP_ENV") != "production" {
	client.SetDebug(true)
	//}

	defaultTimeout, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT"))
	client.SetTimeout(time.Second * time.Duration(defaultTimeout))

	return &HttpClient{client: client}
}

func (hc *HttpClient) UploadFilePickle(name string,fileName multipart.File,replace string)  (int,interface{},error) {
	hc.client.SetDebug(true)
	bufferedReader:=bufio.NewReader(fileName)
	res, err := hc.client.R().
		SetFileReader("pickle_file" ,name,bufferedReader).
		SetFormData(map[string]string{"replace":replace}).
		SetHeader("Content-Type", "multipart/form-data").
		Post(os.Getenv("UPLOAD_PICKLE_BASE_URL")+os.Getenv("PICKLE_UPLOAD_URL"))
	log.Info(" Response Time From Upload Pickle : ", res.Time())
	var response interface{}
	_ = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return res.StatusCode(),response,err
	}
	return res.StatusCode(),response,nil
}
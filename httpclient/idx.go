package httpclient

import (
	"encoding/json"
	"github.com/kreditplus/scorepro/config/credential"
	"github.com/kreditplus/scorepro/models"
	"github.com/labstack/gommon/log"
	"gopkg.in/resty.v1"
	"os"
	"strconv"
	"time"
)

type IdxHTTPClient struct {
	client *resty.Client
}

func NewIdxHTTPClient() *IdxHTTPClient {
	client := resty.New()
	if os.Getenv("APP_ENV") != "production" {
		client.SetDebug(true)
	}

	defaultTimeout, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT"))
	client.SetTimeout(time.Second * time.Duration(defaultTimeout))

	return &IdxHTTPClient{client: client}
}

func (hc *IdxHTTPClient) ScoringPickleModeling(req interface{}, modelingPickle string,url string) (models.PickleResponseIDX, int, error) {

	res, err := hc.client.R().
		SetBody(map[string]string{"prospect_id" : "TEST-HIGH-XXXXXX"}).
		Post(url)
	if err != nil {
		log.Info("err : ",err)
	}
	var response models.PickleResponseIDX
	_ = json.Unmarshal(res.Body(), &response)

	return response, res.StatusCode(), err
}

func (hc *IdxHTTPClient) HitExperian(body interface{}, token string) (*models.ExperianScoreResp, int, error) {
	res, err := hc.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", token).
		SetBody(body).
		Post(credential.ExperianBaseUrl + os.Getenv("EXPERIAN_CHECK_CREDIT_URL"))

	if err != nil {
		log.Info("err : ",err)
	}

	var response models.ExperianScoreResp
	_ = json.Unmarshal(res.Body(), &response)

	statusCode := res.StatusCode()
	return &response, statusCode, err
}
package scorepro

import (
	"encoding/json"
	"github.com/kreditplus/scorepro/config/credential"
	"github.com/kreditplus/scorepro/utils/monitoring"
	"github.com/kreditplus/scorepro/models"
	"github.com/labstack/gommon/log"
	"gopkg.in/resty.v1"
	"os"
	"strconv"
	"time"
)

type HTTPClient struct {
	client *resty.Client
}

func NewHTTPClient() *HTTPClient {
	client := resty.New()
	if os.Getenv("APP_ENV") != "production" {
		client.SetDebug(true)
	}

	defaultTimeout, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT"))
	client.SetTimeout(time.Second * time.Duration(defaultTimeout))

	return &HTTPClient{client: client}
}

func (hc *HTTPClient) GetToken() (*models.TokenResponse, error) {
	res, err := hc.client.R().
		SetBasicAuth(credential.ExperianUsername, credential.ExperianPassword).
		Post(credential.ExperianBaseUrl + os.Getenv("EXPERIAN_TOKEN_URL"))
	if err != nil {
		log.Info("err : ", err)
		extra := map[string]interface{}{
			"message":err.Error(),
		}
		monitoring.SendToSentry(nil, extra, "EXPERIAN TOKEN")
	}
	var response models.TokenResponse
	_ = json.Unmarshal(res.Body(), &response)

	return &response, err
}

func (hc *HTTPClient) ScoringPickleModeling(req interface{}, modelingPickle string, url string) (PickleResponseIDX, int, error) {
	res, err := hc.client.R().
		SetHeader("Score-Generator", modelingPickle).
		SetBody(req).
		Post(url)
	if err != nil {
		log.Info("err : ", err)
		extra := map[string]interface{}{
			"message":err.Error(),
		}
		monitoring.SendToSentry(nil, extra, "SCORING PICKLE")
	}
	var response PickleResponseIDX
	_ = json.Unmarshal(res.Body(), &response)

	return response, res.StatusCode(), err
}

func (hc *HTTPClient) HitExperian(body interface{}, token string) (*models.ExperianScoreResp, int, error) {
	res, err := hc.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", token).
		SetBody(body).
		Post(credential.ExperianBaseUrl + os.Getenv("EXPERIAN_CHECK_CREDIT_URL"))

	if err != nil {
		log.Info("err : ", err)
		extra := map[string]interface{}{
			"message":err.Error(),
		}
		monitoring.SendToSentry(nil, extra, "SCORING EXPERIAN")
	}

	var response models.ExperianScoreResp
	_ = json.Unmarshal(res.Body(), &response)

	statusCode := res.StatusCode()
	return &response, statusCode, err
}

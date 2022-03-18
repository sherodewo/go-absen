package httpclient

import (
	"encoding/json"
	"github.com/kreditplus/scorepro/config/credential"
	"github.com/kreditplus/scorepro/constant"
	"github.com/kreditplus/scorepro/dto"
	"github.com/kreditplus/scorepro/models"
	"github.com/labstack/gommon/log"
	"github.com/sony/sonyflake"
	"gopkg.in/resty.v1"
	"net/http"
	"os"
	"strconv"
	"time"
)

type HTTPClient struct {
	client *resty.Client
}

func NewExperianHTTPClient() *HTTPClient {
	client := resty.New()
	if os.Getenv("APP_ENV") != "production" {
		client.SetDebug(true)
	}

	defaultTimeout, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT"))
	client.SetTimeout(time.Second * time.Duration(defaultTimeout))

	return &HTTPClient{client: client}
}

func (hc *HTTPClient) CheckCreditScore(phoneNumber string, token string) (*models.ExperianScoreResp, int, error) {
	high := string(constant.High)
	medium := string(constant.Medium)
	low := string(constant.Low)
	noRes := ""
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	flakeID := strconv.FormatUint(id, 10)
	TransId := phoneNumber + flakeID

	t := time.Now()
	format := t.Format("20060102150405")
	dateTime, _ := strconv.ParseFloat(format, 14)
	reqHeader := models.ReqHeader{
		ClientID:    credential.ExperianClientID,
		PartnerID:   credential.ExperianPartnerID,
		ProductID:   credential.ExperianProductID,
		TransID:     TransId,
		ReqDateTime: dateTime,
		MSISDN:      phoneNumber,
		ServiceName: "CHECK CREDIT",
		Model:       "BUREAU",
		Token:       token,
	}

	data := models.Data{OTP: os.Getenv("EXPERIAN_OTP")}

	reqBody := models.ReqBody{
		Data: data,
	}

	payload := models.PayloadExperian{
		ReqHeader: reqHeader,
		Body:      reqBody,
	}

	body := models.ExperianScore{PayLoad: payload}

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

	score := response.Payload.Body.CreditStatus.Score

	statusCode := res.StatusCode()

	scoreHigh,_ := strconv.Atoi(os.Getenv("SCORE_HIGH"))
	scoreLow,_ := strconv.Atoi(os.Getenv("SCORE_LOW"))

	if score == nil {
		response.Payload.Body.CreditStatus.Result = &noRes
		statusCode = 400
	} else if score != nil {
		response.Payload.Body.CreditStatus.Result = &medium
		if *score >= float64(scoreHigh) {
			response.Payload.Body.CreditStatus.Result = &high
		} else if *score <= float64(scoreLow) {
			response.Payload.Body.CreditStatus.Result = &low
		}
	}
	response.Payload.ResHeader.TransID = &TransId
	return &response, statusCode, err
}

func (hc *HTTPClient) GetToken() (*models.TokenResponse, error) {
	res, err := hc.client.R().
		SetBasicAuth(credential.ExperianUsername, credential.ExperianPassword).
		Post(credential.ExperianBaseUrl + os.Getenv("EXPERIAN_TOKEN_URL"))
	if err != nil {
		log.Info("err : ",err)
	}
	var response models.TokenResponse
	_ = json.Unmarshal(res.Body(), &response)

	return &response, err
}

func (hc *HTTPClient) ScoringPickle(req dto.PickleDto) (models.PickleResponse, int, error) {
	res, err := hc.client.R().
		SetBody(req).
		Post(os.Getenv("PICKLE_BASE_URL") + os.Getenv("PICKLE_SCORE_OFFLINE"))
	if err != nil {
		log.Info("err : ",err)
	}
	var response models.PickleResponse
	_ = json.Unmarshal(res.Body(), &response)

	return response, res.StatusCode(), err
}


func (hc *HTTPClient) ScoringPickleLimit(req dto.PickleLimitDto) (models.PickleLimitResponse, int, error) {
	res, err := hc.client.R().
		SetBody(req).
		Post(os.Getenv("PICKLE_LIMIT_BASE_URL") + os.Getenv("PICKLE_SCORE_LIMIT"))
	if err != nil {
		log.Info("err : ",err)
	}
	var response models.PickleLimitResponse
	_ = json.Unmarshal(res.Body(), &response)

	return response, res.StatusCode(), err
}

package service

import (
	"github.com/allegro/bigcache/v3"
	"github.com/kreditplus/scorepro/config/credential"
	"github.com/kreditplus/scorepro/constant"
	"github.com/kreditplus/scorepro/httpclient"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/repository"
	"github.com/sony/sonyflake"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type TelcoScoreService struct {
	TelcoScoreRepository repository.TelcoScoreRepository
}

func NewTelcoScoreService(repository repository.TelcoScoreRepository) *TelcoScoreService {
	t, _ := strconv.ParseInt(os.Getenv("EXPERIAN_TOKEN_EXPIRED"), 5, 5)
	expired := time.Duration(rand.Int63n(t)) * time.Hour
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(expired))
	return &TelcoScoreService{
		TelcoScoreRepository: repository,
	}
}

func (s *TelcoScoreService) FindAllTelcoScores() (*[]models.Experian, error) {
	data, err := s.TelcoScoreRepository.FindAll()
	return &data, err
}

func (s *TelcoScoreService) FindTelcoScoreById(id string) (*models.Experian, error) {
	data, err := s.TelcoScoreRepository.FindById(id)

	return &data, err
}

func (s *TelcoScoreService) FindTelcoScoreByPhoneNumber(phoneNumber string) (*models.Experian, error) {
	data, err := s.TelcoScoreRepository.FindByPhoneNumber(phoneNumber)

	return &data, err
}

func (s *TelcoScoreService) SaveTelcoScore(req models.Experian) (*models.Experian, error) {
	entity := models.Experian{
		ProspectID:     req.ProspectID,
		Result:         req.Result,
		Status:         req.Status,
		PhoneNumber:    req.PhoneNumber,
		Score:          req.Score,
		Experian:       req.Experian,
		ScorePro:       req.ScorePro,
		ExperianScore:  req.ExperianScore,
		ExperianResult: req.ExperianResult,
		InternalScore:  req.InternalScore,
		InternalResult: req.InternalResult,
		TransID:        req.TransID,
		Type:           req.Type,
		CreatedAt:      time.Time{},
	}
	data, err := s.TelcoScoreRepository.Save(entity)
	return &data, err
}

func (s *TelcoScoreService) UpdateTelcoScore(id string, req models.Experian) (*models.Experian, error) {
	var entity models.Experian
	data, err := s.TelcoScoreRepository.Update(entity)

	return &data, err
}

func (s *TelcoScoreService) DeleteTelcoScore(id string) error {
	entity := models.Experian{
		ExperianID: id,
	}
	err := s.TelcoScoreRepository.Delete(entity)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *TelcoScoreService) QueryDatatable(searchValue string, orderType string, orderBy string, limit int, offset int) (
	recordTotal int64, recordFiltered int64, data []models.Experian, err error) {
	recordTotal, err = s.TelcoScoreRepository.Count()

	if searchValue != "" {
		recordFiltered, err = s.TelcoScoreRepository.CountWhere("or", map[string]interface{}{
			"ProspectID LIKE ?": "%" + searchValue + "%",
			"result LIKE ?":     "%" + searchValue + "%",
		})

		data, err = s.TelcoScoreRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
			"ProspectID LIKE ?": "%" + searchValue + "%",
			"result LIKE ?":     "%" + searchValue + "%",
		})
		return recordTotal, recordFiltered, data, err
	}

	recordFiltered, err = s.TelcoScoreRepository.CountWhere("or", map[string]interface{}{
		"1 =?": 1,
	})

	data, err = s.TelcoScoreRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
		"1= ?": 1,
	})
	return recordTotal, recordFiltered, data, err
}

func (s *TelcoScoreService) IndosatValidate(phoneNumber string) (models.AppConfig, bool) {
	var data models.AppConfig
	var active int
	val, err := cache.Get("indosat-validation")
	if err != nil {
		res, _ := s.TelcoScoreRepository.IndosatValidate()
		data = res
		isActive := strconv.Itoa(res.IsActive)
		cache.Set("indosat-active", []byte(isActive))
		cache.Set("indosat-validation", []byte(res.Value))
		val, _ = cache.Get("indosat-validation")
	}
	isActive, err := cache.Get("indosat-active")
	active, _ = strconv.Atoi(string(isActive))
	data.Value = string(val)
	data.IsActive = active

	arrApp := strings.Split(data.Value, ",")

	validator := contains(arrApp, phoneNumber[:5])

	return data, validator
}

func (s *TelcoScoreService) CheckConfig() (models.AppConfig, bool) {
	data, valid := s.TelcoScoreRepository.CheckConfig()

	return data, valid
}

func (s *TelcoScoreService) GetDbInstance() *gorm.DB {
	return s.TelcoScoreRepository.DbInstance()
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func (s *TelcoScoreService) GetExperianDummy(PhoneNumber string) (*models.Dummy, error) {
	data, err := s.TelcoScoreRepository.GetExperianDummy(PhoneNumber)
	return &data, err
}

func (s *TelcoScoreService) GetPickleModeling(SupplierID string) (*models.SupplierModel, error) {
	data, err := s.TelcoScoreRepository.GetPickleModeling(SupplierID)
	return &data, err
}

func (s *TelcoScoreService) GetPickleModelingOther(ZipCode string) (*models.SupplierOther, error) {
	data, err := s.TelcoScoreRepository.GetPickleModelingOther(ZipCode)
	return &data, err
}

func (s *TelcoScoreService) FinalScoring(isIndosat bool, experianResult string, internalResult string) (*models.ExperianScoring, error) {
	indosat := "NO"
	if isIndosat {
		indosat = "YES"
	} else if &experianResult == nil {
		experianResult = constant.NoHit
	} else if &internalResult == nil {
		internalResult = constant.NoHit
	}
	data, err := s.TelcoScoreRepository.FinalScoring(indosat, experianResult, internalResult)
	return &data, err
}

func (s *TelcoScoreService) HitExperian(phoneNumber string, token string) (*models.ExperianScoreResp, int, error) {
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

	httpClient := httpclient.NewIdxHTTPClient()
	response, statusCode, err := httpClient.HitExperian(body, token)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}

	score := response.Payload.Body.CreditStatus.Score

	scoreHigh, _ := strconv.Atoi(os.Getenv("SCORE_HIGH"))
	scoreLow, _ := strconv.Atoi(os.Getenv("SCORE_LOW"))

	if score == nil {
		response.Payload.Body.CreditStatus.Result = &noRes
		statusCode = http.StatusBadRequest
	} else if score != nil {
		response.Payload.Body.CreditStatus.Result = &medium
		if *score >= float64(scoreHigh) {
			response.Payload.Body.CreditStatus.Result = &high
		} else if *score <= float64(scoreLow) {
			response.Payload.Body.CreditStatus.Result = &low
		}
	}
	response.Payload.ResHeader.TransID = &TransId
	return response, statusCode, err
}

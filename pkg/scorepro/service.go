package scorepro

import (
	"encoding/json"
	"errors"
	"github.com/allegro/bigcache/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/kreditplus/scorepro/config/credential"
	"github.com/kreditplus/scorepro/constant"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/pkg/data_inquiry"
	"github.com/kreditplus/scorepro/pkg/experian_inquiry"
	"github.com/kreditplus/scorepro/pkg/metric_combination"
	"github.com/kreditplus/scorepro/pkg/score_generator"
	"github.com/kreditplus/scorepro/pkg/scorepro_requestor"
	"github.com/kreditplus/scorepro/pkg/scorepro_score"
	"github.com/labstack/echo/v4"
	"github.com/sony/sonyflake"
	"gopkg.in/go-playground/validator.v9"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	cache      *bigcache.BigCache
	httpclient = NewHTTPClient()
)

type ScoreproService struct {
	ScoreproRepository ScoreproRepository
}

func NewScoreproService(repository ScoreproRepository) *ScoreproService {
	t, _ := strconv.ParseInt(os.Getenv("EXPERIAN_TOKEN_EXPIRED"), 5, 5)
	expired := time.Duration(rand.Int63n(t)) * time.Hour
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(expired))
	return &ScoreproService{
		ScoreproRepository: repository,
	}
}

func (s *ScoreproService) CheckRequestor(id string) (*scorepro_requestor.ScoreproRequestor, error) {
	data, err := s.ScoreproRepository.CheckRequestor(id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *ScoreproService) ModelingPickle(id string) (*score_generator.ScoreGenerator, error) {
	data, err := s.ScoreproRepository.GetScoreModelingPickle(id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *ScoreproService) StoreDataInquiry(req PickleModelingDto, IDScoreGenerator string, ScoreproScoreID string, ScoreproRequestor string, id string, metricCombinationID *string) (*data_inquiry.DataInquiry, error) {
	request := data_inquiry.DataInquiry{
		ID:                id,
		ProspectID:        req.ProspectID,
		StatusKonsumen:    req.StatusKonsumen,
		MobilePhone:       req.PhoneNumber,
		SupplierID:        &req.SupplierID,
		Journey:           req.Journey,
		MetricCombination: metricCombinationID,
		TypeModel:         IDScoreGenerator,
		ScoreproScore:     ScoreproScoreID,
		ScoreproRequestor: ScoreproRequestor,
		CreatedAt:         time.Now(),
	}

	data, err := s.ScoreproRepository.StoreDataInquiry(request)
	if err != nil {
		return nil, err
	}

	return &data, err
}

func (s *ScoreproService) HitScoringPickle(req PickleModelingDto, ScoreGenerator score_generator.ScoreGenerator, transactionID string, ctx echo.Context) (*PickleResponseIDX, int, error, *[]models.ErrorValidation, *string) {
	// Check On Off
	_, isOff := s.CheckAppConf(constant.ConfigInternalServiceOn)
	if !isOff {
		//Get Resp Error
		respErr, isOff := s.CheckAppConf(constant.ConfigResponseServiceOff)
		if isOff {
			data := strings.Split(respErr.Value, ",")
			status, _ := strconv.Atoi(data[1])
			return nil, status, errors.New(data[0]), nil, nil
		}

	}

	var hitPickle interface{}
	var validationErrors []models.ErrorValidation

	if ScoreGenerator.ID == constant.ModelingModernMarket {
		var data PickleMM
		dataPickle, _ := json.Marshal(req.Data)
		data.TransactionID = transactionID

		_ = json.Unmarshal(dataPickle, &data)
		if err := ctx.Validate(&data); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				validationErrors = models.WrapValidationApiErrors(errs)
			}
			return nil, 0, nil, &validationErrors, nil
		}

		// Zip Code 2 Digit
		strZip := strconv.Itoa(data.ZipCode)
		zipCode, _ := strconv.Atoi(strZip[:2])
		data.ZipCode = zipCode

		//First Four Of Cell
		data.FirstFourOfCell = data.FirstFourOfCell[:4]

		hitPickle = data
	} else if ScoreGenerator.ID == constant.ModelingJaboJabar {
		var data PickleJJ
		dataPickle, _ := json.Marshal(req.Data)
		_ = json.Unmarshal(dataPickle, &data)
		data.TransactionID = transactionID

		if err := ctx.Validate(&data); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				validationErrors = models.WrapValidationApiErrors(errs)
			}
			return nil, 0, nil, &validationErrors, nil
		}

		// Zip Code 2 Digit
		strZip := strconv.Itoa(data.ZipCode)
		zipCode, _ := strconv.Atoi(strZip[:2])
		data.ZipCode = zipCode

		//First Four Of Cell
		data.FirstFourOfCell = data.FirstFourOfCell[:4]

		hitPickle = data
	} else if ScoreGenerator.ID == constant.ModelingTraditional {
		var data PickleOT
		dataPickle, _ := json.Marshal(req.Data)
		_ = json.Unmarshal(dataPickle, &data)
		data.TransactionID = transactionID

		if err := ctx.Validate(&data); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				validationErrors = models.WrapValidationApiErrors(errs)
			}
			return nil, 0, nil, &validationErrors, nil
		}

		// Zip Code 2 Digit
		strZip := strconv.Itoa(data.ZipCode)
		zipCode, _ := strconv.Atoi(strZip[:2])
		data.ZipCode = zipCode

		//First Four Of Cell
		data.FirstFourOfCell = data.FirstFourOfCell[:4]

		hitPickle = data
	}

	datas, status, err := httpclient.ScoringPickleModeling(hitPickle, ScoreGenerator.ID, ScoreGenerator.Endpoint)
	if err != nil {
		return nil, http.StatusBadGateway, err, nil, nil
	}

	scoreproScore, err := s.ScoreproRepository.ScoreproScore(datas.Data.Score, constant.ScoringInternalWgOff)
	if err != nil {
		return nil, http.StatusBadGateway, err, nil, nil
	}

	datas.Data.Result = scoreproScore.Result

	return &datas, status, nil, nil, &scoreproScore.ID
}

// HitScoringPickleKmb Pickle For KMB
func (s *ScoreproService) HitScoringPickleKmb(req PickleModelingDto, ScoreGenerator score_generator.ScoreGenerator, transactionID string, ctx echo.Context) (*PickleResponseIDX, int, error, *[]models.ErrorValidation, *string) {
	// Check On Off
	_, isOff := s.CheckAppConf(constant.ConfigInternalServiceOn)
	if !isOff {
		//Get Resp Error
		respErr, isOff := s.CheckAppConf(constant.ConfigResponseServiceOff)
		if isOff {
			data := strings.Split(respErr.Value, ",")
			status, _ := strconv.Atoi(data[1])
			return nil, status, errors.New(data[0]), nil, nil
		}

	}

	var hitPickle interface{}
	var validationErrors []models.ErrorValidation

	if ScoreGenerator.ID == constant.Modeling2WJabo {
		var data Pickle2WJabo
		dataPickle, _ := json.Marshal(req.Data)
		data.TransactionID = transactionID

		_ = json.Unmarshal(dataPickle, &data)
		if err := ctx.Validate(&data); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				validationErrors = models.WrapValidationApiErrors(errs)
			}
			return nil, 0, nil, &validationErrors, nil
		}

		// Zip Code 2 Digit
		strZip := strconv.Itoa(data.ZipCode)
		zipCode, _ := strconv.Atoi(strZip[:2])
		data.ZipCode = zipCode

		hitPickle = data
	} else if ScoreGenerator.ID == constant.Modeling2WOther {
		var data Pickle2WOther
		dataPickle, _ := json.Marshal(req.Data)
		_ = json.Unmarshal(dataPickle, &data)
		data.TransactionID = transactionID

		if err := ctx.Validate(&data); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				validationErrors = models.WrapValidationApiErrors(errs)
			}
			return nil, 0, nil, &validationErrors, nil
		}

		// Zip Code 2 Digit
		strZip := strconv.Itoa(data.ZipCode)
		zipCode, _ := strconv.Atoi(strZip[:2])
		data.ZipCode = zipCode

		//First Four Of Cell
		data.FirstFourOfCell = data.FirstFourOfCell[:4]

		hitPickle = data
	}

	datas, status, err := httpclient.ScoringPickleModeling(hitPickle, ScoreGenerator.ID, ScoreGenerator.Endpoint)
	if err != nil {
		return nil, http.StatusBadGateway, err, nil, nil
	}

	scoreproScore, err := s.ScoreproRepository.ScoreproScore(datas.Data.Score, constant.ScoringInternalKmbOff)
	if err != nil {
		return nil, http.StatusBadGateway, err, nil, nil
	}

	datas.Data.Result = scoreproScore.Result

	return &datas, status, nil, nil, &scoreproScore.ID
}


// HitScoringPickleWg Pickle For WG Online
func (s *ScoreproService) HitScoringPickleWg(req PickleModelingDto, ScoreGenerator score_generator.ScoreGenerator, transactionID string, ctx echo.Context, transactionType string) (*PickleResponseIDX, int, error, *[]models.ErrorValidation, *scorepro_score.ScoreproScore) {
	// Check On Off
	_, isOff := s.CheckAppConf(constant.ConfigInternalServiceOn)
	if !isOff {
		//Get Resp Error
		respErr, isOff := s.CheckAppConf(constant.ConfigResponseServiceOff)
		if isOff {
			data := strings.Split(respErr.Value, ",")
			status, _ := strconv.Atoi(data[1])
			return nil, status, errors.New(data[0]), nil, nil
		}

	}

	var hitPickle interface{}
	var validationErrors []models.ErrorValidation
	if transactionType == constant.TransactionLimit {
		if ScoreGenerator.ID == constant.ModelingWgJabo {
			var data ModelingPickleWgJabo
			dataPickle, _ := json.Marshal(req.Data)
			_ = json.Unmarshal(dataPickle, &data)
			data.TransactionID = transactionID

			if err := ctx.Validate(&data); err != nil {
				if errs, ok := err.(validator.ValidationErrors); ok {
					validationErrors = models.WrapValidationApiErrors(errs)
				}
				return nil, 0, nil, &validationErrors, nil
			}

			// Zip Code 2 Digit
			strZip := strconv.Itoa(data.ZipCode)
			zipCode, _ := strconv.Atoi(strZip[:2])
			data.ZipCode = zipCode

			hitPickle = data
		} else if ScoreGenerator.ID == constant.ModelingWgOther {
			var data ModelingPickleWgOther
			dataPickle, _ := json.Marshal(req.Data)
			_ = json.Unmarshal(dataPickle, &data)
			data.TransactionID = transactionID

			if err := ctx.Validate(&data); err != nil {
				if errs, ok := err.(validator.ValidationErrors); ok {
					validationErrors = models.WrapValidationApiErrors(errs)
				}
				return nil, 0, nil, &validationErrors, nil
			}

			// Zip Code 2 Digit
			strZip := strconv.Itoa(data.ZipCode)
			zipCode, _ := strconv.Atoi(strZip[:2])
			data.ZipCode = zipCode

			//First Four Of Cell
			data.FirstFourOfCell = data.FirstFourOfCell[:4]

			hitPickle = data
		}

		datas, status, err := httpclient.ScoringPickleModeling(hitPickle, ScoreGenerator.ID, ScoreGenerator.Endpoint)
		if err != nil {
			return nil, http.StatusBadGateway, err, nil, nil
		}

		scoreproScore, err := s.ScoreproRepository.ScoreproScoreWithTrxType(datas.Data.Score, constant.ScoringInternalWgOnl, constant.TransactionLimit)

		if err != nil {
			return nil, http.StatusBadGateway, err, nil, nil
		}

		datas.Data.Result = scoreproScore.Result
		return &datas, status, nil, nil, &scoreproScore

	} else if transactionType == constant.TransactionProductLimit {
		var data ModelingPickleWgProductLimit
		dataPickle, _ := json.Marshal(req.Data)
		_ = json.Unmarshal(dataPickle, &data)
		data.TransactionID = transactionID

		if err := ctx.Validate(&data); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				validationErrors = models.WrapValidationApiErrors(errs)
			}
			return nil, 0, nil, &validationErrors, nil
		}

		// Zip Code 2 Digit
		strZip := strconv.Itoa(data.ZipCode)
		zipCode, _ := strconv.Atoi(strZip[:2])
		data.ZipCode = zipCode

		hitPickle = data

		datas, status, err := httpclient.ScoringPickleModeling(hitPickle, ScoreGenerator.ID, ScoreGenerator.Endpoint)
		if err != nil {
			return nil, http.StatusBadGateway, err, nil, nil
		}

		scoreproScore, err := s.ScoreproRepository.ScoreproScoreWithTrxType(datas.Data.Score, constant.ScoringInternalWgOnl, constant.TransactionProductLimit)

		if err != nil {
			return nil, http.StatusBadGateway, err, nil, nil
		}

		datas.Data.Result = scoreproScore.Result
		return &datas, status, nil, nil, &scoreproScore
	}

	return nil, http.StatusServiceUnavailable, errors.New(constant.UNDEFINED_TRANSACTION_TYPE), nil, nil
}

func (s *ScoreproService) IndosatValidate(phoneNumber string) (models.AppConfig, bool) {
	var data models.AppConfig
	var active int
	val, err := cache.Get("indosat-validation")
	if err != nil {
		res, _ := s.ScoreproRepository.AppConf("indosat")
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

func (s *ScoreproService) MetricCombintaion(isIndosat bool, experian *models.ExperianScoreResp, internal *PickleResponseIDX,lob string) (*metric_combination.MetricCombination, *float64, error) {
	var score *float64
	var data metric_combination.MetricCombination
	var err error
	indosat := "YES"
	if !isIndosat {
		indosat = "NO"
	}
	if internal != nil && experian != nil {
		data, err = s.ScoreproRepository.MetricCombination(indosat, *experian.Payload.Body.CreditStatus.Result, internal.Data.Result, lob)
		if err != nil {
			return nil, nil, err
		}
	} else if internal == nil && experian != nil {
		data, err = s.ScoreproRepository.MetricCombination(indosat, *experian.Payload.Body.CreditStatus.Result, constant.NoHit, lob)
		if err != nil {
			return nil, nil, err
		}

	} else if experian == nil && internal != nil {
		data, err = s.ScoreproRepository.MetricCombination(indosat, constant.NoHit, internal.Data.Result, lob)
		if err != nil {
			return nil, nil, err
		}
	} else {
		data, err = s.ScoreproRepository.MetricCombination(indosat, constant.NoHit, constant.NoHit, lob)
		if err != nil {
			return nil, nil, err
		}
	}

	if data.FinalScore == constant.ASSTEL_HIGH || data.FinalScore == constant.ASSTEL_MEDIUM || data.FinalScore == constant.ASSTEL_LOW {
		score = experian.Payload.Body.CreditStatus.Score
	} else if data.FinalScore == constant.ASS_HIGH || data.FinalScore == constant.ASS_MEDIUM || data.FinalScore == constant.ASS_LOW {
		score = &internal.Data.Score
	} else {
		score = nil
	}

	return &data, score, nil
}

func (s *ScoreproService) UseDummy() (*models.AppConfig, bool) {
	data, err := s.ScoreproRepository.AppConf("dummy_experian")
	if err != nil {
		return nil, false
	}

	if data.IsActive == 1 {
		return &data, true
	}
	return nil, false
}

func (s *ScoreproService) GetExperianDummy(PhoneNumber string) (*models.Dummy, error) {
	data, err := s.ScoreproRepository.GetExperianDummy(PhoneNumber)
	return &data, err
}

func (s *ScoreproService) SaveTelcoScore(req PickleModelingDto, response ScoreproResponse, experianResp *models.ExperianScoreResp, internalResp *PickleResponseIDX, id string) (*experian_inquiry.Experian, error) {
	expData, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(experianResp)
	spData, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(internalResp)
	experianRes := string(expData)
	internalRes := string(spData)

	var entity experian_inquiry.Experian
	if experianResp != nil && internalResp != nil {
		entity = experian_inquiry.Experian{
			ProspectID:     &req.ProspectID,
			Result:         &response.Result,
			Status:         &response.Status,
			PhoneNumber:    &req.PhoneNumber,
			Score:          response.Score,
			Experian:       &experianRes,
			ScorePro:       &internalRes,
			ExperianScore:  experianResp.Payload.Body.CreditStatus.Score,
			ExperianResult: experianResp.Payload.Body.CreditStatus.Result,
			InternalScore:  &internalResp.Data.Score,
			InternalResult: &internalResp.Data.Result,
			TransID:        experianResp.Payload.ResHeader.TransID,
			Type:           &req.Journey,
			DataInquiry:    id,
			CreatedAt:      time.Now(),
		}
	} else if experianResp == nil && internalResp != nil {
		entity = experian_inquiry.Experian{
			ProspectID:     &req.ProspectID,
			Result:         &response.Result,
			Status:         &response.Status,
			PhoneNumber:    &req.PhoneNumber,
			Score:          response.Score,
			ScorePro:       &internalRes,
			InternalScore:  &internalResp.Data.Score,
			InternalResult: &internalResp.Data.Result,
			Type:           &req.Journey,
			DataInquiry:    id,
			CreatedAt:      time.Now(),
		}
	} else if internalResp == nil && experianResp != nil {
		entity = experian_inquiry.Experian{
			ProspectID:     &req.ProspectID,
			Result:         &response.Result,
			TransID:        experianResp.Payload.ResHeader.TransID,
			Status:         &response.Status,
			PhoneNumber:    &req.PhoneNumber,
			Score:          response.Score,
			Experian:       &experianRes,
			ExperianScore:  experianResp.Payload.Body.CreditStatus.Score,
			ExperianResult: experianResp.Payload.Body.CreditStatus.Result,
			Type:           &req.Journey,
			DataInquiry:    id,
			CreatedAt:      time.Now(),
		}
	} else {
		entity = experian_inquiry.Experian{
			ProspectID:  &req.ProspectID,
			Result:      &response.Result,
			Status:      &response.Status,
			PhoneNumber: &req.PhoneNumber,
			Score:       response.Score,
			Type:        &req.Journey,
			DataInquiry: id,
			CreatedAt:   time.Now(),
		}
	}

	data, err := s.ScoreproRepository.Save(&entity)
	return data, err
}

func (s *ScoreproService) HitExperian(phoneNumber string, token string, prefix string) (*models.ExperianScoreResp, int, *string, error) {

	// Check On Off
	_, isOff := s.CheckAppConf(constant.ConfigExperianServiceOn)
	if !isOff {
		//Get Resp Error
		respErr, isOff := s.CheckAppConf(constant.ConfigResponseServiceOff)
		if isOff {
			data := strings.Split(respErr.Value, ",")
			status, _ := strconv.Atoi(data[1])
			return nil, status, nil, errors.New(data[0])
		}

	}

	var idScore *string

	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return nil, http.StatusInternalServerError, nil, err
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
	resp, status, err := httpclient.HitExperian(body, token)

	if err != nil {
		return nil, http.StatusBadGateway, nil, err
	}

	if resp.Payload.Body.CreditStatus.Score != nil {
		// ScoreproScore
		scoreproScore, err := s.ScoreproRepository.ScoreproScore(*resp.Payload.Body.CreditStatus.Score,prefix)
		if err != nil {
			return nil, http.StatusBadGateway, nil, err
		}
		resp.Payload.Body.CreditStatus.Result = &scoreproScore.Result
		idScore = &scoreproScore.ID
	}else {
		status = http.StatusBadGateway
		resp = nil
		idScore  = nil
	}

	return resp, status, idScore, nil
}

func (s *ScoreproService) QueryDatatable(searchValue string, orderType string, orderBy string, limit int, offset int) (
	recordTotal int64, recordFiltered int64, data []experian_inquiry.Experian, err error) {
	recordTotal, err = s.ScoreproRepository.Count()

	if searchValue != "" {
		recordFiltered, err = s.ScoreproRepository.CountWhere("or", map[string]interface{}{
			"ProspectID LIKE ?": "%" + searchValue + "%",
			"result LIKE ?":     "%" + searchValue + "%",
		})

		data, err = s.ScoreproRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
			"ProspectID LIKE ?": "%" + searchValue + "%",
			"result LIKE ?":     "%" + searchValue + "%",
		})
		return recordTotal, recordFiltered, data, err
	}

	recordFiltered, err = s.ScoreproRepository.CountWhere("or", map[string]interface{}{
		"1 =?": 1,
	})

	data, err = s.ScoreproRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
		"1= ?": 1,
	})
	return recordTotal, recordFiltered, data, err
}

func (s *ScoreproService) FindTelcoScoreById(id string) (*experian_inquiry.Experian, error) {
	data, err := s.ScoreproRepository.FindById(id)

	return &data, err
}

func (s *ScoreproService) CheckAppConf(name string) (*models.AppConfig, bool) {
	data, err := s.ScoreproRepository.AppConf(name)
	if err != nil {
		return nil, false
	}

	if data.IsActive == 1 {
		return &data, true
	}
	return nil, false
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

package controllers

import (
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/kreditplus/scorepro/constant"
	"github.com/kreditplus/scorepro/dto"
	"github.com/kreditplus/scorepro/httpclient"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/service"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var cache *bigcache.BigCache

type TelcoScoreController struct {
	Controller
	BaseBackendController
	service *service.TelcoScoreService
}

func NewTelcoScoreController(service *service.TelcoScoreService) TelcoScoreController {
	t, _ := strconv.ParseInt(os.Getenv("EXPERIAN_TOKEN_EXPIRED"), 5, 5)
	expired := time.Duration(rand.Int63n(t)) * time.Hour
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(expired))
	return TelcoScoreController{
		BaseBackendController: BaseBackendController{
			Menu:        "Dashboard",
			BreadCrumbs: []map[string]interface{}{},
		},
		service: service,
	}
}

func (c *TelcoScoreController) Index(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "List Data",
		"link": "/scorepro/admin/dashboard",
	}
	return Render(ctx, "Dashboard", "telcoscore/index", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}

func (c *TelcoScoreController) List(ctx echo.Context) error {

	draw, err := strconv.Atoi(ctx.Request().URL.Query().Get("draw"))
	search := ctx.Request().URL.Query().Get("search[value]")
	start, err := strconv.Atoi(ctx.Request().URL.Query().Get("start"))
	length, err := strconv.Atoi(ctx.Request().URL.Query().Get("length"))
	order, err := strconv.Atoi(ctx.Request().URL.Query().Get("order[0][column]"))
	orderName := ctx.Request().URL.Query().Get("columns[" + strconv.Itoa(order) + "][name]")
	orderAscDesc := ctx.Request().URL.Query().Get("order[0][dir]")

	recordTotal, recordFiltered, data, err := c.service.QueryDatatable(search, orderAscDesc, orderName, length, start)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var action string
	var createdAt string
	listOfData := make([]map[string]interface{}, len(data))
	for k, v := range data {

		action = `<button class="btn btn-success btn-bold btn-upper" href="javascript:;" onclick="Detail('` + v.ExperianID + `')" data-toggle="modal" data-target="#detail" data-id="` + v.ExperianID + `" style="text-decoration: none;font-weight: 400; color: white;">
                      Detail </button>`
		time := v.CreatedAt
		createdAt = time.Format("2006-01-02T15:04:05+07:00")
		listOfData[k] = map[string]interface{}{
			"experian_id":  v.ExperianID,
			"phone_number": v.PhoneNumber,
			"ProspectID":   v.ProspectID,
			"result":       v.Result,
			"status":       v.Status,
			"score":        v.Score,
			"created_at":   createdAt,
			"action":       action,
		}
	}

	result := models.ResponseDatatable{
		Draw:            draw,
		RecordsTotal:    recordTotal,
		RecordsFiltered: recordFiltered,
		Data:            listOfData,
	}
	return ctx.JSON(http.StatusOK, &result)
}

func (c *TelcoScoreController) Experian(ctx echo.Context) error {
	var req dto.ExperianDto

	if err := ctx.Bind(&req); err != nil {
		return c.InternalServerError(ctx, err)
	}
	var validationErrors []models.ErrorValidation

	if err := ctx.Validate(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationApiErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}

	httpClient := httpclient.NewExperianHTTPClient()

	//GET TOKEN FROM INDOSAT
	token, err := httpClient.GetToken()
	if err != nil {
		return c.InternalServerError(ctx, err)
	}

	//CHECKING SCORE TO EXPERIAN
	resp, _, err := httpClient.CheckCreditScore(req.PhoneNumber, token.AccessToken)
	if err != nil {
		return c.InternalServerError(ctx, err)
	}

	return c.Ok(ctx, resp)
}

func (c *TelcoScoreController) GetToken(ctx echo.Context) error {
	httpClient := httpclient.NewExperianHTTPClient()
	resp, err := httpClient.GetToken()
	if err != nil {
		return c.InternalServerError(ctx, err)
	}

	return c.Ok(ctx, resp)
}

func (c *TelcoScoreController) InternalScoring(ctx echo.Context) error {
	var req dto.PickleDto

	if err := ctx.Bind(&req); err != nil {
		return c.InternalServerError(ctx, err)
	}
	if err := ctx.Validate(&req); err != nil {
		return c.InternalServerError(ctx, err)
	}

	httpClient := httpclient.NewExperianHTTPClient()
	resp, _, err := httpClient.ScoringPickle(req)
	if err != nil {
		return c.InternalServerError(ctx, err)
	}

	return c.Ok(ctx, resp)
}

// Integrator godoc
// @Description Los Scorepro
// @Tags SCORE PRO
// @Produce json
// @Param scs query string false "scoring type by scs"
// @Param phoneNumber path string true "Mobile Phone"
// @Param body body dto.PickleDto true "Body payload"
// @Success 200 {object} ApiResponse{data=models.CreditScoreResp}
// @Failure 400 {object} ApiResponse{}
// @Failure 500 {object} ApiResponse{}
// @Router /score/credit/{phoneNumber} [post]
func (c *TelcoScoreController) CreditScore(ctx echo.Context) error {
	var experian dto.ExperianDto
	var req dto.PickleDto
	var scoringStatusInternal int
	var ExperianScore *float64
	var ExperianResult *string
	var TransID *string
	var isIndosatPrefix bool
	var isScoringInternalActive bool
	var ExperianStatus int
	var token string
	var scorePro models.PickleResponse
	var response models.CreditScoreResp
	var ExperianResp *models.ExperianScoreResp
	var exData []byte
	var ExperianStatusIssue string

	high := string(constant.High)
	medium := string(constant.Medium)
	low := string(constant.Low)
	med2 := string(constant.Med2nd)

	param := ctx.Param("phoneNumber")
	tipe := ctx.QueryParam("type")
	Tipe := strings.ToUpper(tipe)
	if param[:1] == "0" {
		res := strings.Replace(param, "0", "62", 1)
		param = res
	} else if param[:1] != "0" && param[:2] != "62" {
		return c.BadRequestWithSpecificFieldResponses(ctx, "error validation", []string{"phone_number", "start=08,628"})
	}

	experian.PhoneNumber = param

	if err := ctx.Bind(&req); err != nil {
		return c.BadRequest(ctx, err)
	}

	//Validate Phone Number
	var validationErrors []models.ErrorValidation
	if err := ctx.Validate(&experian); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationApiErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}

	//Validate Request
	if err := ctx.Validate(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationApiErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}

	if len(req.Zip) > 3 {
		req.Zip = req.Zip[:3]
	}

	_, isIndosatPrefix = c.service.IndosatValidate(experian.PhoneNumber)

	httpClient := httpclient.NewExperianHTTPClient()
	response = models.CreditScoreResp{
		ProspectID:  req.OrderID,
		PhoneNumber: param,
	}
	//Check IF Dummy is true
	var experianDummy *models.Dummy
	dummy := os.Getenv("EXPERIAN_DUMMY")
	if dummy == "true" && strings.ToUpper(os.Getenv("APP_ENV")) != constant.EnvProduction{
		res, err := c.service.GetExperianDummy(experian.PhoneNumber)
		if res == nil {
			return c.NotFound(ctx, err)
		}
		experianDummy = res
	} else {
		//SET BIG CACHE Token Indosat
		Itoken, err := cache.Get("indosat-token")
		token = string(Itoken)
		if Itoken == nil || err != nil {
			res, err := httpClient.GetToken()
			if err != nil {
				log.Info("GetToken Error: ", err)

				// By pass Experian, next process to API Pickle
				isIndosatPrefix = false
			}
			cache.Set("indosat-token", []byte(res.AccessToken))
			Itoken, _ := cache.Get("indosat-token")
			token = string(Itoken)
		}
	}

	if *req.SpsInc > 0 {
		*req.SpsInc = 1
	}
	if *req.VarInc > 0 {
		*req.VarInc = 1
	}
	//Type SRS
	if Tipe == constant.TypeSCS1 || Tipe == "" {
		Tipe = constant.TypeSCS1

		isScoringInternalActive = true
		//SET CACHE USE_SCORING
		useScoring, err := cache.Get("use-scoring")
		if useScoring == nil || err != nil {
			//Check configurable hit internal scoring or Not
			data, val := c.service.CheckConfig()
			isActive := strconv.Itoa(data.IsActive)
			isScoringInternalActive = val
			cache.Set("use-scoring", []byte(isActive))
			useScoring, err = cache.Get("use-scoring")
		}

		if string(useScoring) != "1" {
			isScoringInternalActive = false
		}

		var resp *models.ExperianScoreResp
		if isScoringInternalActive && isIndosatPrefix {
			//INTERNAL ON && EXPERIAN ON
			//IF DUMMY
			if experianDummy != nil {
				_ = json.Unmarshal([]byte(experianDummy.ExperianResponse), &resp)
				ExperianStatus = http.StatusOK
			} else {
				//CHECKING SCORE TO EXPERIAN
				creditScore, status, err := httpClient.CheckCreditScore(experian.PhoneNumber, token)
				if err != nil {
					log.Info("CheckExperian : ", err)
				}
				resp = creditScore
				ExperianStatus = status
			}
			//CHECK INTERNAL SCORING
			res, ScoringStatus, err := httpClient.ScoringPickle(req)
			if err != nil {
				log.Error("Scoring Internal : ", err)
			}

			if resp != nil && ExperianStatus == http.StatusOK {
				response.Score = resp.Payload.Body.CreditStatus.Score
				response.Result = resp.Payload.Body.CreditStatus.Result
				response.Status = constant.Experian

				if *resp.Payload.Body.CreditStatus.Result == medium && res.Result == constant.High {
					response.Score = res.Score
					response.Result = &res.Result
					response.Status = constant.ScoringInternal
				} else if *resp.Payload.Body.CreditStatus.Result == medium && res.Result != constant.High {
					response.Score = resp.Payload.Body.CreditStatus.Score
					response.Result = &low
					response.Status = constant.Experian
				}

				if ExperianStatus != http.StatusOK && ScoringStatus == http.StatusOK {
					response.Score = res.Score
					response.Result = &res.Result
					response.Status = constant.ScoringInternal
					//EXPERIAN ISSUE
					ExperianStatusIssue = constant.ExperianIssue
				} else if ExperianStatus != http.StatusOK && ScoringStatus != http.StatusOK {
					response.Result = &med2
					scoringStatusInternal = http.StatusBadGateway

				}
				scoringStatusInternal = http.StatusOK
				scorePro = res
				ExperianResp = resp
				ExperianResult = ExperianResp.Payload.Body.CreditStatus.Result
				ExperianScore = ExperianResp.Payload.Body.CreditStatus.Score
				TransID = ExperianResp.Payload.ResHeader.TransID
			} else {
				//HIT INTERNAL
				res, status, err := httpClient.ScoringPickle(req)

				scoringStatusInternal = http.StatusOK
				response.Result = &res.Result

				if status != http.StatusOK {
					response.Result = &med2
					scoringStatusInternal = http.StatusBadGateway
				}
				// EXPERIAN ISSUE
				ExperianStatusIssue = constant.ExperianIssue
				response.Status = constant.ScoringInternal
				response.Score = res.Score
				scorePro = res
				if err != nil {
					response.Result = &med2
					response.Status = constant.ScoringInternal
					response.Score = nil
					scorePro = res
				}
				ExperianResult = nil
				ExperianScore = nil
				TransID = nil
			}
			ExperianResp = resp
			ExperianResult = ExperianResp.Payload.Body.CreditStatus.Result
			ExperianScore = ExperianResp.Payload.Body.CreditStatus.Score
			TransID = ExperianResp.Payload.ResHeader.TransID

		} else if isScoringInternalActive != true && isIndosatPrefix {
			//INTERNAL OFF & EXPERIAN ON
			//IF DUMMY
			if experianDummy != nil {
				_ = json.Unmarshal([]byte(experianDummy.ExperianResponse), &resp)

				ExperianStatus = http.StatusOK
			} else {
				//CHECKING SCORE TO EXPERIAN
				creditScore, status, err := httpClient.CheckCreditScore(experian.PhoneNumber, token)
				if err != nil {
					log.Info("CheckExperian : ", err)
				}
				resp = creditScore
				ExperianStatus = status
			}

			if resp != nil {
				response.Score = resp.Payload.Body.CreditStatus.Score
				response.Result = resp.Payload.Body.CreditStatus.Result
				response.Status = constant.Experian

				result := resp.Payload.Body.CreditStatus.Result
				scoringStatusInternal = http.StatusOK

				if ExperianStatus != http.StatusOK {
					scoringStatusInternal = http.StatusOK
					response.Result = result

					//HIT INTERNAL
					res, status, _ := httpClient.ScoringPickle(req)

					response.Result = &res.Result
					response.Status = constant.ScoringInternal
					response.Score = res.Score

					//EXPERIAN ISSUE
					ExperianStatusIssue = constant.ExperianIssue

					if status != http.StatusOK {
						response.Result = &med2
						scoringStatusInternal = http.StatusBadGateway
					}
					scorePro = res

				} else if *result == medium && ExperianStatus == http.StatusOK {
					//HIT INTERNAL
					res, _, _ := httpClient.ScoringPickle(req)

					response.Score = res.Score
					response.Result = &high
					response.Status = constant.ScoringInternal

					if res.Result != constant.High {
						response.Score = resp.Payload.Body.CreditStatus.Score
						response.Result = &low
						response.Status = constant.Experian
					}
					scorePro = res
				} else if *result == low {
					response.Status = constant.Reject
				}
				ExperianResp = resp
				ExperianResult = ExperianResp.Payload.Body.CreditStatus.Result
				ExperianScore = ExperianResp.Payload.Body.CreditStatus.Score
				TransID = ExperianResp.Payload.ResHeader.TransID
			}  else {
				//HIT INTERNAL
				res, status, err := httpClient.ScoringPickle(req)

				scoringStatusInternal = http.StatusOK

				if status != http.StatusOK {
					scoringStatusInternal = http.StatusBadGateway
				}
				response.Result = &res.Result
				response.Status = constant.ScoringInternal
				response.Score = res.Score
				scorePro = res
				if err != nil {
					response.Result = &med2
					response.Status = constant.ScoringInternal
					response.Score = nil
					scorePro = res
				}
				ExperianResult = nil
				ExperianScore = nil
				TransID = nil
			}
		}else if isIndosatPrefix != true {
			//INTERNAL ONLY
			//HIT INTERNAL
			res, status, err := httpClient.ScoringPickle(req)

			scoringStatusInternal = http.StatusOK


			response.Result = &res.Result
			response.Status = constant.ScoringInternal
			response.Score = res.Score
			scorePro = res
			if status != http.StatusOK {
				response.Result = &med2
				scoringStatusInternal = http.StatusBadGateway
			}
			if err != nil {
				response.Result = &med2
				response.Status = constant.ScoringInternal
				response.Score = nil
				scorePro = res
			}
			ExperianResult = nil
			ExperianScore = nil
			TransID = nil
		}

		exData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(ExperianResp)

	} else if Tipe == constant.TypeSCS2 {
		//Hit Pickle
		res, status, err := httpClient.ScoringPickle(req)

		scoringStatusInternal = http.StatusOK

		if status != http.StatusOK {
			scoringStatusInternal = http.StatusBadGateway
		}
		response.Result = &res.Result
		response.Status = constant.ScoringInternal
		response.Score = res.Score
		scorePro = res
		if err != nil {
			response.Result = &med2
			response.Status = constant.ScoringInternal
			response.Score = nil
			scorePro = res
		}

		//Get DB Experian Result if Indosat Number
		if isIndosatPrefix {
			DBExperian, err := c.service.FindTelcoScoreByPhoneNumber(param)
			if err != nil {
				log.Error(err)
			}
			exData = []byte(DBExperian.Experian)
			TransID = DBExperian.TransID
			ExperianScore = DBExperian.ExperianScore
			ExperianResult = DBExperian.ExperianResult
			response.Result = DBExperian.Result
			response.Status = DBExperian.Status
			response.Score = DBExperian.Score

			if DBExperian.Result == nil || DBExperian == nil {
				scoringStatusInternal = http.StatusBadGateway
			}
		}

	} else {
		return c.BadRequestWithSpecificFields(ctx, "Type must be SCS1 or SCS2")
	}

	spData, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(scorePro)

	//SAVE TO DB
	TelcoScore := models.Experian{
		ProspectID:     response.ProspectID,
		Result:         response.Result,
		Status:         response.Status,
		PhoneNumber:    param,
		Score:          response.Score,
		Experian:       string(exData),
		ScorePro:       string(spData),
		ExperianScore:  ExperianScore,
		ExperianResult: ExperianResult,
		InternalScore:  scorePro.Score,
		InternalResult: scorePro.Result,
		TransID:        TransID,
		Type:           Tipe,
	}
	id, err := c.service.SaveTelcoScore(TelcoScore)

	if err != nil {
		return c.InternalServerError(ctx, err)
	}

	response.ExperianID = id.ExperianID
	response.Status = response.Status + ExperianStatusIssue
	if scoringStatusInternal != http.StatusOK {
		return c.BadGateWay(ctx, response)
	}
	return c.Ok(ctx, response)

}

// Integrator godoc
// @Description Los Scorepro LIMIT
// @Tags SCORE PRO
// @Produce json
// @Param phoneNumber path string true "Mobile Phone"
// @Param body body dto.PickleLimitDto true "Body payload"
// @Param scs query string false "scoring type by scs"
// @Success 200 {object} ApiResponse{data=models.CreditScoreResp}
// @Failure 400 {object} ApiResponse{}
// @Failure 500 {object} ApiResponse{}
// @Router /score/credit/{phoneNumber}/limit [post]
func (c *TelcoScoreController) CreditScoreLimit(ctx echo.Context) error {
	var experian dto.ExperianDto
	var req dto.PickleLimitDto
	var scoringStatusInternal int
	var ExperianScore *float64
	var ExperianResult *string
	var TransID *string
	var isIndosatPrefix bool
	var isScoringInternalActive bool
	var ExperianStatus int
	var token string
	var scorePro models.PickleLimitResponse
	var response models.CreditScoreResp
	var ExperianResp *models.ExperianScoreResp
	var exData []byte
	var ExperianStatusIssue string

	high := string(constant.High)
	medium := string(constant.Medium)
	low := string(constant.Low)
	med2 := string(constant.Med2nd)


	param := ctx.Param("phoneNumber")
	tipe := ctx.QueryParam("type")
	Tipe := strings.ToUpper(tipe)
	if param[:1] == "0" {
		res := strings.Replace(param, "0", "62", 1)
		param = res
	} else if param[:1] != "0" && param[:2] != "62" {
		return c.BadRequestWithSpecificFieldResponses(ctx, "error validation", []string{"phone_number", "start=08,628"})
	}

	experian.PhoneNumber = param

	if err := ctx.Bind(&req); err != nil {
		return c.BadRequest(ctx, err)
	}

	//Validate Phone Number
	var validationErrors []models.ErrorValidation
	if err := ctx.Validate(&experian); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationApiErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}

	//Validate Request
	if err := ctx.Validate(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationApiErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}

	if len(req.ZIP3V) > 3 {
		req.ZIP3V = req.ZIP3V[:3]
	}

	_, isIndosatPrefix = c.service.IndosatValidate(experian.PhoneNumber)

	httpClient := httpclient.NewExperianHTTPClient()
	response = models.CreditScoreResp{
		ProspectID:  req.OrderID,
		PhoneNumber: param,
	}
	//Check IF Dummy is true
	var experianDummy *models.Dummy
	dummy := os.Getenv("EXPERIAN_DUMMY")
	if dummy == "true" {
		res, err := c.service.GetExperianDummy(experian.PhoneNumber)
		if res == nil {
			return c.NotFound(ctx, err)
		}
		experianDummy = res
	} else {
		//SET BIG CACHE Token Indosat
		Itoken, err := cache.Get("indosat-token")
		token = string(Itoken)
		if Itoken == nil || err != nil {
			res, err := httpClient.GetToken()
			if err != nil {
				log.Info("GetToken Error: ", err)

				// By pass Experian, next process to API Pickle
				isIndosatPrefix = false
			}
			cache.Set("indosat-token", []byte(res.AccessToken))
			Itoken, _ := cache.Get("indosat-token")
			token = string(Itoken)
		}
	}

	//Type SRS
	if Tipe == constant.TypeSCS1 || Tipe == "" {
		Tipe = constant.TypeSCS1

		isScoringInternalActive = true
		//SET CACHE USE_SCORING
		useScoring, err := cache.Get("use-scoring")
		if useScoring == nil || err != nil {
			//Check configurable hit internal scoring or Not
			data, val := c.service.CheckConfig()
			isActive := strconv.Itoa(data.IsActive)
			isScoringInternalActive = val
			cache.Set("use-scoring", []byte(isActive))
			useScoring, err = cache.Get("use-scoring")
		}

		if string(useScoring) != "1" {
			isScoringInternalActive = false
		}

		var resp *models.ExperianScoreResp
		if isScoringInternalActive && isIndosatPrefix {
			//INTERNAL ON && EXPERIAN ON
			//IF DUMMY
			if experianDummy != nil {
				_ = json.Unmarshal([]byte(experianDummy.ExperianResponse), &resp)
				ExperianStatus = http.StatusOK
			} else {
				//CHECKING SCORE TO EXPERIAN

				creditScore, status, err := httpClient.CheckCreditScore(experian.PhoneNumber, token)
				if err != nil {
					log.Info("CheckExperian : ", err)
				}
				resp = creditScore
				ExperianStatus = status

			}
			//CHECK INTERNAL SCORING
			res, ScoringStatus, err := httpClient.ScoringPickleLimit(req)
			if err != nil {
				log.Error("Scoring Internal : ", err)
			}

			if resp != nil && ExperianStatus == http.StatusOK {
				response.Score = resp.Payload.Body.CreditStatus.Score
				response.Result = resp.Payload.Body.CreditStatus.Result
				response.Status = constant.Experian

				if *resp.Payload.Body.CreditStatus.Result == medium && res.Result == constant.High {
					response.Score = res.Score
					response.Result = &res.Result
					response.Status = constant.ScoringInternal
				} else if *resp.Payload.Body.CreditStatus.Result == medium && res.Result != constant.High {
					response.Score = resp.Payload.Body.CreditStatus.Score
					response.Result = &low
					response.Status = constant.Experian
				}

				if ExperianStatus != http.StatusOK && ScoringStatus == http.StatusOK {
					response.Score = res.Score
					response.Result = &res.Result
					response.Status = constant.ScoringInternal
					//EXPERIAN ISSUE
					ExperianStatusIssue = constant.ExperianIssue

				} else if ExperianStatus != http.StatusOK && ScoringStatus != http.StatusOK {
					response.Result = &medium
					scoringStatusInternal = http.StatusBadGateway

				}
				scoringStatusInternal = http.StatusOK
				scorePro = res
				ExperianResp = resp
				ExperianResult = ExperianResp.Payload.Body.CreditStatus.Result
				ExperianScore = ExperianResp.Payload.Body.CreditStatus.Score
				TransID = ExperianResp.Payload.ResHeader.TransID
			} else {
				//HIT INTERNAL
				res, status, err := httpClient.ScoringPickleLimit(req)

				scoringStatusInternal = http.StatusOK
				response.Result = &res.Result

				if status != http.StatusOK {
					response.Result = &medium
					scoringStatusInternal = http.StatusBadGateway
				}
				// EXPERIAN ISSUE
				ExperianStatusIssue = constant.ExperianIssue
				response.Status = constant.ScoringInternal
				response.Score = res.Score
				scorePro = res
				if err != nil {
					response.Result = &medium
					response.Status = constant.ScoringInternal
					response.Score = nil
					scorePro = res
				}
				ExperianResult = nil
				ExperianScore = nil
				TransID = nil
			}
			ExperianResp = resp
			ExperianResult = ExperianResp.Payload.Body.CreditStatus.Result
			ExperianScore = ExperianResp.Payload.Body.CreditStatus.Score
			TransID = ExperianResp.Payload.ResHeader.TransID

		} else if isScoringInternalActive != true && isIndosatPrefix {
			//INTERNAL OFF & EXPERIAN ON
			//IF DUMMY
			if experianDummy != nil {
				_ = json.Unmarshal([]byte(experianDummy.ExperianResponse), &resp)

				ExperianStatus = http.StatusOK
			} else {
				//CHECKING SCORE TO EXPERIAN
				creditScore, status, err := httpClient.CheckCreditScore(experian.PhoneNumber, token)
				if err != nil {
					log.Info("CheckExperian : ", err)
				}
				resp = creditScore
				ExperianStatus = status
			}

			if resp != nil {
				response.Score = resp.Payload.Body.CreditStatus.Score
				response.Result = resp.Payload.Body.CreditStatus.Result
				response.Status = constant.Experian

				result := resp.Payload.Body.CreditStatus.Result
				scoringStatusInternal = http.StatusOK

				if ExperianStatus != http.StatusOK {
					scoringStatusInternal = http.StatusOK
					response.Result = result

					//HIT INTERNAL
					res, status, _ := httpClient.ScoringPickleLimit(req)

					response.Result = &res.Result
					response.Status = constant.ScoringInternal
					response.Score = res.Score

					//EXPERIAN ISSUE
					ExperianStatusIssue = constant.ExperianIssue

					if status != http.StatusOK {
						response.Result = &medium
						scoringStatusInternal = http.StatusBadGateway
					}
					scorePro = res

				} else if *result == medium && ExperianStatus == http.StatusOK {
					//HIT INTERNAL
					res, _, _ := httpClient.ScoringPickleLimit(req)

					response.Score = res.Score
					response.Result = &high
					response.Status = constant.ScoringInternal

					if res.Result != constant.High {
						response.Score = resp.Payload.Body.CreditStatus.Score
						response.Result = &low
						response.Status = constant.Experian
					}
					scorePro = res
				} else if *result == low {
					response.Status = constant.Reject
				}
				ExperianResp = resp
				ExperianResult = ExperianResp.Payload.Body.CreditStatus.Result
				ExperianScore = ExperianResp.Payload.Body.CreditStatus.Score
				TransID = ExperianResp.Payload.ResHeader.TransID
			}  else {
				//HIT INTERNAL
				res, status, err := httpClient.ScoringPickleLimit(req)

				scoringStatusInternal = http.StatusOK

				if status != http.StatusOK {
					scoringStatusInternal = http.StatusBadGateway
				}
				response.Result = &res.Result
				response.Status = constant.ScoringInternal
				response.Score = res.Score
				scorePro = res
				if err != nil {
					response.Result = &medium
					response.Status = constant.ScoringInternal
					response.Score = nil
					scorePro = res
				}
				ExperianResult = nil
				ExperianScore = nil
				TransID = nil
			}
		}else if isIndosatPrefix != true {
			//INTERNAL ONLY
			//HIT INTERNAL
			res, status, err := httpClient.ScoringPickleLimit(req)

			scoringStatusInternal = http.StatusOK

			response.Result = &res.Result
			response.Status = constant.ScoringInternal
			response.Score = res.Score
			scorePro = res

			if status != http.StatusOK {
				response.Result = &medium
				scoringStatusInternal = http.StatusBadGateway
			}

			if err != nil {
				response.Result = &medium
				response.Status = constant.ScoringInternal
				response.Score = nil
				scorePro = res
			}
			ExperianResult = nil
			ExperianScore = nil
			TransID = nil
		}

		exData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(ExperianResp)

	} else if Tipe == constant.TypeSCS2 {
		//Hit Pickle
		res, status, err := httpClient.ScoringPickleLimit(req)

		scoringStatusInternal = http.StatusOK

		if status != http.StatusOK {
			scoringStatusInternal = http.StatusBadGateway
		}
		response.Result = &res.Result
		response.Status = constant.ScoringInternal
		response.Score = res.Score
		scorePro = res
		if err != nil {
			response.Result = &medium
			response.Status = constant.ScoringInternal
			response.Score = nil
			scorePro = res
		}

		//Get DB Experian Result if Indosat Number
		if isIndosatPrefix {
			DBExperian, err := c.service.FindTelcoScoreByPhoneNumber(param)
			if err != nil {
				log.Error(err)
			}
			exData = []byte(DBExperian.Experian)
			TransID = DBExperian.TransID
			ExperianScore = DBExperian.ExperianScore
			ExperianResult = DBExperian.ExperianResult
			response.Result = DBExperian.Result
			response.Status = DBExperian.Status
			response.Score = DBExperian.Score

			if DBExperian.Result == nil || DBExperian == nil {
				scoringStatusInternal = http.StatusBadGateway
			}
		}

	} else {
		return c.BadRequestWithSpecificFields(ctx, "Type must be SCS1 or SCS2")
	}

	spData, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(scorePro)

	//SAVE TO DB
	result := response.Result
	if scoringStatusInternal != http.StatusOK && ExperianStatus != http.StatusOK{
		result = &med2
	}
	TelcoScore := models.Experian{
		ProspectID:     response.ProspectID,
		Result:         result,
		Status:         response.Status,
		PhoneNumber:    param,
		Score:          response.Score,
		Experian:       string(exData),
		ScorePro:       string(spData),
		ExperianScore:  ExperianScore,
		ExperianResult: ExperianResult,
		InternalScore:  scorePro.Score,
		InternalResult: scorePro.Result,
		TransID:        TransID,
		Type:           Tipe,
	}
	id, err := c.service.SaveTelcoScore(TelcoScore)

	if err != nil {
		return c.InternalServerError(ctx, err)
	}

	response.ExperianID = id.ExperianID
	response.Status = response.Status + ExperianStatusIssue
	if scoringStatusInternal != http.StatusOK {
		return c.BadGateWay(ctx, response)
	}
	return c.Ok(ctx, response)

}

//GET ALL SCORE CREDIT
func (c *TelcoScoreController) Detail(ctx echo.Context) error {
	id := ctx.Param("id")
	data, err := c.service.FindTelcoScoreById(id)
	if err != nil {
		return c.NotFound(ctx, err)
	}
	return c.Ok(ctx, data)
}

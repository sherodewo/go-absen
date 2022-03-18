package scorepro

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kreditplus/scorepro/constant"
	"github.com/kreditplus/scorepro/controllers"
	"github.com/kreditplus/scorepro/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type WgScoreproController struct {
	controllers.Controller
	controllers.BaseBackendController
	service *ScoreproService
}

func NewWgScoreproController(service *ScoreproService) WgScoreproController {
	return WgScoreproController{
		BaseBackendController: controllers.BaseBackendController{
			Menu:        "Wg Dashboard",
			BreadCrumbs: []map[string]interface{}{},
		},
		service: service,
	}
}

// Scoring godoc
// @Summary Telcoscore WG Scoring
// @Description Get Scoring Pickle & Experian
// @Tags TELCOSCORE
// @Produce json
// @Param body body PickleModelingDto true "Body payload"
// @Success 200 {object} controllers.ApiResponse{data=WgScoreproResponse}
// @Failure 400 {object} controllers.ApiResponse{}
// @Failure 500 {object} controllers.ApiResponse{}
// @Router /scorepro/wg/idx [post]
func (c *WgScoreproController) Scoring(ctx echo.Context) error {
	var request PickleModelingDto
	var response ScoreproResponse
	var isIndosatPrefix bool
	var resp *models.ExperianScoreResp
	var scoreproScoreID string
	var scoreCombinationID string

	if err := ctx.Bind(&request); err != nil {
		return c.BadRequest(ctx, err)
	}

	var validationErrors []models.ErrorValidation
	//  Validate Request Pickle
	if err := ctx.Validate(&request); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationApiErrors(errs)
		}
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "error validation", "errors": validationErrors})
	}

	// CHECK Requestor
	requestor, err := c.service.CheckRequestor(request.RequestorID)
	if err != nil || &requestor.ID == nil {
		return c.InternalServerError(ctx, err)
	}

	response.PhoneNumber = request.PhoneNumber
	response.ProspectID = request.ProspectID
	param := request.PhoneNumber

	// Validate Phone Number
	if param[:1] == "0" {
		res := strings.Replace(param, "0", "62", 1)
		param = res
	} else if param[:1] != "0" && param[:2] != "62" {
		return c.BadRequestWithSpecificFieldResponses(ctx, "error validation", []string{"phone_number", "start=08,628"})
	}

	// Check score requestor

	// Modeling Pickle By ScoreGeneratorID
	modelPickle, err := c.service.ModelingPickle(request.ScoreGeneratorID)
	if err != nil {
		log.Error(err)
	}

	//Initiate HTTPCLIENT
	httpClient := NewHTTPClient()
	transactionID := uuid.New().String()

	// Hit Pickle
	dataPickle, status, err, validatePickle, scoreproScore := c.service.HitScoringPickleWg(request, *modelPickle, transactionID, ctx, request.TransactionType)
	if validatePickle != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "error validation", "errors": validatePickle})
	}
	if err == nil {
		scoreproScoreID = scoreproScore.ID
		response.ScoreResult = dataPickle.Data.Result
	}

	if *request.CbFound  {
		if status != http.StatusOK {
			response.Result = constant.PASS
			response.ScoreResult = constant.MEDIUM2ND
			response.Status = constant.ASS_MEDIUM2ND
			if *request.CbFound {
				response.Status = constant.ASSCB_MEDIUM2ND
			}
		} else {
			response.Status = constant.ASSCB_SCORE
			if dataPickle.Data.Result == constant.High || dataPickle.Data.Result == constant.Medium {
				response.Result = constant.PASS
			} else {
				response.Result = constant.Reject
			}

			response.Score = &dataPickle.Data.Score
		}

	} else if !*request.CbFound {
		var token string
		//Check IF Dummy is true
		var experianDummy *models.Dummy
		_, dummy := c.service.UseDummy()
		if dummy && strings.ToUpper(os.Getenv("APP_ENV")) != constant.EnvProduction{
			res, err := c.service.GetExperianDummy(param)
			if err != nil {
				return c.InternalServerError(ctx, err)
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

		// Check Indosat Number
		_, isIndosatPrefix = c.service.IndosatValidate(param)
		if isIndosatPrefix {
			response.Status = constant.ASS_TELCO

			if experianDummy != nil {
				_ = json.Unmarshal([]byte(experianDummy.ExperianResponse), &resp)
			} else {
				// Hit To Experian
				experian, statusExperian, _, err := c.service.HitExperian(param, token,  constant.ScoringExperianWgOnl)
				if err != nil {
					log.Error(err)
					//return c.InternalServerError(ctx, err)
				}
				if statusExperian == http.StatusOK {
					resp = experian
				}
			}
			response.Status = constant.ASS_TELCO

		}else {
			response.Status = constant.ASS_SCORE
			response.Result = constant.PASS
			if dataPickle.Data.Result == constant.Low {
				response.Result = constant.Reject
			}
			response.Score = &dataPickle.Data.Score
		}


		if status != http.StatusOK {
			dataPickle = nil
		}

		// Score Combination
		scoreCombination, score, err := c.service.MetricCombintaion(isIndosatPrefix, resp, dataPickle, constant.WgOnline)
		if err != nil {
			log.Error(err)
			return c.InternalServerError(ctx, err)
		}
		scoreCombinationID = scoreCombination.ID
		response.Result = scoreCombination.ScoreLos
		response.Status = scoreCombination.FinalScore
		response.ScoreResult = scoreCombination.ExpectedScoreLos
		response.Score = score

	}

	//Store Data Inquiry
	dataInquiry, err := c.service.StoreDataInquiry(request, modelPickle.ID, scoreproScoreID, requestor.ID, transactionID, &scoreCombinationID)
	if err != nil {
		log.Error(err)
		return c.InternalServerError(ctx, err)
	}

	// Save to DB ExperianInquiry
	if status != http.StatusOK {
		dataPickle = nil
	}
	re:=regexp.MustCompile("[0-9]+")
	_, err = c.service.SaveTelcoScore(request, response, resp, dataPickle, dataInquiry.ID)
	if err != nil {
		return c.BadGateWay(ctx, err)
	}
	wgResponse := WgScoreproResponse{
		ProspectID:  response.ProspectID,
		Score:       response.Score,
		Result:      response.Result,
		MaxDSR:      scoreproScore.MaxDSR,
		ScoreBand:  re.FindString(scoreproScore.ScoreBand),
		ScoreResult: response.ScoreResult,
		Status:      response.Status,
		PhoneNumber: response.PhoneNumber,
	}
	return c.Ok(ctx, wgResponse)
}
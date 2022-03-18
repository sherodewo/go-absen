package scorepro

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kreditplus/scorepro/constant"
	"github.com/kreditplus/scorepro/controllers"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ScoreproController struct {
	controllers.Controller
	controllers.BaseBackendController
	service *ScoreproService
}

func NewTelcoScoreController(service *ScoreproService) ScoreproController {
	return ScoreproController{
		BaseBackendController: controllers.BaseBackendController{
			Menu:        "Dashboard",
			BreadCrumbs: []map[string]interface{}{},
		},
		service: service,
	}
}

func (c *ScoreproController) Index(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "List Data",
		"link": "/scorepro/admin/dashboard",
	}
	return controllers.Render(ctx, "Dashboard Incoming Offline", "incoming_offline/index", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}

func (c *ScoreproController) List(ctx echo.Context) error {

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

		action = `<button class="btn btn-success btn-bold btn-upper" href="javascript:;" onclick="Detail('` + v.ID + `')" data-toggle="modal" data-target="#detail" data-id="` + v.ID + `" style="text-decoration: none;font-weight: 400; color: white;">
                      Detail </button>`
		time := v.CreatedAt
		createdAt = time.Format("2006-01-02T15:04:05+07:00")
		listOfData[k] = map[string]interface{}{
			"experian_id":  v.ID,
			"phone_number": v.PhoneNumber,
			"ProspectID":   v.ProspectID,
			"result":       v.Result,
			"status":       v.Status,
			"score":        v.Score,
			"type":         v.Type,
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

func (c *ScoreproController) Scoring(ctx echo.Context) error {
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
	dataPickle, status, err, validatePickle, idScore := c.service.HitScoringPickle(request, *modelPickle, transactionID, ctx)
	if validatePickle != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "error validation", "errors": validatePickle})
	}
	if err == nil {
		scoreproScoreID = *idScore
		response.ScoreResult = dataPickle.Data.Result
	}

	if request.Journey == constant.SCP1 {
		if status != http.StatusOK {
			response.Result = constant.PASS
			response.ScoreResult = constant.MEDIUM2ND
			response.Status = constant.ASS_MEDIUM2ND
			if *request.CbFound {
				response.Status = constant.ASSCB_MEDIUM2ND
			}
		} else {
			// Check CB Found
			if *request.CbFound {
				response.Status = constant.ASSCB_SCORE
				if dataPickle.Data.Result == constant.High || dataPickle.Data.Result == constant.Medium {
					response.Result = constant.PASS
				} else {
					response.Result = constant.Reject
				}
			} else {
				response.Status = constant.ASS_SCORE
				response.Result = constant.PASS
				if dataPickle.Data.Result == constant.Low {
					response.Result = constant.Reject
				}
			}
			response.Score = &dataPickle.Data.Score
		}

	} else if request.Journey == constant.SCP2 {
		if status == http.StatusOK {
			if dataPickle.Data.Result == constant.Low {
				response.ScoreResult = dataPickle.Data.Result
				response.Score = &dataPickle.Data.Score
				response.Status = constant.ASS_LOW
				response.Result = constant.Reject
			}
		}
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
			if experianDummy != nil {
				_ = json.Unmarshal([]byte(experianDummy.ExperianResponse), &resp)
			} else {
				// Hit To Experian
				experian, statusExperian, _, err := c.service.HitExperian(param, token, constant.ScoringExperianWgOff)
				if err != nil {
					log.Error(err)
					//return c.InternalServerError(ctx, err)
				}
				if statusExperian == http.StatusOK {
					resp = experian
				}
			}

		}
		// Score Combination
		scoreCombination, score, err := c.service.MetricCombintaion(isIndosatPrefix, resp, dataPickle, constant.WgOffline)
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
	_, err = c.service.SaveTelcoScore(request, response, resp, dataPickle, dataInquiry.ID)
	if err != nil {
		return c.BadGateWay(ctx, err)
	}
	return c.Ok(ctx, response)
}

//GET ALL SCORE CREDIT
func (c *ScoreproController) Detail(ctx echo.Context) error {
	id := ctx.Param("id")
	data, err := c.service.FindTelcoScoreById(id)
	if err != nil {
		return c.NotFound(ctx, err)
	}
	return c.Ok(ctx, data)
}

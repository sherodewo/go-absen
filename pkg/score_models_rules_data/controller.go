package score_models_rules_data

import (
	"github.com/kreditplus/scorepro/controllers"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type ScoreModelsRulesDataController struct {
	controllers.BaseBackendController
	service *ScoreModelsRulesDataService
}

func NewScoreModelsRulesDataController(service *ScoreModelsRulesDataService) ScoreModelsRulesDataController {
	return ScoreModelsRulesDataController{
		BaseBackendController: controllers.BaseBackendController{
			Menu:        "Suppliers",
			BreadCrumbs: []map[string]interface{}{},
		},
		service: service,
	}
}

func (c *ScoreModelsRulesDataController) Index(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "List Score Models Rules Data",
		"link": "/scorepro/admin/supplier",
	}
	return controllers.Render(ctx, "Score Models Rules Data List", "supplier/index", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}

func (c *ScoreModelsRulesDataController) List(ctx echo.Context) error {

	draw, err := strconv.Atoi(ctx.Request().URL.Query().Get("draw"))
	search := ctx.Request().URL.Query().Get("search[value]")
	start, err := strconv.Atoi(ctx.Request().URL.Query().Get("start"))
	length, err := strconv.Atoi(ctx.Request().URL.Query().Get("length"))
	order, err := strconv.Atoi(ctx.Request().URL.Query().Get("order[0][column]"))
	orderName := ctx.Request().URL.Query().Get("columns[" + strconv.Itoa(order) + "][name]")
	//orderAscDesc := ctx.Request().URL.Query().Get("order[0][dir]")

	recordTotal, recordFiltered, data, err := c.service.QueryDatatable(search, "desc", orderName, length, start)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	var action string

	//var role string
	listOfData := make([]map[string]interface{}, len(data))
	for k, v := range data {

		action = `<div class="btn-group open">
		<button class="btn btn-xs dropdown-toggle" type="button" data-toggle="dropdown" aria-expanded="true"> Actions</button>
		<ul class="dropdown-menu" role="menu">
		<li>
		<a href="JavaScript:void(0);" onclick="Edit('` + v.ID + `')" data-toggle="modal" data-target="#edit" data-placement="right" title="Set Active"><i class="fa fa-lock-open"></i>Edit</a>
		</li>
		<li>
		<a href="JavaScript:void(0);" onclick="Delete('` + v.ID + `')" style="text-decoration: none;font-weight: 400; color: #333;" data-toggle="tooltip" data-placement="right" title="Delete"><i class="fa fa-trash" style="color: #ff4d65"></i>Delete</a>
		</li>
		</ul>
		</div>`

		listOfData[k] = map[string]interface{}{
			"key":    v.Key,
			"value":   v.Value,
			"description":   v.Description,
			"score_generators":   v.ScoreGenerator,
			"action": action,
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

func (c *ScoreModelsRulesDataController) GetByID(ctx echo.Context) error {
	id := ctx.Param("id")
	data, err := c.service.FindScoreModelsRulesDataById(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, data)
}

func (c *ScoreModelsRulesDataController) Store(ctx echo.Context) error {
	var request ScoreModelsRulesDataReq
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(&request); err != nil {
		var validationErrors []models.ValidationError
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}

	result, err := c.service.SaveScoreModelsRulesData(request)
	if err != nil {
		return ctx.JSON(400, echo.Map{"message": "error save data ScoreModelsRulesData"})
	}

	session.SetFlashMessage(ctx, "save data success", "success", nil)
	return ctx.JSON(200, echo.Map{"message": "data has been saved", "data": result})
}

func (c *ScoreModelsRulesDataController) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var request ScoreModelsRulesDataReq
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(400, echo.Map{"message": "error binding data"})
	}

	if err := ctx.Validate(&request); err != nil {
		var validationErrors []models.ValidationError
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}

	result, err := c.service.UpdateScoreModelsRulesData(id, request)
	if err != nil {
		return ctx.JSON(400, echo.Map{"message": "error update data ScoreModelsRulesData"})
	}
	session.SetFlashMessage(ctx, "update data success", "success", nil)
	return ctx.JSON(200, echo.Map{"message": "data has been updated", "data": result})
}

func (c *ScoreModelsRulesDataController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	err := c.service.DeleteScoreModelsRulesData(id)
	if err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying delete data"})
	}
	return ctx.JSON(200, echo.Map{"message": "delete data has been deleted"})
}
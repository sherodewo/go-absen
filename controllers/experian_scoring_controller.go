package controllers

import (
	"github.com/jinzhu/gorm"
	"github.com/kreditplus/scorepro/dto"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/service"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type ExperianController struct {
	BaseBackendController
	service *service.ExperianService
}

func NewExperianController(service *service.ExperianService) ExperianController {
	return ExperianController{
		BaseBackendController: BaseBackendController{
			Menu:        "Experian Scoring",
			BreadCrumbs: []map[string]interface{}{},
		},
		service: service,
	}
}

func (c *ExperianController) Index(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "List Data",
		"link": "/scorepro/admin/experian-scoring",
	}
	return Render(ctx, "Experian Scoring", "experian_scoring/index", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}

func (c *ExperianController) List(ctx echo.Context) error {

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
	listOfData := make([]map[string]interface{}, len(data))
	for k, v := range data {

		action = `<div class="btn-group open">
					<button class="btn btn-xs dropdown-toggle" type="button" data-toggle="dropdown" aria-expanded="true"> Actions</button>
                      <ul class="dropdown-menu" role="menu">
                      	<li>
        	 				<a href="JavaScript:void(0);" onclick="Edit('` + v.ID + `')" data-toggle="modal" data-target="#edit-modal" data-placement="right" title="Set Active"><i class="fa fa-lock-open"></i>Edit</a>
      					</li>
      					<li>
         					<a href="JavaScript:void(0);" onclick="Delete('` + v.ID + `')" style="text-decoration: none;font-weight: 400; color: #333;" data-toggle="tooltip" data-placement="right" title="Delete"><i class="fa fa-trash" style="color: #ff4d65"></i>Delete</a>
      					</li>
                      </ul>
                      </div>`
		listOfData[k] = map[string]interface{}{
			"is_indosat":  v.IsIndosat,
			"experian":    v.Experian,
			"internal":    v.Internal,
			"score_los":   v.ScoreLos,
			"final_score": v.FinalScore,
			"notes":       v.Notes,
			"action":      action,
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

func (c *ExperianController) Store(ctx echo.Context) error {
	var req dto.ExperianScoringDto
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, echo.Map{"message": "error binding data"})
	}

	if err := ctx.Validate(&req); err != nil {
		var validationErrors []models.ValidationError
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}

	result, err := c.service.SaveExperian(req)
	if err != nil {
		return ctx.JSON(400, echo.Map{"message": "error save data user"})
	}

	session.SetFlashMessage(ctx, "store data success", "success", result)
	return ctx.Redirect(302, "/scorepro/admin/experian")
}

func (c *ExperianController) Detail(ctx echo.Context) error {
	id := ctx.Param("id")
	data, err := c.service.FindExperianById(id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/scorepro/admin/experian")
		}
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/experian")
	}

	return ctx.JSON(http.StatusOK, data)
}

func (c *ExperianController) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var req dto.ExperianScoringUpdateDto
	if err := ctx.Bind(&req); err != nil {
		session.SetFlashMessage(ctx, "error binding data", "error", nil)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	result, err := c.service.UpdateExperian(id, req)
	if err != nil {
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	session.SetFlashMessage(ctx, "update data success", "success", result)
	return ctx.JSON(http.StatusOK, result)
}
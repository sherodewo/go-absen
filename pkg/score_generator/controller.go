package score_generator

import (
	"github.com/kreditplus/scorepro/controllers"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
)

type ScoreController struct {
	controllers.BaseBackendController
	controllers.Controller
	service *ScoreService
}

func NewScoreController(service *ScoreService) ScoreController {
	return ScoreController{
		BaseBackendController: controllers.BaseBackendController{
			Menu:        "Dashboard",
			BreadCrumbs: []map[string]interface{}{},
		},
		service: service,
	}
}

func (c *ScoreController) Index(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "List Data",
		"link": "/scorepro/admin/score-generator",
	}
	return controllers.Render(ctx, "Dashboard Incoming Offline", "score_generator/index", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}

func (c *ScoreController) Create(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "Create",
		"link": "/scorepro/admin/score-generator/create",
	}
	return controllers.Render(ctx, "Create Score Generator", "score_generator/create", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}

func (c *ScoreController) List(ctx echo.Context) error {

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
		//time := v.CreatedAt
		//createdAt = time.Format("2006-01-02T15:04:05+07:00")
		listOfData[k] = map[string]interface{}{

			"name":                    v.Name,
			"endpoint":                v.Endpoint,
			"file_pickle":             v.FilePickle,
			"save_result_to":          v.SaveResultTo,
			"save_result_object_name": v.SaveResultObjectName,
			"score_models_rules":      v.ScoreModelsRules,
			"description":             v.Description,
			"action":                  action,
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

func (c *ScoreController) Store(ctx echo.Context) error {
	var request WizardDto

	if err := ctx.Bind(&request); err != nil {
		session.SetFlashMessage(ctx, "error binding data", "error", nil)
		return ctx.JSON(http.StatusInternalServerError,err)
	}

	//Store Score Generator
	scoreGenerator, err := c.service.StoreScoreGenerator(request)
	if err != nil {
		session.SetFlashMessage(ctx, "error store score generator", "error", nil)
		return ctx.JSON(http.StatusInternalServerError,err)
	}

	//Store Score Generator Metadata
	res, err := c.service.StoreScoreGeneratorMetadata(request, scoreGenerator.ID)
	if err != nil {
		session.SetFlashMessage(ctx, "error store score generator meta data", "error", nil)
		return ctx.JSON(http.StatusInternalServerError,err)
	}

	// Create Table
	err = c.service.CreateTable(request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,err)
	}

	session.SetFlashMessage(ctx, "Succses Add Data", "success", nil)
	return c.Ok(ctx, res)
}

func (c *ScoreController) Upload(ctx echo.Context) error {
	replace := ctx.Param("replace")
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Error("[Error] : ", err)
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()


	status,messages ,err := c.service.UploadFilePickle(file.Filename,src,replace)
	if err != nil || status != http.StatusOK {
		return ctx.JSON(http.StatusInternalServerError,messages)
	}

	return c.Ok(ctx, nil)
}

func (c *ScoreController) GetScoreModels(ctx echo.Context) error {
	data, err := c.service.GetScoreModelsRules()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return c.Ok(ctx, data)
}

func (c *ScoreController) GetAllScoreGenerator(ctx echo.Context) error {
	data, err := c.service.GetAllScoreGenerator()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return c.Ok(ctx, data)
}

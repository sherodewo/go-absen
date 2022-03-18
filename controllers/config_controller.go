package controllers

import (
	"github.com/jinzhu/gorm"
	"github.com/kreditplus/scorepro/dto"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/service"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ConfigController struct {
	BaseBackendController
	service *service.ConfigService
}

func NewConfigController(service *service.ConfigService) ConfigController {
	return ConfigController{
		BaseBackendController: BaseBackendController{
			Menu: "Config",
			BreadCrumbs: []map[string]interface{}{
				0: {
					"menu": "Menu",
					"link": "/scorepro/admin/config",
				},
			},
		},
		service: service,
	}
}

func (c *ConfigController) Index(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "List Configuration",
		"link": "/scorepro/admin/config/list",
	}
	return Render(ctx, "List Config", "config/index", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}

func (c *ConfigController) Datatable(ctx echo.Context) error {

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
	var active string
	var createdAt string

	listOfData := make([]map[string]interface{}, len(data))
	for k, v := range data {
		//action = `<a href="/scorepro/admin/config/submenu/` + v.ID + `/list" data-toggle="tooltip" data-placement="right" title="Add sub menu"><i class="la la-plus"></i></a>`
		//action += `<a href="/scorepro/admin/config/detail/` + v.ID + `" data-toggle="tooltip" data-placement="right" title="Detail"><i class="fa fa-user"></i> </a>`
		//action += `<a href="/scorepro/admin/config/edit/` + v.ID + `" data-toggle="tooltip" data-placement="right" title="Edit"><i class="fa fa-edit"></i> </a>`
		//action += `<a href="JavaScript:void(0);" onclick="SetActive('` + v.ID + `')" data-toggle="tooltip" data-placement="right" title="Set Active"><i class="fa fa-lock-open"></i></a>`
		//action += `<a href="JavaScript:void(0);" onclick="SetInactive('` + v.ID + `')" data-toggle="tooltip" data-placement="right" title="Set Inactive"><i class="fa fa-lock" style="color: #ff4d65"></i></a>`
		//action += `<a href="JavaScript:void(0);" onclick="Delete('` + v.ID + `')" data-toggle="tooltip" data-placement="right" title="Delete"><i class="fa fa-trash" style="color: #ff4d65"></i></a>`

		action =
			`<div class="btn-group open">
   <button class="btn btn-xs dropdown-toggle" type="button" data-toggle="dropdown" aria-expanded="true"> Actions</button>
   <ul class="dropdown-menu" role="menu">
      <li>
         <a href="JavaScript:void(0);" onclick="Edit('` + v.ID + `')" data-toggle="modal" data-target="#edit" data-placement="right" title="Set Active"><i class="fa fa-lock-open"></i>Edit</a>
      </li>
   </ul>
</div>`

		if v.GroupName != "indosat" && v.GroupName != "response_service_off"{
			if v.IsActive == 1 {
				active = `<span class="kt-switch kt-switch--outline kt-switch--icon kt-switch--warning">
							<label>
								<input id="` + v.ID + `" type="checkbox" checked="checked" value="true" name="is_active" href="JavaScript:void(0);" onclick="SetStatus('` + v.ID + `')">
								<span></span>
							</label>
						</span>`
			} else {
				active = `<span class="kt-switch kt-switch--outline kt-switch--icon kt-switch--warning">
							<label>
								<input id="` + v.ID + `" type="checkbox" value="false" name="is_active" href="JavaScript:void(0);" onclick="SetStatus('` + v.ID + `')">
								<span></span>
							</label>
						</span>`
			}
		}else {
			active =``
		}

		time := v.CreatedAt
		createdAt = time.Format("2006-01-02T15:04:05+07:00")

		listOfData[k] = map[string]interface{}{
			"group_name": v.GroupName,
			"key":        v.Key,
			"value":      v.Value,
			"active":     active,
			"created_at": createdAt,
			"action":     action,
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

func (c *ConfigController) Store(ctx echo.Context) error {
	var menuDto dto.AppConfDto
	if err := ctx.Bind(&menuDto); err != nil {
		session.SetFlashMessage(ctx, "error binding data", "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/config")
	}
	menuDto.IsActive = 1
	if err := ctx.Validate(&menuDto); err != nil {
		session.SetFlashMessage(ctx, "Validation Error", "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/config")
	}
	result, err := c.service.StoreMenu(menuDto)
	if err != nil {
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/config")
	}
	session.SetFlashMessage(ctx, "store data success", "success", result)
	return ctx.Redirect(302, "/scorepro/admin/config/list")
}

func (c *ConfigController) Detail(ctx echo.Context) error {
	id := ctx.Param("id")
	data, err := c.service.FindById(id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/scorepro/admin/config")
		}
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/config")
	}

	return ctx.JSON(http.StatusOK, data)
}

func (c *ConfigController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	err := c.service.DeleteConfig(id)
	if err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying delete data"})
	}
	return ctx.JSON(200, echo.Map{"message": "delete data has been deleted"})
}

func (c *ConfigController) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var req dto.AppConfDto
	if err := ctx.Bind(&req); err != nil {
		session.SetFlashMessage(ctx, "error binding data", "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/config/list")
	}

	if err := ctx.Validate(&req); err != nil {
		session.SetFlashMessage(ctx, "Validation Error", "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/config/edit/"+id)
	}

	result, err := c.service.UpdateConfig(id, req)
	if err != nil {
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/config/list")
	}
	session.SetFlashMessage(ctx, "update data success", "success", result)
	return ctx.Redirect(302, "/scorepro/admin/config/list")
}

func (c *ConfigController) SetActive(ctx echo.Context) error {
	var app models.AppConfig
	id := ctx.Param("id")

	if err := c.service.GetDbInstance().Model(&app).Where("id = ?", id).
		Update("is_active", 1).Error; err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying set active data"})
	}
	data, _ := c.service.FindById(id)

	switch data.GroupName {
	case "indosat":
		cache.Set("indosat-validation", []byte(data.Value))
	case "use_scoring":
		isActive := strconv.Itoa(1)
		cache.Set("use-scoring", []byte(isActive))
	}
	return ctx.JSON(200, echo.Map{"message": "success set active data"})
}

func (c *ConfigController) SetInactive(ctx echo.Context) error {
	var app models.AppConfig
	id := ctx.Param("id")

	if err := c.service.GetDbInstance().Model(&app).Where("id = ?", id).
		Update("is_active", 0).Error; err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying set inactive data"})
	}
	data, _ := c.service.FindById(id)

	switch data.GroupName {
	case "indosat":
		isActive := strconv.Itoa(0)
		cache.Set("indosat-active", []byte(isActive))
		cache.Set("indosat-validation", []byte(data.Value))
	case "use_scoring":
		isActive := strconv.Itoa(0)
		cache.Set("use-scoring", []byte(isActive))
	}
	return ctx.JSON(200, echo.Map{"message": "success set inactive data"})
}

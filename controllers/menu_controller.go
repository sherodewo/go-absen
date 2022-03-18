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

type MenuController struct {
	BaseBackendController
	service *service.MenuService
}

func NewMenuController(service *service.MenuService) MenuController {
	return MenuController{
		BaseBackendController: BaseBackendController{
			Menu: "Menus",
			BreadCrumbs: []map[string]interface{}{
				0: {
					"menu": "Menu",
					"link": "/scorepro/admin/menus",
				},
			},
		},
		service: service,
	}
}

func (c *MenuController) Index(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "List Menu",
		"link": "/scorepro/admin/menus/list",
	}
	return Render(ctx, "List Menu", "menu/index", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}

func (c *MenuController) Add(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "Add",
		"link": "/scorepro/admin/menus/add",
	}
	return Render(ctx, "Add Menu", "menu/add", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}

func (c *MenuController) Datatable(ctx echo.Context) error {

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
		//action = `<a href="/scorepro/admin/menus/submenu/` + v.ID + `/list" data-toggle="tooltip" data-placement="right" title="Add sub menu"><i class="la la-plus"></i></a>`
		//action += `<a href="/scorepro/admin/menus/detail/` + v.ID + `" data-toggle="tooltip" data-placement="right" title="Detail"><i class="fa fa-user"></i> </a>`
		//action += `<a href="/scorepro/admin/menus/edit/` + v.ID + `" data-toggle="tooltip" data-placement="right" title="Edit"><i class="fa fa-edit"></i> </a>`
		//action += `<a href="JavaScript:void(0);" onclick="SetActive('` + v.ID + `')" data-toggle="tooltip" data-placement="right" title="Set Active"><i class="fa fa-lock-open"></i></a>`
		//action += `<a href="JavaScript:void(0);" onclick="SetInactive('` + v.ID + `')" data-toggle="tooltip" data-placement="right" title="Set Inactive"><i class="fa fa-lock" style="color: #ff4d65"></i></a>`
		//action += `<a href="JavaScript:void(0);" onclick="Delete('` + v.ID + `')" data-toggle="tooltip" data-placement="right" title="Delete"><i class="fa fa-trash" style="color: #ff4d65"></i></a>`

		action =
			`<div class="btn-group open">
   <button class="btn btn-xs dropdown-toggle" type="button" data-toggle="dropdown" aria-expanded="true"> Actions</button>
   <ul class="dropdown-menu" role="menu">
     
   </ul>
</div>`
		listOfData[k] = map[string]interface{}{
			"name":       v.Name,
			"route":      v.Route,
			"icon_class": v.IconClass,
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

func (c *MenuController) Store(ctx echo.Context) error {
	var menuDto dto.MenuDto
	if err := ctx.Bind(&menuDto); err != nil {
		session.SetFlashMessage(ctx, "error binding data", "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/menus/list")
	}
	if err := ctx.Validate(&menuDto); err != nil {
		session.SetFlashMessage(ctx, "Validation Error", "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/menus/add")
	}
	result, err := c.service.StoreMenu(menuDto)
	if err != nil {
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/menus/list")
	}
	session.SetFlashMessage(ctx, "store data success", "success", result)
	return ctx.Redirect(302, "/scorepro/admin/menus/list")
}

func (c *MenuController) Detail(ctx echo.Context) error {
	id := ctx.Param("id")
	data, err := c.service.FindById(id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/scorepro/admin/menus/list")
		}
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/menus/list")
	}
	breadCrumbs := map[string]interface{}{
		"menu": "Detail",
		"link": "/scorepro/admin/menus/detail",
	}
	return Render(ctx, "Detail Menu", "menu/detail", c.Menu, session.FlashMessage{},
		append(c.BreadCrumbs, breadCrumbs), data)
}

func (c *MenuController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	err := c.service.DeleteMenu(id)
	if err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying delete data"})
	}
	return ctx.JSON(200, echo.Map{"message": "delete data has been deleted"})
}

func (c *MenuController) Edit(ctx echo.Context) error {
	id := ctx.Param("id")
	data, err := c.service.FindById(id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/scorepro/admin/menus/list")
		}
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/menus/list")
	}
	breadCrumbs := map[string]interface{}{
		"menu": "Edit",
		"link": "/scorepro/admin/menus/edit",
	}
	return Render(ctx, "Edit Role", "menu/edit", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), data)
}

func (c *MenuController) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var menuUpdateDto dto.MenuUpdateDto
	if err := ctx.Bind(&menuUpdateDto); err != nil {
		session.SetFlashMessage(ctx, "error binding data", "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/menus/list")
	}

	if err := ctx.Validate(&menuUpdateDto); err != nil {
		session.SetFlashMessage(ctx, "Validation Error", "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/menus/edit/"+id)
	}

	result, err := c.service.UpdateMenu(id, menuUpdateDto)
	if err != nil {
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/menus/list")
	}
	session.SetFlashMessage(ctx, "update data success", "success", result)
	return ctx.Redirect(302, "/scorepro/admin/menus/list")
}

func (c *MenuController) SetActive(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := c.service.GetDbInstance().Model(&models.Menu{}).Where("user_type =?", id).
		Update("is_active", true).Error; err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying set active data"})
	}
	return ctx.JSON(200, echo.Map{"message": "success set active data"})
}

func (c *MenuController) SetInactive(ctx echo.Context) error {
	id := ctx.Param("id")

	if err := c.service.GetDbInstance().Model(&models.Menu{}).Where("user_type =?", id).
		Update("is_active", false).Error; err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying set active data"})
	}
	return ctx.JSON(200, echo.Map{"message": "success set active data"})
}

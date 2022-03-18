package controllers

import (
	"errors"
	"github.com/kreditplus/scorepro/dto"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/service"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserController struct {
	BaseBackendController
	service *service.UserService
}

func NewUserController(service *service.UserService) UserController {
	return UserController{
		BaseBackendController: BaseBackendController{
			Menu:        "Users",
			BreadCrumbs: []map[string]interface{}{},
		},
		service: service,
	}
}

func (c *UserController) Index(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "List User",
		"link": "/scorepro/admin/users",
	}
	return Render(ctx, "User List", "user/index", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}

func (c *UserController) Add(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "Add",
		"link": "/scorepro/admin/users/add",
	}
	var Role []models.UserRole
	_ = c.service.GetDbInstance().Find(&Role)
	return Render(ctx, "Add User", "user/add", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), Role)
}

func (c *UserController) List(ctx echo.Context) error {

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
	var status string

	//var role string
	listOfData := make([]map[string]interface{}, len(data))
	for k, v := range data {

		action = `<div class="btn-group open">
					<button class="btn btn-xs dropdown-toggle" type="button" data-toggle="dropdown" aria-expanded="true"> Actions</button>
                      <ul class="dropdown-menu" role="menu">
                      <li>
                      <a href="/scorepro/admin/register/edit/` + v.UserID + `" style="text-decoration: none;font-weight: 400; color: #333;">
                      <i class="fa fa-edit"></i>Edit </a>
                      </li>
                      <li>
                      <a href="/scorepro/admin/register/detail/` + v.UserID + `"style="text-decoration: none;font-weight: 400; color: #333;">
                      <i class="fa fa-user"></i>Detail </a>
                      </li>
                      <li>
                      <a href="javascript:;" onclick="Delete('` + v.UserID + `')" data-toggle="tooltip" data-placement="right" title="Delete" style="text-decoration: none;font-weight: 400; color: #333;">
                      <i class="fa fa-trash" style="color: #ff4d65;"></i> Delete </a>
                      </li>
                      </ul>
                      </div>`

		if v.IsActive == 1 {
			status = `Aktif`
		} else {
			status = `Tidak Aktif`
		}
		var Role models.UserRole
		_ = c.service.GetDbInstance().First(&Role, models.UserRole{ID: v.UserRoleID})
		listOfData[k] = map[string]interface{}{
			"id":     v.UserID,
			"nik":    v.Nik,
			"name":   v.Name,
			"email":  v.Email,
			"role":   Role.Name,
			"status": status,
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

func (c *UserController) Store(ctx echo.Context) error {
	var userDto dto.UserDto
	if err := ctx.Bind(&userDto); err != nil {
		return ctx.JSON(400, echo.Map{"message": "error binding data"})
	}
	userDto.IsActive = 1
	if err := ctx.Validate(&userDto); err != nil {
		var validationErrors []models.ValidationError
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}

	result, err := c.service.SaveUser(userDto)
	if err != nil {
		return ctx.JSON(400, echo.Map{"message": "error save data user"})
	}

	session.SetFlashMessage(ctx, "save data success", "success", nil)
	return ctx.JSON(200, echo.Map{"message": "data has been saved", "data": result})
}

func (c *UserController) Edit(ctx echo.Context) error {
	id := ctx.Param("id")
	data, err := c.service.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/scorepro/admin/users")
		}
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/users")
	}
	breadCrumbs := map[string]interface{}{
		"menu": "Edit",
		"link": "/scorepro/admin/users/edit",
	}

	dataUser := models.UserViewDetail{
		UserID:     data.UserID,
		Nik:        data.Nik,
		Name:       data.Name,
		Email:      data.Email,
		Password:   data.Password,
		IsActive:   `unchecked`,
		UserRoleID: data.UserRoleID,
		TypeUser:   data.TypeUser,
	}
	if data.IsActive == 1 {
		dataUser.IsActive = `checked`
	}
	return Render(ctx, "Edit User", "user/edit", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), dataUser)
}

func (c *UserController) View(ctx echo.Context) error {
	id := ctx.Param("id")
	var data models.User
	err := c.service.GetDbInstance().First(&data, "user_id =?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			session.SetFlashMessage(ctx, err.Error(), "error", nil)
			return ctx.Redirect(302, "/scorepro/admin/users")
		}
		session.SetFlashMessage(ctx, err.Error(), "error", nil)
		return ctx.Redirect(302, "/scorepro/admin/users")
	}
	breadCrumbs := map[string]interface{}{
		"menu": "Detail User",
		"link": "/scorepro/admin/users/detail",
	}
	return Render(ctx, "Detail User ", "user/view", c.Menu, session.FlashMessage{},
		append(c.BreadCrumbs, breadCrumbs), echo.Map{"User": data})
}

func (c *UserController) Update(ctx echo.Context) error {
	id := ctx.Param("id")
	var userDto dto.UserUpdateDto
	if err := ctx.Bind(&userDto); err != nil {
		return ctx.JSON(400, echo.Map{"message": "error binding data"})
	}

	if err := ctx.Validate(&userDto); err != nil {
		var validationErrors []models.ValidationError
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors = models.WrapValidationErrors(errs)
		}
		return ctx.JSON(400, echo.Map{"message": "error validation", "errors": validationErrors})
	}
	//remove file before upload
	//f, err := c.service.FindUserById(id)
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return ctx.JSON(400, echo.Map{"message": "data user not found"})
	//	}
	//	return ctx.JSON(400, echo.Map{"message": "error get data user"})
	//}
	//if f.ImageUrl != "" {
	//	if err := utils.RemoveFile(f.ImageUrl); err != nil {
	//		return ctx.JSON(400, echo.Map{"message": "error remove image user"})
	//	}
	//}
	//
	//var fileName string
	//if userDto.Image != "" {
	//	f, err := utils.UploadImageV2(userDto.Image)
	//	if err != nil {
	//		return ctx.JSON(400, echo.Map{"message": "error upload image"})
	//	}
	//	fileName = f
	//}

	result, err := c.service.UpdateUser(id, userDto)
	if err != nil {
		return ctx.JSON(400, echo.Map{"message": "error update data user"})
	}
	session.SetFlashMessage(ctx, "update data success", "success", nil)
	return ctx.JSON(200, echo.Map{"message": "data has been updated", "data": result})
}

func (c *UserController) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	err := c.service.DeleteUser(id)
	if err != nil {
		return ctx.JSON(500, echo.Map{"message": "error when trying delete data"})
	}
	return ctx.JSON(200, echo.Map{"message": "delete data has been deleted"})
}

func (c *UserController) Profile(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "profile",
		"link": "/scorepro/admin/users/profile",
	}
	return Render(ctx, "Profile User", "user/profile", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), nil)
}

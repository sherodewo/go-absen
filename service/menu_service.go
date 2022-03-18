package service

import (
	"github.com/kreditplus/scorepro/dto"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/repository"
	"gorm.io/gorm"
)

type MenuService struct {
	MenuRepository repository.MenuRepository
}

func NewMenuService(repository repository.MenuRepository) *MenuService {
	return &MenuService{
		MenuRepository: repository,
	}
}

func (s *MenuService) QueryDatatable(searchValue string, orderType string, orderBy string, limit int, offset int) (
	recordTotal int64, recordFiltered int64, data []models.Menu, err error) {

	recordTotal, err = s.MenuRepository.CountWhere("and", map[string]interface{}{
		//"1 =?":         1,
	})

	if searchValue != "" {
		var dataCount []models.Menu
		qs := s.MenuRepository.DbInstance().Preload("Role").
			Where("parent_id =?", "").
			Where("(menu_title LIKE ? OR slug LIKE ?	OR url LIKE ?)", "%"+searchValue+"%", "%"+searchValue+"%", "%"+searchValue+"%")

		recordFiltered = qs.Find(&dataCount).RowsAffected

		err = qs.Order(orderBy + " " + orderType).
			Limit(limit).
			Offset(offset).Find(&data).Error

		return recordTotal, recordFiltered, data, err
	}

	recordFiltered, err = s.MenuRepository.CountWhere("and", map[string]interface{}{
		//"parent_id =?": "",
		//"1 =?":         1,
	})
	data, err = s.MenuRepository.FindAllWhere("and", orderType, orderBy, limit, offset, map[string]interface{}{
		//"parent_id =?": "",
		//"1 =?":         1,
	})

	return recordTotal, recordFiltered, data, err
}

func (s *MenuService) QueryDatatableSubMenu(id string, searchValue string, orderType string, orderBy string, limit int, offset int) (
	recordTotal int64, recordFiltered int64, data []models.Menu, err error) {

	recordTotal, err = s.MenuRepository.CountWhere("and", map[string]interface{}{
		"parent_id =?": id,
		"1 =?":         1,
	})

	if searchValue != "" {
		var dataCount []models.Menu
		qs := s.MenuRepository.DbInstance().Preload("Role").
			Where("parent_id =?", id).
			Where("(menu_title LIKE ? OR slug LIKE ?	OR url LIKE ?)", "%"+searchValue+"%", "%"+searchValue+"%", "%"+searchValue+"%")

		recordFiltered = qs.Find(&dataCount).RowsAffected

		err = qs.Order(orderBy + " " + orderType).
			Limit(limit).
			Offset(offset).Find(&data).Error

		return recordTotal, recordFiltered, data, err
	}

	recordFiltered, err = s.MenuRepository.CountWhere("and", map[string]interface{}{
		"parent_id =?": id,
		"1 =?":         1,
	})
	data, err = s.MenuRepository.FindAllWhere("and", orderType, orderBy, limit, offset, map[string]interface{}{
		"parent_id =?": id,
		"1 =?":         1,
	})

	return recordTotal, recordFiltered, data, err
}

func (s *MenuService) DeleteMenu(id string) error {
	//entity := models.Menu{
	//	ID: id,
	//}
	//err := s.MenuRepository.Delete(entity)
	//return err
	return nil
}

func (s *MenuService) FindById(id string) (*models.Menu, error) {
	data, err := s.MenuRepository.FindById(id)
	return &data, err
}

func (s *MenuService) FindAllMenu() (*[]models.Menu, error) {
	data, err := s.MenuRepository.FindAll()
	return &data, err
}

func (s *MenuService) FindMenuSelect2(search string) ([]map[string]interface{}, error) {
	if search != "" {
		data, err := s.MenuRepository.FindAllWhere("and", "asc", "created_at", 999, 0, map[string]interface{}{
			"name LIKE ?": "%" + search + "%",
		})
		result := make([]map[string]interface{}, len(data))
		for k, v := range data {
			result[k] = map[string]interface{}{
				"name":  v.Name,
			}
		}
		return result, err
	}

	data, err := s.MenuRepository.FindAll()
	result := make([]map[string]interface{}, len(data))
	for k, v := range data {
		result[k] = map[string]interface{}{
			"name":  v.Name,
		}
	}
	return result, err
}

func (s *MenuService) StoreMenu(dto dto.MenuDto) (data models.Menu, err error) {
	//if err := s.MenuRepository.DbInstance().Transaction(func(tx *gorm.DB) error {
	//	data = models.Menu{
	//		ID:          uuid.New().String(),
	//		ParentID:    dto.ParentID,
	//		MenuTitle:   dto.MenuTitle,
	//		Slug:        dto.Slug,
	//		Url:         dto.Url,
	//		Icon:        dto.Icon,
	//		MenuOrder:   dto.MenuOrder,
	//		Description: dto.Description,
	//		IsActive:    dto.IsActive,
	//	}
	//	if err = tx.Create(&data).Error; err != nil {
	//		return err
	//	}
	//	listMenuRole := make([]models.MenuRole, len(dto.Role))
	//	for k, v := range dto.Role {
	//		listMenuRole[k] = models.MenuRole{
	//			RoleID: v,
	//			MenuID: data.ID,
	//		}
	//	}
	//	if err = tx.Model(&models.MenuRole{}).Create(listMenuRole).Error; err != nil {
	//		return err
	//	}
	//	return nil
	//}); err != nil {
	//	return models.Menu{}, err
	//}
	return data, err
}

func (s *MenuService) UpdateMenu(id string, dto dto.MenuUpdateDto) (data models.Menu, err error) {

	//if err := s.MenuRepository.DbInstance().Transaction(func(tx *gorm.DB) error {
	//	if len(dto.Role) > 0 {
	//		if err := tx.Unscoped().Delete(&models.MenuRole{}, "menu_id =?", id).Error; err != nil {
	//			return err
	//		}
	//	}
	//	listMenuRole := make([]models.MenuRole, len(dto.Role))
	//	for k, v := range dto.Role {
	//		listMenuRole[k] = models.MenuRole{
	//			RoleID: v,
	//			MenuID: id,
	//		}
	//	}
	//	if err := tx.Model(&models.MenuRole{}).Create(listMenuRole).Error; err != nil {
	//		return err
	//	}
	//
	//	data = models.Menu{
	//		ID:          id,
	//		ParentID:    dto.ParentID,
	//		MenuTitle:   dto.MenuTitle,
	//		Slug:        dto.Slug,
	//		Url:         dto.Url,
	//		Icon:        dto.Icon,
	//		MenuOrder:   dto.MenuOrder,
	//		Description: dto.Description,
	//		IsActive:    dto.IsActive,
	//	}
	//
	//	if err := tx.UpdateColumns(&data).Error; err != nil {
	//		return err
	//	}
	//	return nil
	//}); err != nil {
	//	return models.Menu{}, err
	//}

	return data, err
}

func (s *MenuService) GetDbInstance() *gorm.DB {
	return s.MenuRepository.DbInstance()
}

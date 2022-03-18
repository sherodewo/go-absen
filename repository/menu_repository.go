package repository

import (
	"github.com/kreditplus/scorepro/models"
	"gorm.io/gorm"
)

type MenuRepository interface {
	FindAll() ([]models.Menu, error)
	FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.Menu, error)
	FindById(id string) (models.Menu, error)
	FindWhere(map[string]interface{}) (models.Menu, error)
	Save(models.Menu) (models.Menu, error)
	//Update(models.Menu) (models.Menu, error)
	Delete(models.Menu) error
	Count() (int64, error)
	CountWhere(operation string, keyVal map[string]interface{}) (int64, error)
	DbInstance() *gorm.DB
}

type menuRepository struct {
	*gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{DB: db}
}

func (r menuRepository) FindAll() ([]models.Menu, error) {
	var entities []models.Menu
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r menuRepository) FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.Menu, error) {
	var entity []models.Menu
	q := r.DB.Order("[order]" + " " + orderType).Limit(limit).Offset(offset)

	for k, v := range keyVal {
		switch operation {
		case "and":
			q = q.Where(k, v)
		case "or":
			q = q.Or(k, v)
		}
	}

	err := q.Find(&entity).Error
	return entity, err
}

func (r menuRepository) FindById(id string) (models.Menu, error) {
	var entity models.Menu
	err := r.DB.Where("user_type = ?", id).First(&entity).Error
	return entity, err
}

func (r menuRepository) FindWhere(keyVal map[string]interface{}) (models.Menu, error) {
	var entity models.Menu
	q := r.DB
	for k, v := range keyVal {
		q = q.Where(k, v)
	}
	err := q.First(&entity).Error
	return entity, err
}

func (r menuRepository) Save(entity models.Menu) (models.Menu, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

//func (r menuRepository) Update(entity models.Menu) (models.Menu, error) {
//	err := r.DB.Model(&models.Menu{ID: entity.ID}).UpdateColumns(&entity).Error
//	return entity, err
//}

func (r menuRepository) Delete(entity models.Menu) error {
	return r.DB.Delete(&entity).Error
}

func (r menuRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Table("web_menu").Count(&count).Error
	return count, err
}

func (r menuRepository) CountWhere(operation string, keyVal map[string]interface{}) (int64, error) {
	var count int64
	q := r.DB.Model(&models.Menu{})
	for k, v := range keyVal {
		switch operation {
		case "and":
			q = q.Where(k, v)
		case "or":
			q = q.Or(k, v)
		}
	}

	err := q.Count(&count).Error
	return count, err
}

func (r menuRepository) DbInstance() *gorm.DB {
	return r.DB
}

package repository

import (
	"github.com/kreditplus/scorepro/models"
	"gorm.io/gorm"
)

type ConfigRepository interface {
	FindAll() ([]models.AppConfig, error)
	FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.AppConfig, error)
	FindById(id string) (models.AppConfig, error)
	FindWhere(map[string]interface{}) (models.AppConfig, error)
	Save(models.AppConfig) (models.AppConfig, error)
	Update(models.AppConfig) (models.AppConfig, error)
	Delete(models.AppConfig) error
	Count() (int64, error)
	CountWhere(operation string, keyVal map[string]interface{}) (int64, error)
	DbInstance() *gorm.DB
}

type configRepository struct {
	*gorm.DB
}

func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepository{DB: db}
}

func (r configRepository) FindAll() ([]models.AppConfig, error) {
	var entities []models.AppConfig
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r configRepository) FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.AppConfig, error) {
	var entity []models.AppConfig
	q := r.DB.Order(orderBy + " " + orderType).Limit(limit).Offset(offset)

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

func (r configRepository) FindById(id string) (models.AppConfig, error) {
	var entity models.AppConfig
	err := r.DB.Where("id = ?", id).First(&entity).Error
	return entity, err
}

func (r configRepository) FindWhere(keyVal map[string]interface{}) (models.AppConfig, error) {
	var entity models.AppConfig
	q := r.DB
	for k, v := range keyVal {
		q = q.Where(k, v)
	}
	err := q.First(&entity).Error
	return entity, err
}

func (r configRepository) Save(entity models.AppConfig) (models.AppConfig, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r configRepository) Update(entity models.AppConfig) (models.AppConfig, error) {
	err := r.DB.Model(&models.AppConfig{ID: entity.ID}).UpdateColumns(&entity).Error
	return entity, err
}

func (r configRepository) Delete(entity models.AppConfig) error {
	return r.DB.Delete(&entity).Error
}

func (r configRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Table("web_menu").Count(&count).Error
	return count, err
}

func (r configRepository) CountWhere(operation string, keyVal map[string]interface{}) (int64, error) {
	var count int64
	q := r.DB.Model(&models.AppConfig{})
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

func (r configRepository) DbInstance() *gorm.DB {
	return r.DB
}

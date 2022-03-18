package repository

import (
	"github.com/kreditplus/scorepro/models"
	"gorm.io/gorm"
)

type ExperianRepository interface {
	FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.ExperianScoring, error)
	FindById(id string) (models.ExperianScoring, error)
	Save(models.ExperianScoring) (models.ExperianScoring, error)
	Update(models.ExperianScoring) (models.ExperianScoring, error)
	Delete(models.ExperianScoring) error
	Count() (int64, error)
	CountWhere(operation string, keyVal map[string]interface{}) (int64, error)
	DbInstance() *gorm.DB
}

type experianRepository struct {
	*gorm.DB
}

func NewExperianRepository(db *gorm.DB) ExperianRepository {
	return &experianRepository{
		DB: db,
	}
}

func (r experianRepository) FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.ExperianScoring, error) {
	var entity []models.ExperianScoring
	q := r.DB.Order(orderBy + " desc").Limit(limit).Offset(offset)

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

func (r experianRepository) FindById(id string) (models.ExperianScoring, error) {
	var entity models.ExperianScoring
	err := r.DB.Where("[experian_scoring].id = ?", id).First(&entity).Error
	return entity, err
}

func (r experianRepository) Save(entity models.ExperianScoring) (models.ExperianScoring, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r experianRepository) Update(entity models.ExperianScoring) (models.ExperianScoring, error) {
	err := r.DB.Table("experian_scoring").UpdateColumns(&entity).Error
	return entity, err
}

func (r experianRepository) Delete(entity models.ExperianScoring) error {
	return r.DB.Delete(&entity).Error
}


func (r experianRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Table("experian_scoring").Count(&count).Error
	return count, err
}

func (r experianRepository) CountWhere(operation string, keyVal map[string]interface{}) (int64, error) {
	var count int64
	q := r.DB.Model(&models.ExperianScoring{})
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

func (r experianRepository) DbInstance() *gorm.DB {
	return r.DB
}
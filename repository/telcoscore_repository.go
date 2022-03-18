package repository

import (
	"github.com/kreditplus/scorepro/models"
	"gorm.io/gorm"
)

type TelcoScoreRepository interface {
	FindAll() ([]models.Experian, error)
	FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.Experian, error)
	FindById(id string) (models.Experian, error)
	FindByPhoneNumber(phoneNumber string) (models.Experian, error)
	GetExperianDummy(PhoneNumber string) (models.Dummy, error)
	GetPickleModeling(SupplierID string) (models.SupplierModel, error)
	GetPickleModelingOther(ZipCode string) (models.SupplierOther, error)
	FinalScoring(isIndosat string, experianResult string, internalResult string) (models.ExperianScoring, error)
	IndosatValidate() (models.AppConfig, error)
	CheckConfig() (models.AppConfig, bool)
	FindWhere(email string) (models.Experian, error)
	Save(models.Experian) (models.Experian, error)
	Update(models.Experian) (models.Experian, error)
	Delete(models.Experian) error
	Count() (int64, error)
	CountWhere(operation string, keyVal map[string]interface{}) (int64, error)
	DbInstance() *gorm.DB
}

type telcoScoreRepository struct {
	*gorm.DB
}

func NewTelcoScoreRepository(db *gorm.DB) TelcoScoreRepository {
	return &telcoScoreRepository{
		DB: db,
	}
}

func (r telcoScoreRepository) FindAll() ([]models.Experian, error) {
	var entities []models.Experian
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r telcoScoreRepository) FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.Experian, error) {
	var entity []models.Experian
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

func (r telcoScoreRepository) FindById(id string) (models.Experian, error) {
	var entity models.Experian
	err := r.DB.Where("[experian].experian_id = ?", id).First(&entity).Error
	return entity, err
}

func (r telcoScoreRepository) FindByPhoneNumber(phoneNumber string) (models.Experian, error) {
	var entity models.Experian
	err := r.DB.Where("[experian].phone_number = ?", phoneNumber).Order("created_at desc").First(&entity).Error
	return entity, err
}

func (r telcoScoreRepository) GetExperianDummy(PhoneNumber string) (models.Dummy, error) {
	var entity models.Dummy
	err := r.DB.Where("experian_dummy.phone_number = ?", PhoneNumber).Find(&entity).Error
	return entity, err
}

func (r telcoScoreRepository) FindWhere(email string) (models.Experian, error) {
	var entity models.Experian
	err := r.DB.Where("email = ?",email).Find(&entity).Error
	return entity, err
}

func (r telcoScoreRepository) Save(entity models.Experian) (models.Experian, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r telcoScoreRepository) Update(entity models.Experian) (models.Experian, error) {
	err := r.DB.Model(models.Experian{ExperianID: entity.ExperianID}).UpdateColumns(&entity).Error
	return entity, err
}

func (r telcoScoreRepository) Delete(entity models.Experian) error {
	return r.DB.Delete(&entity).Error
}

func (r telcoScoreRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Table("experian").Count(&count).Error
	return count, err
}

func (r telcoScoreRepository) CountWhere(operation string, keyVal map[string]interface{}) (int64, error) {
	var count int64
	q := r.DB.Model(&models.Experian{})
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

func (r telcoScoreRepository) IndosatValidate() (models.AppConfig, error) {
	var app models.AppConfig

	err := r.DB.Where("group_name = 'indosat' AND is_active = 1").
		Find(&app).Error
	return app,err
}

func (r telcoScoreRepository) CheckConfig() (models.AppConfig, bool) {
	var app models.AppConfig
	valid := false
	err := r.DB.Where("group_name = 'use_scoring'").
		Find(&app)
	if err.RowsAffected >= 1 {
		valid = true
	}
	return app,valid
}

func (r telcoScoreRepository) GetPickleModeling(SupplierID string) (models.SupplierModel, error) {
	var entity models.SupplierModel
	err := r.DB.Where("supplier_id = ?", SupplierID).Find(&entity).Error
	return entity, err
}

func (r telcoScoreRepository) GetPickleModelingOther(ZipCode string) (models.SupplierOther, error) {
	var entity models.SupplierOther
	err := r.DB.Where("supplier_others.zip_code = ?", ZipCode).Find(&entity).Error
	return entity, err
}

func (r telcoScoreRepository) FinalScoring(isIndosat string,experianResult string, internalResult string) (models.ExperianScoring, error) {
	var entity models.ExperianScoring
	err := r.DB.Where("is_indosat = ? AND experian = ? AND internal = ? AND is_active = 1", isIndosat,experianResult,internalResult).
		Find(&entity).Error
	return entity, err
}

func (r telcoScoreRepository) DbInstance() *gorm.DB {
	return r.DB
}
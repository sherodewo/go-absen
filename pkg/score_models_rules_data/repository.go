package score_models_rules_data

import (
	"gorm.io/gorm"
)

type ScoreModelsRulesDataRepository interface {
	FindAll() ([]ScoreModelsRulesData, error)
	FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]ScoreModelsRulesData, error)
	FindById(id string) (ScoreModelsRulesData, error)
	Save(ScoreModelsRulesData) (ScoreModelsRulesData, error)
	Update(ScoreModelsRulesData) (ScoreModelsRulesData, error)
	Delete(ScoreModelsRulesData) error
	Count() (int64, error)
	CountWhere(operation string, keyVal map[string]interface{}) (int64, error)
	DbInstance() *gorm.DB
}

type supplierRepository struct {
	*gorm.DB
}

func NewScoreModelsRulesDataRepository(db *gorm.DB) ScoreModelsRulesDataRepository {
	return &supplierRepository{DB: db}
}

func (r supplierRepository) FindAll() ([]ScoreModelsRulesData, error) {
	var entities []ScoreModelsRulesData
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r supplierRepository) FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]ScoreModelsRulesData, error) {
	var entity []ScoreModelsRulesData
	q := r.DB.Table("score_models_rules_data a").
		Select("a.*, b.name as score_generators").
		Joins("LEFT JOIN score_generators b on a.score_generators = b.id").Order(orderBy + " " + orderType).Limit(limit).Offset(offset)

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

func (r supplierRepository) FindById(id string) (ScoreModelsRulesData, error) {
	var entity ScoreModelsRulesData
	err := r.DB.Where("[score_models_rules_data].id = ?", id).First(&entity).Error
	return entity, err
}

func (r supplierRepository) Save(entity ScoreModelsRulesData) (ScoreModelsRulesData, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r supplierRepository) Update(entity ScoreModelsRulesData) (ScoreModelsRulesData, error) {
	err := r.DB.Model(ScoreModelsRulesData{ID: entity.ID}).UpdateColumns(&entity).Error
	return entity, err
}

func (r supplierRepository) Delete(entity ScoreModelsRulesData) error {
	return r.DB.Delete(&entity).Error
}

func (r supplierRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Table("score_models_rules_data").Count(&count).Error
	return count, err
}

func (r supplierRepository) CountWhere(operation string, keyVal map[string]interface{}) (int64, error) {
	var count int64
	q := r.DB.Model(&ScoreModelsRulesData{})
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

func (r supplierRepository) DbInstance() *gorm.DB {
	return r.DB
}
package scorepro

import (
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/pkg/data_inquiry"
	"github.com/kreditplus/scorepro/pkg/experian_inquiry"
	"github.com/kreditplus/scorepro/pkg/metric_combination"
	"github.com/kreditplus/scorepro/pkg/score_generator"
	"github.com/kreditplus/scorepro/pkg/scorepro_requestor"
	"github.com/kreditplus/scorepro/pkg/scorepro_score"
	"gorm.io/gorm"
)

type ScoreproRepository interface {
	GetScoreModelingPickle(id string) (score_generator.ScoreGenerator, error)
	StoreDataInquiry(request data_inquiry.DataInquiry) (data_inquiry.DataInquiry, error)
	AppConf(groupName string) (models.AppConfig, error)
	GetExperianDummy(groupName string) (models.Dummy, error)
	Save(request *experian_inquiry.Experian) (*experian_inquiry.Experian, error)
	ScoreproScore(number float64, flag string) (scorepro_score.ScoreproScore, error)
	ScoreproScoreWithTrxType(number float64, flag string, trxType string) (scorepro_score.ScoreproScore, error)
	CheckRequestor(id string) (scorepro_requestor.ScoreproRequestor, error)
	MetricCombination(isIndosat string, experian string, internal string,lob string) (metric_combination.MetricCombination, error)
	CountWhere(operation string, keyVal map[string]interface{}) (int64, error)
	FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]experian_inquiry.Experian, error)
	Count() (int64, error)
	FindById(id string) (experian_inquiry.Experian, error)

	DbInstance() *gorm.DB
}

type scoreproRepository struct {
	*gorm.DB
}

func NewTelcoScoreRepository(db *gorm.DB) ScoreproRepository {
	return &scoreproRepository{
		DB: db,
	}
}

func (r scoreproRepository) GetScoreModelingPickle(id string) (data score_generator.ScoreGenerator, err error) {
	err = r.DB.Find(&data, score_generator.ScoreGenerator{ID: id}).Error

	return data, err
}

func (r scoreproRepository) StoreDataInquiry(request data_inquiry.DataInquiry) (data_inquiry.DataInquiry, error) {
	err := r.DB.Create(&request).Error

	return request, err
}

func (r scoreproRepository) AppConf(groupName string) (models.AppConfig, error) {
	var app models.AppConfig

	err := r.DB.Where("group_name = ? AND is_active = 1", groupName).
		Find(&app).Error
	return app, err
}

func (r scoreproRepository) MetricCombination(isIndosat string, experian string, internal string, lob string) (metric_combination.MetricCombination, error) {
	var entity metric_combination.MetricCombination
	err := r.DB.Model(&entity).
		Find(&entity, "is_indosat = ? AND experian = ? AND internal = ? AND is_active  = 1 AND lob = ?", isIndosat, experian, internal,lob).
		Error

	return entity, err
}

func (r scoreproRepository) GetExperianDummy(PhoneNumber string) (models.Dummy, error) {
	var entity models.Dummy
	err := r.DB.Where("experian_dummy.phone_number = ?", PhoneNumber).First(&entity)
	return entity, err.Error
}

func (r scoreproRepository) Save(entity *experian_inquiry.Experian) (*experian_inquiry.Experian, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r scoreproRepository) ScoreproScore(number float64, flag string) (scorepro_score.ScoreproScore, error) {
	var entity scorepro_score.ScoreproScore
	err := r.DB.Find(&entity, "score_start <= ? AND score_end >= ? AND threshold_owners = ?", number, number, flag).Error

	return entity, err
}

func (r scoreproRepository) ScoreproScoreWithTrxType(number float64, flag string, transactionType string) (scorepro_score.ScoreproScore, error) {
	var entity scorepro_score.ScoreproScore
	err := r.DB.Find(&entity, "score_start <= ? AND score_end >= ? AND threshold_owners = ? AND transaction_type = ?", number, number, flag, transactionType).Error

	return entity, err
}

func (r scoreproRepository) CheckRequestor(ID string) (scorepro_requestor.ScoreproRequestor, error) {
	var entity scorepro_requestor.ScoreproRequestor
	err := r.DB.Find(&entity, "ID = ? ", ID).Error

	return entity, err
}

func (r scoreproRepository) FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]experian_inquiry.Experian, error) {
	var entity []experian_inquiry.Experian
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

func (r scoreproRepository) CountWhere(operation string, keyVal map[string]interface{}) (int64, error) {
	var count int64
	q := r.DB.Model(&experian_inquiry.Experian{})
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

func (r scoreproRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Table("scorepro_experian_inquiry").Count(&count).Error
	return count, err
}

func (r scoreproRepository) FindById(id string) (experian_inquiry.Experian, error) {
	var entity experian_inquiry.Experian
	err := r.DB.Where("[scorepro_experian_inquiry].id = ?", id).First(&entity).Error
	return entity, err
}

func (r scoreproRepository) DbInstance() *gorm.DB {
	return r.DB
}

package score_generator

import (
	"github.com/kreditplus/scorepro/pkg/score_generator_meta_data"
	"github.com/kreditplus/scorepro/pkg/score_models_rules"
	"gorm.io/gorm"
)

type ScoreRepository interface {
	CountWhere(operation string, keyVal map[string]interface{}) (int64, error)
	FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]ScoreGenerator, error)
	Count() (int64, error)
	DbInstance() *gorm.DB
	StoreScoreGenerator(generator ScoreGenerator) (ScoreGenerator,error)
	StoreScoreGeneratorMetadata(generator score_generator_meta_data.ScoreGeneratorsMetadata) (score_generator_meta_data.ScoreGeneratorsMetadata,error)
	GetScoreModelRules() ([]score_models_rules.ScoreModelsRules,error)
	GetAllScoreGenerator() ([]ScoreGenerator,error)
	CreateTable(dbName string,query string) error
}

type scoreRepository struct {
	*gorm.DB
}

func NewScoreRepository(db *gorm.DB) ScoreRepository {
	return &scoreRepository{
		DB: db,
	}
}

func (r scoreRepository) DbInstance() *gorm.DB {
	return r.DB
}

func (r scoreRepository) FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]ScoreGenerator, error) {
	var entity []ScoreGenerator
	q := r.DB.Select("a.*, b.[key] as score_models_rules").
		Table("score_generators a").
		Joins("LEFT JOIN score_models_rules b on a.score_models_rules = b.id").
		Order("a."+orderBy + " desc").Limit(limit).Offset(offset)

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

func (r scoreRepository) CountWhere(operation string, keyVal map[string]interface{}) (int64, error) {
	var count int64
	q := r.DB.Model(&ScoreGenerator{})
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

func (r scoreRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Table("scorepro_experian_inquiry").Count(&count).Error
	return count, err
}

func (r scoreRepository) StoreScoreGenerator(entity ScoreGenerator) (ScoreGenerator,error) {
	err := r.DB.Create(&entity).Error
	return entity,err
}

func (r scoreRepository) StoreScoreGeneratorMetadata(entity score_generator_meta_data.ScoreGeneratorsMetadata) (score_generator_meta_data.ScoreGeneratorsMetadata,error) {
	err := r.DB.Create(&entity).Error
	return entity,err
}

func (r scoreRepository) GetScoreModelRules()([]score_models_rules.ScoreModelsRules,error) {
	var data []score_models_rules.ScoreModelsRules
	err := r.DB.Find(&data).Error
	return data,err
}

func (r scoreRepository) CreateTable(dbName string,query string) error {
	rawQuery := ` CREATE TABLE [dbo].[`+dbName+`]
	 ([id] varchar(50) COLLATE SQL_Latin1_General_CP1_CI_AS  NOT NULL
	 `+query+`,
	 [score] int NULL,
	 [created_at] datetime2(7)  NULL,
	 CONSTRAINT [PK__`+dbName+`] PRIMARY KEY CLUSTERED ([id]));
	 CREATE NONCLUSTERED INDEX [idx_`+dbName+`_transaction_id] ON [dbo].[`+dbName+`] ([transaction_id] ASC);
	 CREATE NONCLUSTERED INDEX [idx_`+dbName+`_created_at] ON [dbo].[`+dbName+`] ([created_at] ASC);
`
	err := r.DB.Exec(rawQuery).Error

	return err
}

func (r scoreRepository) GetAllScoreGenerator()([]ScoreGenerator,error) {
	var data []ScoreGenerator
	err := r.DB.Find(&data).Error
	return data,err
}
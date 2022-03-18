package score_models_rules_data

import (
	"gorm.io/gorm"
	"time"
)

type ScoreModelsRulesDataService struct {
	ScoreModelsRulesDataRepository ScoreModelsRulesDataRepository
}

func NewScoreModelsRulesDataService(repository ScoreModelsRulesDataRepository) *ScoreModelsRulesDataService {
	return &ScoreModelsRulesDataService{
		ScoreModelsRulesDataRepository: repository,
	}
}

func (s *ScoreModelsRulesDataService) FindAllScoreModelsRulesDatas() (*[]ScoreModelsRulesData, error) {
	data, err := s.ScoreModelsRulesDataRepository.FindAll()
	return &data, err
}

func (s *ScoreModelsRulesDataService) FindScoreModelsRulesDataById(id string) (*ScoreModelsRulesData, error) {
	data, err := s.ScoreModelsRulesDataRepository.FindById(id)

	return &data, err
}

func (s *ScoreModelsRulesDataService) SaveScoreModelsRulesData(dto ScoreModelsRulesDataReq) (*ScoreModelsRulesData, error) {
	entity := ScoreModelsRulesData{
		Key:            dto.Key,
		Value:          dto.Value,
		Description:    dto.Description,
		CreatedAt:      time.Now(),
		ScoreGenerator: dto.ScoreGenerator,
	}

	data, err := s.ScoreModelsRulesDataRepository.Save(entity)
	return &data, err
}

func (s *ScoreModelsRulesDataService) UpdateScoreModelsRulesData(id string, dto ScoreModelsRulesDataReq) (*ScoreModelsRulesData, error) {
	entity := ScoreModelsRulesData{
		ID:             id,
		Key:            dto.Key,
		Value:          dto.Value,
		Description:    dto.Description,
		ScoreGenerator: dto.ScoreGenerator,
	}

	data, err := s.ScoreModelsRulesDataRepository.Update(entity)

	return &data, err
}

func (s *ScoreModelsRulesDataService) DeleteScoreModelsRulesData(id string) error {
	entity := ScoreModelsRulesData{
		ID: id,
	}
	err := s.ScoreModelsRulesDataRepository.Delete(entity)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *ScoreModelsRulesDataService) QueryDatatable(searchValue string, orderType string, orderBy string, limit int, offset int) (
	recordTotal int64, recordFiltered int64, data []ScoreModelsRulesData, err error) {
	recordTotal, err = s.ScoreModelsRulesDataRepository.Count()

	if searchValue != "" {
		recordFiltered, err = s.ScoreModelsRulesDataRepository.CountWhere("or", map[string]interface{}{
			"key LIKE ?": "%" + searchValue + "%",
			"value LIKE ?":   "%" + searchValue + "%",
		})

		data, err = s.ScoreModelsRulesDataRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
			"key LIKE ?": "%" + searchValue + "%",
			"value LIKE ?":   "%" + searchValue + "%",
		})
		return recordTotal, recordFiltered, data, err
	}

	recordFiltered, err = s.ScoreModelsRulesDataRepository.CountWhere("or", map[string]interface{}{
		"1 =?": 1,
	})

	data, err = s.ScoreModelsRulesDataRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
		"1= ?": 1,
	})
	return recordTotal, recordFiltered, data, err
}

func (s *ScoreModelsRulesDataService) GetDbInstance() *gorm.DB {
	return s.ScoreModelsRulesDataRepository.DbInstance()
}

package service

import (
	"github.com/kreditplus/scorepro/dto"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/repository"
	"gorm.io/gorm"
)

type ExperianService struct {
	ExperianRepository repository.ExperianRepository
}

func NewExperianService(repository repository.ExperianRepository) *ExperianService {
	return &ExperianService{
		ExperianRepository: repository,
	}
}

func (s *ExperianService) FindExperianById(id string) (*models.ExperianScoring, error) {
	data, err := s.ExperianRepository.FindById(id)

	return &data, err
}

func (s *ExperianService) SaveExperian(req dto.ExperianScoringDto) (*models.ExperianScoring, error) {
	entity := models.ExperianScoring{
		IsIndosat:  req.IsIndosat,
		Experian:   req.Experian,
		Internal:   req.Internal,
		ScoreLos:   req.ScoreLos,
		FinalScore: req.FinalScore,
		Notes:      req.Notes,
	}
	data, err := s.ExperianRepository.Save(entity)
	return &data, err
}

func (s *ExperianService) UpdateExperian(id string, req dto.ExperianScoringUpdateDto) (*models.ExperianScoring, error) {
	entity := models.ExperianScoring{
		ID:         id,
		IsIndosat:  req.IsIndosat,
		Experian:   req.Experian,
		Internal:   req.Internal,
		ScoreLos:   req.ScoreLos,
		FinalScore: req.FinalScore,
		Notes:      req.Notes,
	}
	data, err := s.ExperianRepository.Update(entity)

	return &data, err
}

func (s *ExperianService) DeleteExperian(id string) error {
	entity := models.ExperianScoring{
		ID: id,
	}
	err := s.ExperianRepository.Delete(entity)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *ExperianService) QueryDatatable(searchValue string, orderType string, orderBy string, limit int, offset int) (
	recordTotal int64, recordFiltered int64, data []models.ExperianScoring, err error) {
	recordTotal, err = s.ExperianRepository.Count()

	if searchValue != "" {
		recordFiltered, err = s.ExperianRepository.CountWhere("or", map[string]interface{}{
			"experian LIKE ?": "%" + searchValue + "%",
			"internal LIKE ?":     "%" + searchValue + "%",
		})

		data, err = s.ExperianRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
			"experian LIKE ?": "%" + searchValue + "%",
			"internal LIKE ?":     "%" + searchValue + "%",
		})
		return recordTotal, recordFiltered, data, err
	}

	recordFiltered, err = s.ExperianRepository.CountWhere("or", map[string]interface{}{
		"1 =?": 1,
	})

	data, err = s.ExperianRepository.FindAllWhere("or", orderType, "id", limit, offset, map[string]interface{}{
		"1= ?": 1,
	})
	return recordTotal, recordFiltered, data, err
}

func (s *ExperianService) GetDbInstance() *gorm.DB {
	return s.ExperianRepository.DbInstance()
}


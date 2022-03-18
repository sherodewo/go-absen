package service

import (
	"github.com/allegro/bigcache/v3"
	"github.com/kreditplus/scorepro/dto"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/repository"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var cache *bigcache.BigCache

type ConfigService struct {
	ConfigRepository repository.ConfigRepository
}

func NewConfigService(repository repository.ConfigRepository) *ConfigService {
	return &ConfigService{
		ConfigRepository: repository,
	}
}

func (s *ConfigService) QueryDatatable(searchValue string, orderType string, orderBy string, limit int, offset int) (
	recordTotal int64, recordFiltered int64, data []models.AppConfig, err error) {
	recordTotal, err = s.ConfigRepository.Count()

	if searchValue != "" {
		recordFiltered, err = s.ConfigRepository.CountWhere("or", map[string]interface{}{
			"group_name LIKE ?": "%" + searchValue + "%",
			"key LIKE ?":        "%" + searchValue + "%",
		})

		data, err = s.ConfigRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
			"group_name LIKE ?": "%" + searchValue + "%",
			"key LIKE ?":        "%" + searchValue + "%",
		})
		return recordTotal, recordFiltered, data, err
	}

	recordFiltered, err = s.ConfigRepository.CountWhere("or", map[string]interface{}{
		"1 =?": 1,
	})

	data, err = s.ConfigRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
		"1= ?": 1,
	})
	return recordTotal, recordFiltered, data, err
}

func (s *ConfigService) DeleteConfig(id string) error {
	entity := models.AppConfig{
		ID: id,
	}
	err := s.ConfigRepository.Delete(entity)
	return err
}

func (s *ConfigService) FindById(id string) (*models.AppConfig, error) {
	data, err := s.ConfigRepository.FindById(id)
	return &data, err
}

func (s *ConfigService) FindAllMenu() (*[]models.AppConfig, error) {
	data, err := s.ConfigRepository.FindAll()
	return &data, err
}

func (s *ConfigService) StoreMenu(dto dto.AppConfDto) (*models.AppConfig, error) {
	data := models.AppConfig{
		GroupName: dto.GroupName,
		Key:       dto.Key,
		Value:     dto.Value,
		IsActive:  dto.IsActive,
		CreatedAt: time.Time{},
	}
	res, err := s.ConfigRepository.Save(data)
	if err != nil {
		return &res, err
	}
	return &res, err
}

func (s *ConfigService) UpdateConfig(id string, dto dto.AppConfDto) (*models.AppConfig, error) {
	entity := models.AppConfig{
		ID:        id,
		GroupName: dto.GroupName,
		Key:       dto.Key,
		Value:     dto.Value,
		IsActive:  dto.IsActive,
		UpdatedAt: time.Time{},
	}
	data, err := s.ConfigRepository.Update(entity)
	switch data.GroupName {
	case "indosat":
		isActive := strconv.Itoa(dto.IsActive)
		cache.Set("indosat-active", []byte(isActive))
		cache.Set("indosat-validation", []byte(dto.Value))
	case "use_scoring":
		isActive := strconv.Itoa(dto.IsActive)
		cache.Set("use-scoring", []byte(isActive))
		_, _ = cache.Get("use-scoring")
	}

	return &data, err
}

func (s *ConfigService) GetDbInstance() *gorm.DB {
	return s.ConfigRepository.DbInstance()
}

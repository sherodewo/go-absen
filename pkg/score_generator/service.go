package score_generator

import (
	"fmt"
	"github.com/kreditplus/scorepro/pkg/score_generator_meta_data"
	"github.com/kreditplus/scorepro/pkg/score_models_rules"
	"mime/multipart"
	"strings"
	"time"
)

type ScoreService struct {
	ScoreRepository ScoreRepository
}

func NewScoreService(repository ScoreRepository) *ScoreService {
	return &ScoreService{
		ScoreRepository: repository,
	}
}

func (s *ScoreService) QueryDatatable(searchValue string, orderType string, orderBy string, limit int, offset int) (
	recordTotal int64, recordFiltered int64, data []ScoreGenerator, err error) {
	recordTotal, err = s.ScoreRepository.Count()

	if searchValue != "" {
		recordFiltered, err = s.ScoreRepository.CountWhere("or", map[string]interface{}{
			"ProspectID LIKE ?": "%" + searchValue + "%",
			"result LIKE ?":     "%" + searchValue + "%",
		})

		data, err = s.ScoreRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
			"ProspectID LIKE ?": "%" + searchValue + "%",
			"result LIKE ?":     "%" + searchValue + "%",
		})
		return recordTotal, recordFiltered, data, err
	}

	recordFiltered, err = s.ScoreRepository.CountWhere("or", map[string]interface{}{
		"1 =?": 1,
	})

	data, err = s.ScoreRepository.FindAllWhere("or", orderType, "created_at", limit, offset, map[string]interface{}{
		"1= ?": 1,
	})
	return recordTotal, recordFiltered, data, err
}

func (s *ScoreService) StoreScoreGenerator(req WizardDto) (*ScoreGenerator, error) {
	entity := ScoreGenerator{
		Name:                 req.Name,
		Description:          req.Description,
		Endpoint:             req.Endpoint,
		FilePickle:           req.FilePickle,
		SaveResultTo:         req.SaveResultTo,
		SaveResultObjectName: req.SaveResultName,
		CreatedAt:            time.Now(),
		ScoreModelsRules:     req.ScoreModelsRules,
	}

	data, err := s.ScoreRepository.StoreScoreGenerator(entity)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *ScoreService) StoreScoreGeneratorMetadata(req WizardDto,id string) (*score_generator_meta_data.ScoreGeneratorsMetadata, error) {
	var response score_generator_meta_data.ScoreGeneratorsMetadata

	cleaned := strings.Replace(req.ParamRequestPayload, ",", " ", -1)
	reqPayload := strings.Fields(cleaned)

	cleaned = strings.Replace(req.FieldPickle, ",", " ", -1)
	fieldPickle := strings.Fields(cleaned)

	cleaned = strings.Replace(req.DataType, ",", " ", -1)
	dataType := strings.Fields(cleaned)

	cleaned = strings.Replace(req.Length, ",", " ", -1)
	length := strings.Fields(cleaned)

	cleaned = strings.Replace(req.IsRequired, ",", " ", -1)
	isRequired := strings.Fields(cleaned)

	for i:=0; i<len(reqPayload); i++{
		if dataType[i] == "int" {
			dataType[i] = "INTEGER"
		}else {
			dataType[i] = "VARCHAR"
		}

		data := score_generator_meta_data.ScoreGeneratorsMetadata{
			ParamRequestPayload: reqPayload[i],
			FieldPickle:         fieldPickle[i],
			DataType:            dataType[i],
			Length:              length[i],
			IsRequired:          isRequired[i],
			CreatedAt:           time.Now(),
			ScoreGenerators:     id,
		}

		res, err := s.ScoreRepository.StoreScoreGeneratorMetadata(data)
		if err != nil {
			return nil, err
		}
		response = res
	}



	return &response, nil
}

func (s *ScoreService) UploadFilePickle (name string, fileName multipart.File ,replace string) (int,interface{},error) {
	httpClient := NewUploadHTTPClient()
	status ,message, err := httpClient.UploadFilePickle(name,fileName,replace)

	return status,message, err
}

func (s *ScoreService) GetScoreModelsRules() ([]score_models_rules.ScoreModelsRules,error) {
	data,err := s.ScoreRepository.GetScoreModelRules()

	return data,err
}

func (s *ScoreService) CreateTable(req WizardDto) error{
	cleaned := strings.Replace(req.ParamRequestPayload, ",", " ", -1)
	reqPayload := strings.Fields(cleaned)

	cleaned = strings.Replace(req.DataType, ",", " ", -1)
	dataType := strings.Fields(cleaned)

	cleaned = strings.Replace(req.Length, ",", " ", -1)
	length := strings.Fields(cleaned)

	cleaned = strings.Replace(req.IsRequired, ",", " ", -1)
	isRequired := strings.Fields(cleaned)

	var query string

	for i:=0; i<len(reqPayload); i++ {
		if isRequired[i] == "TRUE" && dataType[i] != "int" && dataType[i] != "bigint"{
			query += fmt.Sprintf(",[%s] %s (%s) NOT NULL",reqPayload[i],dataType[i],length[i])
		}else if dataType[i] != "int" && dataType[i] != "bigint"{
			query += fmt.Sprintf(",[%s] %s (%s) NULL",reqPayload[i],dataType[i],length[i])
		}else if dataType[i] == "int" && isRequired[i] == "TRUE" || dataType[i] == "bigint" && isRequired[i] == "TRUE" {
			query += fmt.Sprintf(",[%s] %s NOT NULL",reqPayload[i],dataType[i])
		}else {
			query += fmt.Sprintf(",[%s] %s NOT NULL",reqPayload[i],dataType[i])
		}

	}

		err :=  s.ScoreRepository.CreateTable(req.SaveResultName,query)

		return err
}


func (s *ScoreService) GetAllScoreGenerator() ([]ScoreGenerator,error) {
	data,err := s.ScoreRepository.GetAllScoreGenerator()

	return data,err
}
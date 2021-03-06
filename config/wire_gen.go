// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package config

import (
	"github.com/kreditplus/scorepro/controllers"
	"github.com/kreditplus/scorepro/pkg/score_generator"
	"github.com/kreditplus/scorepro/pkg/scorepro"
	"github.com/kreditplus/scorepro/pkg/score_models_rules_data"
	"github.com/kreditplus/scorepro/repository"
	"github.com/kreditplus/scorepro/service"
	"gorm.io/gorm"
)

// Injectors from di.go:

func InjectUserController(db *gorm.DB) controllers.UserController {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)
	return userController
}

func InjectAuthController(db *gorm.DB) controllers.AuthController {
	userRepository := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepository)
	authController := controllers.NewAuthController(authService)
	return authController
}

func InjectMenuController(db *gorm.DB) controllers.MenuController {
	menuRepository := repository.NewMenuRepository(db)
	menuService := service.NewMenuService(menuRepository)
	menuController := controllers.NewMenuController(menuService)
	return menuController
}

func InjectConfigController(db *gorm.DB) controllers.ConfigController {
	repository := repository.NewConfigRepository(db)
	service := service.NewConfigService(repository)
	controller := controllers.NewConfigController(service)
	return controller
}

func InjectTelcoScoreController(db *gorm.DB) controllers.TelcoScoreController {
	repository := repository.NewTelcoScoreRepository(db)
	service := service.NewTelcoScoreService(repository)
	controller := controllers.NewTelcoScoreController(service)
	return controller
}

func InjectTExperianScoringController(db *gorm.DB) controllers.ExperianController {
	repository := repository.NewExperianRepository(db)
	service := service.NewExperianService(repository)
	controller := controllers.NewExperianController(service)
	return controller
}

func InjectScoreproController(db *gorm.DB) scorepro.ScoreproController {
	repository := scorepro.NewTelcoScoreRepository(db)
	service := scorepro.NewScoreproService(repository)
	controller := scorepro.NewTelcoScoreController(service)
	return controller
}

func InjectGeneratorController(db *gorm.DB) score_generator.ScoreController {
	repository := score_generator.NewScoreRepository(db)
	service := score_generator.NewScoreService(repository)
	controller := score_generator.NewScoreController(service)
	return controller
}

func InjectKmbScoreproController(db *gorm.DB) scorepro.KmbScoreproController {
	repository := scorepro.NewTelcoScoreRepository(db)
	service := scorepro.NewScoreproService(repository)
	controller := scorepro.NewKmbScoreproController(service)
	return controller
}

func InjectWgScoreproController(db *gorm.DB) scorepro.WgScoreproController {
	repository := scorepro.NewTelcoScoreRepository(db)
	service := scorepro.NewScoreproService(repository)
	controller := scorepro.NewWgScoreproController(service)
	return controller
}


func InjectSupplierController(db *gorm.DB) score_models_rules_data.ScoreModelsRulesDataController {
	repository := score_models_rules_data.NewScoreModelsRulesDataRepository(db)
	service := score_models_rules_data.NewScoreModelsRulesDataService(repository)
	controller := score_models_rules_data.NewScoreModelsRulesDataController(service)
	return controller
}


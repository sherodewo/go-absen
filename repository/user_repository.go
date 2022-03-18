package repository

import (
	"github.com/kreditplus/scorepro/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.User, error)
	FindById(id string) (models.User, error)
	FindWhere(email string) (models.User, error)
	FindBranch(id string) (models.UserDetail, error)
	Save(models.User) (models.User, error)
	Update(models.User) (models.User, error)
	Delete(models.User) error
	Count() (int64, error)
	CountWhere(operation string, keyVal map[string]interface{}) (int64, error)
	DbInstance() *gorm.DB
}

type userRepository struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r userRepository) FindAll() ([]models.User, error) {
	var entities []models.User
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r userRepository) FindAllWhere(operation string, orderType string, orderBy string, limit int, offset int, keyVal map[string]interface{}) ([]models.User, error) {
	var entity []models.User
	q := r.DB.Order(orderBy + " " + orderType).Limit(limit).Offset(offset)

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

func (r userRepository) FindById(id string) (models.User, error) {
	var entity models.User
	err := r.DB.Where("[user].user_id = ?", id).First(&entity).Error
	return entity, err
}

func (r userRepository) FindWhere(email string) (models.User, error) {
	var entity models.User
	err := r.DB.Where("email = ?",email).Find(&entity).Error
	return entity, err
}

func (r userRepository) FindBranch(id string) (models.UserDetail, error) {
	var entity models.UserDetail
	err := r.DB.Where("user_id = ?",id).Find(&entity).Error
	return entity, err
}

func (r userRepository) Save(entity models.User) (models.User, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r userRepository) Update(entity models.User) (models.User, error) {
	err := r.DB.Model(models.User{UserID: entity.UserID}).UpdateColumns(&entity).Error
	return entity, err
}

func (r userRepository) Delete(entity models.User) error {
	return r.DB.Delete(&entity).Error
}

func (r userRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Table("user").Count(&count).Error
	return count, err
}

func (r userRepository) CountWhere(operation string, keyVal map[string]interface{}) (int64, error) {
	var count int64
	q := r.DB.Model(&models.User{})
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

func (r userRepository) DbInstance() *gorm.DB {
	return r.DB
}
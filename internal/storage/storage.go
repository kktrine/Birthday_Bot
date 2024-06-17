package storage

import (
	"birthday_bot/internal/model"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	db *gorm.DB
}

func New(cfg string) *Storage {
	db, err := gorm.Open(postgres.Open(cfg), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("couldn't connect to database: " + err.Error())
	}
	db = db.Debug()

	err = db.AutoMigrate(&model.Employee{}, &model.Subscription{}, &model.User{})
	if err != nil {
		panic("failed to migrate tables: " + err.Error())
	}
	return &Storage{db: db}
}

func (d Storage) Stop() error {
	val, err := d.db.DB()
	if err != nil {
		return errors.New("failed to get database error: " + err.Error())
	}
	return val.Close()
}

func (d Storage) GetEmployees() (*[]model.Employee, error) {
	return nil, nil
}

func (d Storage) CreateUser(username, password string) error {
	tx := d.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	tx.Model(&model.User{}).Create(&model.User{
		Username: username,
		Password: password,
	})
	tx.Commit()
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d Storage) GetHashedPassword(username string) (string, error) {
	var hashedPassword string
	err := d.db.Model(&model.User{}).Where("username = ?", username).Select("password").Scan(&hashedPassword).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return hashedPassword, nil

}

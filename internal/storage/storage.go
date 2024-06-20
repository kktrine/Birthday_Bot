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
	var res []model.Employee
	err := d.db.Model(&model.Employee{}).Find(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return &res, nil
}

func (d Storage) CreateUser(user model.User) error {
	tx := d.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	tx.Model(&model.User{}).Create(&user)
	tx.Commit()
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//func (d Storage) AddInfo(info model.Employee) error {
//
//}

func (d Storage) GetHashedPassword(username string) (string, *int64, error) {
	var hashedPassword string
	var chatId *int64
	err := d.db.Model(&model.User{}).Where("username = ?", username).Select("password, chat_id").Scan(&hashedPassword).Scan(chatId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, nil
		}
		return "", nil, err
	}
	return hashedPassword, chatId, nil

}

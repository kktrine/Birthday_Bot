package storage

import (
	"birthday_bot/internal/model"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
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

func (d Storage) GetHashedPassword(username string) (model.User, error) {
	var user model.User
	err := d.db.Model(&model.User{}).Where("username = ?", username).Select("password, chat_id, id").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, nil
		}
		return model.User{}, err
	}
	return user, nil
}

func (d Storage) AddInfo(info model.Employee) error {
	tx := d.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	var num int64
	tx.Model(&model.User{}).Where("id = (?)", info.UserId).Count(&num)
	if num == 0 {
		return errors.New("пользователь с этим id не найден")
	}
	query := tx.Model(&model.Employee{}).Where("user_id = (?)", info.UserId).Select("id")
	query.Count(&num)
	if num == 0 {
		err := tx.Create(&info).Error
		if err != nil {
			return err
		}
	} else {
		query.Scan(&info.Id)
		tx.Save(&info)
	}
	tx.Commit()
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d Storage) Subscribe(data model.Subscribe) error {
	tx := d.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	var num int64
	tx.Model(&model.User{}).Where("id = (?)", data.Id).Count(&num)
	if num == 0 {
		return errors.New("пользователь с этим id не найден")
	}
	tx.Model(&model.User{}).Where("id in ?", *data.SubscribeTo).Count(&num)
	if num != int64(len(*data.SubscribeTo)) {
		return errors.New("один или несколько пользователей для подписки не найдены")
	}
	requests := make([]model.Subscription, len(*data.SubscribeTo))
	for i, id := range *data.SubscribeTo {
		requests[i].UserId = *data.Id
		requests[i].SubscribedTo = id
	}
	err := tx.Create(requests).Error
	if err != nil {
		return err
	}
	tx.Commit()
	if tx.Error != nil {
		return err
	}
	return nil
}

func (d Storage) CheckBDays() *[]model.Message {
	var bdays []model.Employee
	today := time.Now().Format("2006-01-02")
	err := d.db.Where("birth = ?", today).Select("user_id, name, surname").Find(&model.Employee{}).Error
	if err != nil {
		return nil
	}
	bdaysMap := make(map[int]struct {
		name    string
		surname string
	})

	var bdayIds []int
	for _, employee := range bdays {
		bdayIds = append(bdayIds, *employee.UserId)
		bdaysMap[employee.Id] = struct {
			name    string
			surname string
		}{name: employee.Name, surname: employee.Surname}
	}
	var subscribers []model.Subscription
	err = d.db.Where("subscribed_at in ?", bdayIds).Select("user_id").Find(&subscribers).Error
	if err != nil {
		return nil
	}
	var subscribersIds []int
	for _, sub := range subscribers {
		subscribersIds = append(subscribersIds, sub.UserId)
	}

	var users []model.User
	err = d.db.Model(&model.User{}).Where("id in ?", subscribersIds).Select("chat_id").Find(&users).Error
	if err != nil {
		return nil
	}
}

package model

import "time"

type Employee struct {
	Id               int       `gorm:"primaryKey;autoincrement:1"`
	Name             string    `gorm:"size:50;not null"`
	Surname          string    `gorm:"size:50;not null"`
	Birth            time.Time `gorm:"type:date not null"`
	OrganizationName string    `gorm:"size:50"`
}

type Subscription struct {
	Id           int       `gorm:"primaryKey, autoincrement:1"`
	UserId       int       `gorm:"not null"`
	SubscribedTo int       `gorm:"not null"`
	User         *User     `gorm:"foreignKey:UserId;references:Id"`
	Employee     *Employee `gorm:"foreignKey:SubscribedTo;references:Id"`
}

type User struct {
	Id         int    `gorm:"primaryKey;autoincrement" json:"-"`
	EmployeeId int    `gorm:"foreignKey:Id" json:"-"`
	Username   string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password   string `gorm:"size:100;not null" json:"password"`
	ChatId     *int64 `gorm:"uniqueIndex" json:"chatId"`
}

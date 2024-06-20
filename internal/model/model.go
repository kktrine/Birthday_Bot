package model

import "time"

type Employee struct {
	Id      int        `gorm:"primaryKey;autoincrement:1" json:"-"`
	UserId  *int64     `gorm:"not null" json:"userId"`
	Name    string     `gorm:"size:50;not null" json:"name"`
	Surname string     `gorm:"size:50;not null" json:"surname"`
	Birth   *time.Time `gorm:"type:date not null" json:"birth"`
	User    *User      `gorm:"foreignKey:UserId;references:Id" json:"-"`
}

type Subscription struct {
	Id           int       `gorm:"primaryKey, autoincrement:1"`
	UserId       int       `gorm:"not null"`
	SubscribedTo int       `gorm:"not null"`
	User         *User     `gorm:"foreignKey:UserId;references:Id"`
	Employee     *Employee `gorm:"foreignKey:SubscribedTo;references:Id"`
}

type User struct {
	Id       int    `gorm:"primaryKey;autoincrement" json:"-"`
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:100;not null" json:"password"`
	ChatId   *int64 `gorm:"uniqueIndex" json:"chatId"`
}

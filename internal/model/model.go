package model

import "time"

type Employee struct {
	Id      int        `gorm:"primaryKey;autoincrement:1" json:"-"`
	UserId  int        `gorm:"not null;uniqueIndex" json:"userId"`
	Name    string     `gorm:"size:50;not null" json:"name"`
	Surname string     `gorm:"size:50;not null" json:"surname"`
	Birth   *time.Time `gorm:"type:date not null" json:"birth"`
	User    *User      `gorm:"foreignKey:UserId;references:Id" json:"-"`
}

type Subscription struct {
	Id           int       `gorm:"primaryKey, autoincrement:1"`
	UserId       int       `gorm:"not null;uniqueIndex:idx_user_subscribed"`
	SubscribedTo int       `gorm:"not null;uniqueIndex:idx_user_subscribed"`
	User         *User     `gorm:"foreignKey:UserId;references:Id"`
	Employee     *Employee `gorm:"foreignKey:SubscribedTo;references:UserId"`
}

type User struct {
	Id       int    `gorm:"primaryKey;autoincrement" json:"-"`
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:100;not null" json:"password"`
	ChatId   *int64 `gorm:"uniqueIndex" json:"chatId"`
}

type Subscribe struct {
	Id          *int   `json:"id"`
	SubscribeTo *[]int `json:"subscribeTo"`
}

type Messages struct {
	ChatId     int64
	BdayPeople []string
}

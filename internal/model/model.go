package model

import "time"

type Employee struct {
	Id               int       `gorm:"primaryKey, autoincrement:1"`
	Name             string    `gorm:"size:50, not null"`
	Surname          string    `gorm:"size:50, not null"`
	Birth            time.Time `gorm:"type:date, not null"`
	OrganizationName string    `gorm:"size:50"`
}

type Subscription struct {
	Id           int `gorm:"primaryKey, autoincrement:1"`
	EmployeeId   int `gorm:"foreignKey:EmployeeId"`
	SubscribedTo int `gorm:"foreignKey:EmployeeId"`
}

type User struct {
	Id         int    `gorm:"primaryKey, autoincrement" json:"-"`
	EmployeeId int    `gorm:"foreignKey:Id" json:"-"`
	Username   string `gorm:"size:50, index:unique_index, not null" json:"username"`
	Password   string `gorm:"size:100, not null" json:"password"`
}

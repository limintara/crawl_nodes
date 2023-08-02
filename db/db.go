package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func init() {
	cmd := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(cmd), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	DB.AutoMigrate(&Nodes{})
}

type Nodes struct {
	Id           int
	Url          string
	Link         string `gorm:"unique;size:600"`
	Ping         int
	AvgSpeed     int
	MaxSpeed     int
	FailCount    int
	SuccessCount int
	UpdateTime   time.Time
	CreateTime   time.Time
}

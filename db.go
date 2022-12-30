package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type TableStrings struct {
	gorm.Model
	Name          string `gorm:"type:varchar(100); not null" json:"name" form:"name" binding:"required"`
	TranslateStrs string `gorm:"type:longtext; not null" json:"strs" form:"strs" binding:"required"` // json结构数组，存所有的多语言文案
}

func ConnectDB() *gorm.DB {
	dsn := "root:Aa123456@tcp(127.0.0.1:3306)/stringsDB?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	fmt.Println(db)
	fmt.Println(err)

	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(30 * time.Second)

	db.AutoMigrate(&TableStrings{})
	return db
}

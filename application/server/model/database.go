package model

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	// PostgreSQL连接配置
	dsn := "host=127.0.0.1 user=postgres password=123456 dbname=blockchain port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		return fmt.Errorf("连接数据库失败：%v", err)
	}

	// 自动迁移表结构
	err = DB.AutoMigrate(&User{}, &Token{}, &Message{}, &ChatSession{}, &MarketListing{}, &MarketOffer{}, &Lot{}, &Bid{}, &AuctionResult{})
	if err != nil {
		return fmt.Errorf("数据库迁移失败：%v", err)
	}

	log.Println("数据库连接成功，表结构已迁移")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

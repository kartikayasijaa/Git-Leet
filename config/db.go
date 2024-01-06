package config

import (
	"context"
	"fmt"
	"gitleet/structs"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection(config *DBConfig) (*gorm.DB, context.Context, error) {
	//create context
	ctx, _ := context.WithCancel(context.Background())

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, ctx, err
	}

	db = db.WithContext(ctx)
	return db, ctx, nil
}

func InitDB() (*gorm.DB, context.Context) {
	db, ctx, err := NewConnection(&DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.ExpandEnv("DB_SSL"),
		DBName:   os.Getenv("DB_NAME"),
	})

	if err != nil {
		InitLogger().Panic("Can not load DB")
	} else {
		fmt.Println("Database Connected")
	}

	if err = structs.MigrateUser(db); err != nil {
		InitLogger().Panic("User not Migrated")
	}

	return db, ctx
}

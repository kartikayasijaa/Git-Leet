package services

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type Services struct {
	DBService       *DBService
	LeetcodeService *LeetcodeService
}

func ServiceHandler(DBService *DBService, LeetcodeService *LeetcodeService) *Services {
	return &Services{
		DBService:       DBService,
		LeetcodeService: LeetcodeService,
	}
}

func InitServices(db *gorm.DB, logger *log.Logger, ctx context.Context) *Services {
	dbService := DBServicesHandler(db, ctx)
	leetcodeService := &LeetcodeService{}
	services := ServiceHandler(dbService, leetcodeService)

	return services
}

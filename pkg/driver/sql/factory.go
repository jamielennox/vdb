package sql

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
)

type sqlFactory struct {
	db *gorm.DB
}

func (s *sqlFactory) Build(ctx context.Context, name common.CollectionName, value common.DriverData) (driver.Driver, error) {
	return &sqlDriver{
		collectionName: name,
		db:             s.db,
	}, nil
}

func NewSqlDriverFactory(db *gorm.DB) (driver.Factory, error) {
	if err := db.AutoMigrate(&DbData{}); err != nil {
		return nil, fmt.Errorf("failed to migrate dbdata: %w", err)
	}

	return &sqlFactory{
		db: db,
	}, nil
}

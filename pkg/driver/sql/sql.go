package sql

import (
	"context"
	"fmt"
	"vdb/pkg/common"

	"gorm.io/gorm"

	driver "vdb/pkg/driver/base"
)

type sqlDriver struct {
	collectionName common.CollectionName
	db             *gorm.DB
}

type DbData struct {
	gorm.Model

	CollectionName string `gorm:"primaryKey"`
	Id             string `gorm:"primaryKey"`
	Revision       uint   `gorm:"primaryKey;autoIncrement:true"`
	Value          any    `gorm:"serializer:json"`
}

func (s *sqlDriver) GetLatest(ctx context.Context, id common.CollectionId) (driver.Revision, error) {
	ormRet := DbData{}
	result := s.db.WithContext(ctx).Order("revision desc").First(&ormRet, DbData{Id: string(id)})

	if result.Error != nil {
		return driver.Revision{}, result.Error
	}

	return driver.Revision{
		Meta: driver.Meta{
			Id:       common.CollectionId(ormRet.Id),
			Revision: common.RevisionID(ormRet.Revision),
			Version:  1,
		},
		Value: ormRet.Value,
	}, nil
}

func (s *sqlDriver) GetRevisions(ctx context.Context, id common.CollectionId) ([]driver.Revision, error) {
	ormRet := []DbData{}
	result := s.db.WithContext(ctx).Order("revision asc").Find(&ormRet, DbData{Id: string(id)})

	if result.Error != nil {
		return nil, result.Error
	}

	ret := make([]driver.Revision, result.RowsAffected)
	for i, r := range ormRet {
		ret[i] = driver.Revision{
			Meta: driver.Meta{
				Id:       id,
				Revision: common.RevisionID(r.Revision),
				Version:  1,
			},
			Value: r.Value,
		}
	}

	return ret, nil
}

func (s *sqlDriver) Set(ctx context.Context, id common.CollectionId, value common.CollectionValue) (driver.Revision, error) {
	d := DbData{
		CollectionName: string(s.collectionName),
		Id:             string(id),
		Value:          value,
	}

	result := s.db.WithContext(ctx).Create(&d)

	if result.Error != nil {
		return driver.Revision{}, fmt.Errorf("Failed to write record to db: %w", result.Error)
	}

	fmt.Println("created primary key: ", d.Id, d.Revision)

	return driver.Revision{
		Meta: driver.Meta{
			Id:       id,
			Revision: common.RevisionID(d.Revision),
			Version:  1,
		},
		Value: value,
	}, nil
}

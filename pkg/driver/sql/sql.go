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

type Label struct {
	gorm.Model

	CollectionName string `gorm:"primaryKey"`
	Id             string `gorm:"primaryKey"`
	Revision       uint   `gorm:"primaryKey"`
	Key            string `gorm:"primaryKey"`

	Value string
}

type DbData struct {
	gorm.Model

	CollectionName string  `gorm:"primaryKey"`
	Id             string  `gorm:"primaryKey"`
	Labels         []Label `gorm:"foreignKey:CollectionName,Id,Revision"`
	Revision       uint    `gorm:"primaryKey;autoIncrement:true"`

	TransactionId *string

	Value any `gorm:"serializer:json"`
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
	result := s.db.WithContext(ctx).Preload("Labels").Order("revision asc").Find(&ormRet, DbData{Id: string(id)})

	if result.Error != nil {
		return nil, result.Error
	}

	ret := make([]driver.Revision, result.RowsAffected)
	for i, r := range ormRet {
		labels := map[string]string{"hello": "world"}
		for _, l := range r.Labels {
			labels[l.Key] = l.Value
		}

		ret[i] = driver.Revision{
			Meta: driver.Meta{
				Id:       id,
				Revision: common.RevisionID(r.Revision),
				Version:  1,
			},
			Labels: labels,
			Value:  r.Value,
		}
	}

	return ret, nil
}

func (s *sqlDriver) Set(ctx context.Context, transId common.TransactionId, data ...driver.CollectionData) (driver.Transaction, error) {
	var dbd = make([]DbData, len(data))

	// FIXME: There's a bug here that the revisionid is always the max value.

	for i, d := range data {
		dbd[i].CollectionName = string(s.collectionName)
		dbd[i].Id = string(d.Id)
		dbd[i].Value = d.Value
		dbd[i].TransactionId = transId
	}

	result := s.db.WithContext(ctx).Create(&dbd)

	if result.Error != nil {
		return driver.Transaction{}, fmt.Errorf("failed to write record to db: %w", result.Error)
	}

	trans := driver.Transaction{
		Id:        transId,
		Revisions: make([]driver.Revision, len(data)),
	}

	for i, d := range dbd {
		if dbd[i].Labels != nil {
			trans.Revisions[i].Labels = make(common.Labels, len(dbd[i].Labels))

			for _, l := range dbd[i].Labels {
				trans.Revisions[i].Labels[l.Key] = l.Value
			}
		}

		trans.Revisions[i].Meta = driver.Meta{
			Id:       common.CollectionId(d.Id),
			Revision: common.RevisionID(d.Revision),
			Version:  1,
		}

		trans.Revisions[i].Value = dbd[i].Value
	}

	return trans, nil
}

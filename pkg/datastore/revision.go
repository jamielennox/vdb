package datastore

import (
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
)

type Meta struct {
	driver.Meta
	Type common.CollectionName
}

type Revision struct {
	Meta  Meta
	Value common.CollectionValue
}

func convertRevision(typ common.CollectionName, revision *driver.Revision) (Revision, error) {
	return Revision{
		Meta: Meta{
			Meta: revision.Meta,
			Type: typ,
		},
		Value: revision.Value,
	}, nil
}

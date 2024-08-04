package datastore

import (
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
)

type Meta struct {
	driver.Meta
	Type common.TypeName
}

type Revision struct {
	Meta  Meta
	Value common.Value
}

func convertRevision(typ common.TypeName, revision *driver.Revision) (Revision, error) {
	return Revision{
		Meta: Meta{
			Meta: revision.Meta,
			Type: typ,
		},
		Value: revision.Value,
	}, nil
}

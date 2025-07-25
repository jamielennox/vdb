package collection

import (
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
)

type Meta struct {
	driver.Meta
	Type common.CollectionName
}

type Revision struct {
	Meta   Meta
	Labels common.Labels
	Value  common.CollectionValue
}

type Transaction struct {
	Id        common.TransactionId
	Revisions []Revision
}

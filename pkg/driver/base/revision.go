package base

import (
	"vdb/pkg/common"
)

const DefaultVersion common.VersionID = 1

type Meta struct {
	Id       common.CollectionId
	Revision common.RevisionID
	Version  common.VersionID
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

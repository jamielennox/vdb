package datastore

import (
	"fmt"

	"vdb/pkg/common"
)

type ErrUnknownType struct {
	Type common.CollectionName
}

func (e ErrUnknownType) Error() string {
	return fmt.Sprintf("unknown type: %s", e.Type)
}

type ErrIdNotFound struct {
	Type common.CollectionName
	Id   common.CollectionId
}

func (e ErrIdNotFound) Error() string {
	return fmt.Sprintf("Not found id: %s for type: %s", e.Id, e.Type)
}

type ErrRevisionNotFound struct {
	Type       common.CollectionName
	Id         common.CollectionId
	RevisionID common.RevisionID
}

func (e ErrRevisionNotFound) Error() string {
	return fmt.Sprintf("Not found revision: %s for id: %s for type: %s", e.RevisionID, e.Id, e.Type)
}

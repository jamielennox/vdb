package datastore

import (
	"fmt"

	"vdb/pkg/common"
)

type ErrUnknownType struct {
	Type common.TypeName
}

func (e ErrUnknownType) Error() string {
	return fmt.Sprintf("unknown type: %s", e.Type)
}

type ErrIdNotFound struct {
	Type common.TypeName
	Id   common.TypeID
}

func (e ErrIdNotFound) Error() string {
	return fmt.Sprintf("Not found id: %s for type: %s", e.Id, e.Type)
}

type ErrRevisionNotFound struct {
	Type       common.TypeName
	Id         common.TypeID
	RevisionID common.RevisionID
}

func (e ErrRevisionNotFound) Error() string {
	return fmt.Sprintf("Not found revision: %s for id: %s for type: %s", e.RevisionID, e.Id, e.Type)
}

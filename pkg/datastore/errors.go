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

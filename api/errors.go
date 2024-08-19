package api

import (
	"vdb/pkg/collection"
	"vdb/pkg/datastore"
)

func RenderErrUnknownType(e datastore.ErrUnknownType) ErrNotFound {
	t := "type"
	v := string(e.Type)

	return ErrNotFound{
		Message: e.Error(),
		Type:    &t,
		Value:   &v,
	}
}

func RenderErrIdNotFound(e collection.ErrIdNotFound) ErrNotFound {
	t := string(e.Type)
	v := string(e.Id)

	return ErrNotFound{
		Message: e.Error(),
		Type:    &t,
		Value:   &v,
	}
}

func RenderRevisionIdNotFound(e collection.ErrRevisionNotFound) ErrNotFound {
	t := string(e.Type)
	v := string(e.Id)
	r := RevisionId(e.RevisionID)

	return ErrNotFound{
		Message:  e.Error(),
		Type:     &t,
		Value:    &v,
		Revision: &r,
	}
}

func RenderServerError(err error) ErrServerError {
	return ErrServerError{
		Message: err.Error(),
	}
}

package http

import (
	"vdb/pkg/collection"
)

func renderRevision(revision collection.Revision) Revision {
	r := Revision{
		Meta: Meta{
			Id:       string(revision.Meta.Id),
			Revision: uint64(revision.Meta.Revision),
			Type:     string(revision.Meta.Type),
			Version:  int64(revision.Meta.Version),
		},
		Value: revision.Value,
	}

	if len(revision.Labels) > 0 {
		mr := Labels(revision.Labels)
		r.Labels = &mr
	}

	return r
}

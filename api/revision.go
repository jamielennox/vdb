package api

import "vdb/pkg/datastore"

func renderRevision(revision datastore.Revision) Revision {
	return Revision{
		Meta: Meta{
			Id:       string(revision.Meta.Id),
			Revision: uint64(revision.Meta.Revision),
			Type:     string(revision.Meta.Type),
			Version:  int64(revision.Meta.Version),
		},
		Value: revision.Value,
	}
}

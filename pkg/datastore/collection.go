package datastore

import (
	"context"
	"fmt"
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
	validator "vdb/pkg/validator/base"
)

type Collection struct {
	Name common.CollectionName
	dri  driver.Driver
	vali validator.Validator
}

func (c *Collection) Get(ctx context.Context, id common.CollectionId) (Revision, error) {
	rev, err := c.dri.GetLatest(ctx, id)
	if err != nil {
		return Revision{}, ErrIdNotFound{Type: c.Name, Id: id}
	}

	if err = c.vali.Validate(ctx, rev.Value); err != nil {
		return Revision{}, err
	}

	return convertRevision(c.Name, &rev)
}

func (c *Collection) GetRevisions(ctx context.Context, id common.CollectionId) ([]Revision, error) {
	revs, err := c.dri.GetRevisions(ctx, id)
	if err != nil {
		return nil, ErrIdNotFound{Type: c.Name, Id: id}
	}

	ret := make([]Revision, len(revs))

	for i, x := range revs {
		if err = c.vali.Validate(ctx, x.Value); err != nil {
			return nil, fmt.Errorf("validate failure in index %d: %w", i, err)
		}

		r, err := convertRevision(c.Name, &x)
		if err != nil {
			return nil, fmt.Errorf("convert revision failure in index %d: %w", i, err)
		}

		ret[i] = r
	}

	return ret, nil
}

func (c *Collection) GetRevision(ctx context.Context, id common.CollectionId, revisionId common.RevisionID) (Revision, error) {
	revs, err := c.dri.GetRevisions(ctx, id)
	if err != nil {
		return Revision{}, ErrIdNotFound{Type: c.Name, Id: id}
	}

	if int(revisionId) >= len(revs) {
		return Revision{}, ErrRevisionNotFound{Type: c.Name, Id: id, RevisionID: revisionId}
	}

	rev := revs[revisionId]

	if err = c.vali.Validate(ctx, rev.Value); err != nil {
		return Revision{}, err
	}

	return convertRevision(c.Name, &rev)
}

func (c *Collection) Set(ctx context.Context, id common.CollectionId, value common.CollectionValue) (Revision, error) {
	if err := c.vali.Validate(ctx, value); err != nil {
		return Revision{}, err
	}

	rev, err := c.dri.Set(ctx, id, value)
	if err != nil {
		return Revision{}, err
	}

	return convertRevision(c.Name, &rev)
}

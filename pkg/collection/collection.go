package collection

import (
	"context"
	"fmt"
	"log/slog"
	audit "vdb/pkg/audit/base"
	authz "vdb/pkg/authz/base"
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
	validator "vdb/pkg/validator/base"
)

type Collection struct {
	Name   common.CollectionName
	Labels common.Labels

	logger *slog.Logger
	aud    audit.Auditor
	dri    driver.Driver
	vali   validator.Validator
	authz  authz.Authorizer
}

type CollectionData struct {
	Id    common.CollectionId
	Value common.CollectionValue
}

func NewCollection(name common.CollectionName, aud audit.Auditor, dri driver.Driver, opts ...Option) (*Collection, error) {
	o := getCollectionOptions(opts...)

	return &Collection{
		Name:   name,
		Labels: o.labels,
		aud:    aud,
		dri:    dri,
		vali:   o.vali,
		authz:  o.authz,
		logger: o.logger.With(slog.String("collection", string(name))),
	}, nil
}

func (c *Collection) Get(ctx context.Context, id common.CollectionId) (Revision, error) {
	user := common.UserInfo{
		UserName: "jamie",
		Roles:    []string{"admin"},
	}

	slog.Debug(
		"get",
		slog.String("operation", string(common.OperationRead)),
		slog.String("id", string(id)),
		slog.String("user", user.UserName),
	)

	rev, err := c.dri.GetLatest(ctx, id)
	if err != nil {
		return Revision{}, ErrIdNotFound{Type: c.Name, Id: id}
	}

	if err = c.vali.Validate(ctx, rev.Value); err != nil {
		return Revision{}, err
	}

	newRev, err := convertRevision(c.Name, &rev)
	if err != nil {
		return Revision{}, err
	}

	event := common.Event{
		Operation: common.OperationRead,
		Target: common.CollectionTarget{
			Name:     c.Name,
			Id:       id,
			Revision: rev.Meta.Revision,
			Labels:   c.Labels,
			Type:     "collection",
		},
		Subject: user,
	}

	c.aud.Event(event)
	return newRev, nil
}

func (c *Collection) GetRevisions(ctx context.Context, id common.CollectionId) ([]Revision, error) {
	revs, err := c.dri.GetRevisions(ctx, id)
	if err != nil {
		return nil, ErrIdNotFound{Type: c.Name, Id: id}
	}

	ret := make([]Revision, 0, len(revs))

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

	newRev, err := convertRevision(c.Name, &rev)

	event := common.Event{
		Operation: common.OperationRead,
		Target: common.CollectionTarget{
			Name:     c.Name,
			Id:       id,
			Revision: rev.Meta.Revision,
			Labels:   c.Labels,
			Type:     "collection",
		},
		Subject: common.UserInfo{
			UserName: "jamie",
			Roles:    []string{"admin"},
		},
	}

	c.aud.Event(event)
	return newRev, nil
}

func (c *Collection) Set(ctx context.Context, transId common.TransactionId, data ...CollectionData) (Transaction, error) {

	driverData := make([]driver.CollectionData, 0, len(data))
	for _, d := range data {
		if err := c.vali.Validate(ctx, d.Value); err != nil {
			return Transaction{}, err
		}

		driverData = append(driverData, driver.CollectionData{
			Id:    d.Id,
			Value: d.Value,
		})
	}

	t, err := c.dri.Set(ctx, nil, driverData...)
	if err != nil {
		return Transaction{Id: transId}, err
	}

	if len(t.Revisions) != 1 {
		return Transaction{Id: transId}, fmt.Errorf("unexpected revision response count, expected (1), found (%d)", len(t.Revisions))
	}

	trans := Transaction{
		Id:        t.Id,
		Revisions: make([]Revision, 0, len(t.Revisions)),
	}
	events := []common.Event{}

	for _, r := range t.Revisions {
		newRev, err := convertRevision(c.Name, &r)

		if err != nil {
			return trans, nil
		}

		trans.Revisions = append(trans.Revisions, newRev)
		events = append(events, common.Event{
			Operation: common.OperationUpdate,
			Target: common.CollectionTarget{
				Name:     c.Name,
				Id:       r.Meta.Id,
				Revision: r.Meta.Revision,
				Labels:   c.Labels,
				Type:     "collection",
			},
			Subject: common.UserInfo{
				UserName: "jamie",
				Roles:    []string{"admin"},
			},
		})
	}

	c.aud.Event(events...)
	return trans, nil
}

func convertRevision(typ common.CollectionName, revision *driver.Revision) (Revision, error) {
	return Revision{
		Meta: Meta{
			Meta: revision.Meta,
			Type: typ,
		},
		Labels: revision.Labels,
		Value:  revision.Value,
	}, nil
}

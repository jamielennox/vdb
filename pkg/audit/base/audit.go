package base

import "vdb/pkg/common"

type Auditor interface {
	Event(event ...common.Event)
}

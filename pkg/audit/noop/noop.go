package noop

import (
	"fmt"
	"vdb/pkg/audit/base"
	"vdb/pkg/common"
)

type noopAuditor struct {
	prefix string
}

func (n *noopAuditor) Event(event common.Event) {
	fmt.Printf(
		"%s%s %s %s - %s/%s/%d\n",
		n.prefix,
		event.Subject.UserName,
		event.Operation,
		event.Target.Type,
		event.Target.Name,
		event.Target.Id,
		event.Target.Revision,
	)
}

func NewNoopAuditor() base.Auditor {
	return &noopAuditor{
		prefix: "audit: ",
	}
}

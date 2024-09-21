package factory

import (
	"vdb/pkg/common"
	driver "vdb/pkg/driver/base"
)

type Output[K ~string, V, O any] struct {
	Object[K, V, O]

	driver     driver.Driver
	collection common.CollectionName
	Meta       driver.Meta
}

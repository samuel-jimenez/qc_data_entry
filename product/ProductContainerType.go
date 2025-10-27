package product

import "database/sql"

// bs.container_types
type ProductContainerType int

const (
	CONTAINER_SAMPLE ProductContainerType = 1 + iota
	CONTAINER_TOTE
	CONTAINER_RAILCAR
	CONTAINER_ISO
)

// DONTTODO go:generate stringer -type=ProductContainer Thanks I hate it
var ProductContainerTypes = []string{"Sample", "Tote", "Railcar", "ISO"}

func (container ProductContainerType) String() string {
	return ProductContainerTypes[container-1]
}

func DiscreteFromContainer(Container ProductContainerType) Discrete {
	return Discrete{sql.NullInt32{Int32: int32(Container), Valid: Container != 0}}
}

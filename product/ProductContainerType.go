package product

// bs.container_types
type ProductContainerType int

const (
	CONTAINER_SAMPLE ProductContainerType = iota + 1
	CONTAINER_TOTE
	CONTAINER_RAILCAR
	CONTAINER_ISO
)

// DONTTODO go:generate stringer -type=ProductContainer Thanks I hate it
var ProductContainerTypes = []string{"Sample", "Tote", "Railcar", "ISO"}

func (container ProductContainerType) String() string {
	return ProductContainerTypes[container-1]
}

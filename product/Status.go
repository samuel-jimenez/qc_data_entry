package product

// bs.internal_status_list
type Status string

const (
	Status_REQUESTED Status = "REQUESTED"
	Status_PRINTED          = "PRINTED"
	Status_BLENDED          = "BLENDED"
	Status_TESTED           = "TESTED"
	Status_SHIPPED          = "SHIPPED"
)

// // TODO? type01 create type
// // bs.internal_status_list
// type Status int
//
// const (
// 	Status_REQUESTED Status = iota
// 	Status_PRINTED
// 	Status_BLENDED
// 	Status_TESTED
// 	Status_SHIPPED
// )
//
// var (
// 	INTERNAL_STATUS_LIST = []string{"REQUESTED", "PRINTED", "BLENDED", "TESTED", "SHIPPED"}
// )
//
// func (status Status) String() string {
// 	return INTERNAL_STATUS_LIST[status]
// }

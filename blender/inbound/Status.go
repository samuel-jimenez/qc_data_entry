package inbound

// TODO type00 export type
// bs.inbound_status_list
type Status string

const (
	Status_AVAILABLE   Status = "AVAILABLE"
	Status_SAMPLED            = "SAMPLED"
	Status_TESTED             = "TESTED"
	Status_UNAVAILABLE        = "UNAVAILABLE"
)

// TODO? type00 create type
// bs.inbound_status_list
// type Status int
//
// const (
// 	Status_AVAILABLE Status = iota
// 	Status_SAMPLED
// 	Status_TESTED
// 	Status_UNAVAILABLE
// )
//
// var (
// 	INTERNAL_STATUS_LIST = []string{"AVAILABLE", "SAMPLED", "TESTED", "UNAVAILABLE"}
// )
//
// func (status Status) String() string {
// 	return INTERNAL_STATUS_LIST[status]
// }

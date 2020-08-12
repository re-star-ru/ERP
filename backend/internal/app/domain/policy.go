package domain

type Action int

// Action attributes
const (
	ActionProductCreate Action = iota
	ActionProductDelete
)

//func CheckPermission() bool {
//	return true
//}

//
//type RBAC struct {
//}
//
//type RBACusecase interface {
//}
//
//type Role struct {
//	Name string
//	Permissions []
//}
//
//type RoleUsecase interface {
//}
//
//type Permission struct {
//
//}

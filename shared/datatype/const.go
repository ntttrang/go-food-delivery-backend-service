package datatype

type UserRole string

const (
	KeyRequester          = "requester"
	RoleUser     UserRole = "USER"
	RoleAdmin    UserRole = "ADMIN"
	RoleShipper  UserRole = "SHIPPER"
)

package datatype

type UserRole string

const (
	KeyRequester          = "requester"
	RoleUser     UserRole = "USER"
	RoleAdmin    UserRole = "ADMIN"
	RoleShipper  UserRole = "SHIPPER"

	// Cart Status
	CartStatusActive    = "ACTIVE"
	CartStatusUpdated   = "UPDATED"   // When Frontend update quantity
	CartStatusProcessed = "PROCESSED" // Auto updated by Backend. All items go to Order

)

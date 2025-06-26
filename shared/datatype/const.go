package datatype

type UserRole string

const (
	RoleUser    UserRole = "USER"
	RoleAdmin   UserRole = "ADMIN"
	RoleShipper UserRole = "SHIPPER"
)

type RecordStatus string

const (
	RecordStatusActive   RecordStatus = "ACTIVE"
	RecordStatusInactive RecordStatus = "INACTIVE"
)

type UserStatus string

const (
	StatusPending  UserStatus = "PENDING"
	StatusActive   UserStatus = "ACTIVE"
	StatusInactive UserStatus = "INACTIVE"
	StatusBanned   UserStatus = "BANNED"
	StatusDeleted  UserStatus = "DELETED"
)

type UserType string

const (
	TypeEmailPassword UserType = "EMAIL_PASSWORD"
	TypeFacebook      UserType = "FACEBOOK"
	TypeGmail         UserType = "GMAIL"
)

const (
	KeyRequester = "requester"
)

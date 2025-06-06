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

const (
	EvtNotifyOrderCreate         = "order-create"
	EvtNotifyOrderStateChange    = "order-state-change"
	EvtNotifyOrderCancel         = "order-cancel"
	EvtNotifyShipperAssign       = "order-shipper-assign"
	EvtNotifyPaymentStatusChange = "order-payment-status-change"
)

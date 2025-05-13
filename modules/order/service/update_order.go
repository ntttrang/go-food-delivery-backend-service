package service

import (
	"context"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type OrderUpdateReq struct {
	ID         string  `json:"-"`
	ShipperID  *string `json:"shipperId,omitempty"`
	State      *string `json:"state,omitempty"`
	PaymentStatus *string `json:"paymentStatus,omitempty"`
}

func (o *OrderUpdateReq) Validate() error {
	if o.ID == "" {
		return ordermodel.ErrOrderIdRequired
	}

	return nil
}

// Initialize service
type IUpdateOrderRepo interface {
	FindById(ctx context.Context, id string) (*ordermodel.Order, *ordermodel.OrderTracking, []ordermodel.OrderDetail, error)
	Update(ctx context.Context, order *ordermodel.Order, tracking *ordermodel.OrderTracking) error
}

type UpdateCommandHandler struct {
	repo IUpdateOrderRepo
}

func NewUpdateCommandHandler(repo IUpdateOrderRepo) *UpdateCommandHandler {
	return &UpdateCommandHandler{repo: repo}
}

// Implement
func (hdl *UpdateCommandHandler) Execute(ctx context.Context, req OrderUpdateReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	// Get order from database
	order, tracking, _, err := hdl.repo.FindById(ctx, req.ID)
	if err != nil {
		if err == ordermodel.ErrOrderNotFound {
			return datatype.ErrNotFound.WithWrap(err).WithDebug(err.Error())
		}
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Update order fields
	if req.ShipperID != nil {
		order.ShipperID = req.ShipperID
	}

	// Update tracking fields
	if req.State != nil {
		// Validate state transitions
		if tracking.State == "delivered" && *req.State != "delivered" {
			return datatype.ErrBadRequest.WithWrap(ordermodel.ErrOrderIsDelivered).WithDebug(ordermodel.ErrOrderIsDelivered.Error())
		}
		if tracking.State == "cancel" && *req.State != "cancel" {
			return datatype.ErrBadRequest.WithWrap(ordermodel.ErrOrderIsCancelled).WithDebug(ordermodel.ErrOrderIsCancelled.Error())
		}
		tracking.State = *req.State
	}

	if req.PaymentStatus != nil {
		tracking.PaymentStatus = *req.PaymentStatus
	}

	// Update in database
	if err := hdl.repo.Update(ctx, order, tracking); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}

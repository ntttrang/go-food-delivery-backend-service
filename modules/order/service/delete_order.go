package service

import (
	"context"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type OrderDeleteReq struct {
	ID string `json:"-"`
}

func (o *OrderDeleteReq) Validate() error {
	if o.ID == "" {
		return ordermodel.ErrOrderIdRequired
	}

	return nil
}

// Initialize service
type IDeleteOrderRepo interface {
	FindById(ctx context.Context, id string) (*ordermodel.Order, *ordermodel.OrderTracking, []ordermodel.OrderDetail, error)
	Delete(ctx context.Context, id string) error
}

type DeleteCommandHandler struct {
	repo IDeleteOrderRepo
}

func NewDeleteCommandHandler(repo IDeleteOrderRepo) *DeleteCommandHandler {
	return &DeleteCommandHandler{repo: repo}
}

// Implement
func (hdl *DeleteCommandHandler) Execute(ctx context.Context, req OrderDeleteReq) error {
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

	// Check if order can be deleted
	if tracking.State != "waiting_for_shipper" && tracking.State != "preparing" {
		return datatype.ErrBadRequest.WithWrap(ordermodel.ErrOrderIsProcessed).WithDebug(ordermodel.ErrOrderIsProcessed.Error())
	}

	// Soft delete by updating status
	order.Status = string(datatype.StatusDeleted)
	tracking.Status = string(datatype.StatusDeleted)

	// Update in database
	if err := hdl.repo.Delete(ctx, req.ID); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}

package repository

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/pkg/errors"
)

type UpdateFoodDto struct {
	Name         string  `json:"name"` // Update when it has value
	Description  string  `json:"description"`
	Status       string  `json:"status"`
	RestaurantId *string `json:"restaurantId"` // Pointer => Update value => empty
	CategoryId   *string `json:"categoryId"`
}

func (r *FoodRepo) Update(ctx context.Context, id uuid.UUID, req foodmodel.FoodUpdateReq) error {
	db := r.dbCtx.GetMainConnection().Begin()

	var status = ""
	if req.Status != nil {
		status = *req.Status
	}
	updateDto := UpdateFoodDto{
		Name:         req.Name,
		Description:  req.Description,
		Status:       status,
		RestaurantId: req.RestaurantId,
		CategoryId:   req.CategoryId,
	}
	if err := db.Table(req.TableName()).Where("id = ?", id).Updates(updateDto).Error; err != nil {
		db.Rollback()
		return errors.WithStack(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return errors.WithStack(err)
	}

	return nil
}

package foodgormmysql

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type FoodInfoDto struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Images      string    `json:"images"`
	Price       float64   `json:"price"`
	AvgPoint    float64   `json:"avgPoint"`
	CommentQty  int       `json:"commentQty"`
	CategoryId  uuid.UUID `json:"categoryId"`
	Status      string    `json:"status"`
}

func (r *FoodRepo) FindByIds(ctx context.Context, ids []uuid.UUID) ([]foodmodel.FoodInfoDto, error) {
	var foods []foodmodel.FoodInfoDto
	db := r.dbCtx.GetMainConnection()
	if err := db.Raw(`SELECT f.id, 
				f.name, 
				f.description, 
				f.images, 
				f.price, 
				f.category_id,
				f.status,
				COUNT(fr.comment) AS comment_qty,
				AVG(fr.point) AS avg_point,
				f.restaurant_id
				FROM foods f
			LEFT JOIN food_ratings fr
			ON f.id = fr.food_id
			WHERE f.id IN (?) AND f.status = ?
			GROUP BY f.id, f.name, f.description, f.images, f.price, f.category_id, f.restaurant_id`, ids, sharedModel.StatusActive).
		Find(&foods).Error; err != nil {
		return nil, err
	}
	return foods, nil
}

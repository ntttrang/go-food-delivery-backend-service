package repository

import (
	"context"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"github.com/pkg/errors"
)

// Shouldn't use Goroutines here. It makes your code more complex and dificult to maintain
// func (r *CategoryRepo) ListCategories(ctx context.Context, req categorymodel.ListCategoryReq) ([]categorymodel.Category, int64, error) {

// 	var categories []categorymodel.Category
// 	var total int64

// 	query := r.db.Model(&categorymodel.Category{})

// 	if req.Name != "" {
// 		query = query.Where("name LIKE ?", "%"+req.Name+"%")
// 	}

// 	if req.Description != "" {
// 		query = query.Where("description LIKE ?", "%"+req.Description+"%")
// 	}
// 	query = query.Where("status in (?)", []string{sharedModel.StatusActive})

// 	totalChan := make(chan int64)
// 	categoriesChan := make(chan []categorymodel.Category)
// 	rsChan := make(chan error)
// 	go func(rsChan chan error) {
// 		var tempTotal int64
// 		rs := query.Count(&tempTotal)
// 		if rs.Error != nil {
// 			rsChan <- errors.New("COUNT_ERR: " + rs.Error.Error())
// 		}
// 		totalChan <- tempTotal
// 	}(rsChan)

// 	go func(rsChan chan error) {
// 		var tempCategories []categorymodel.Category
// 		rs := query.Limit(req.Paging.Limit).Offset(req.Paging.Page).Find(&tempCategories)
// 		if rs.Error != nil {
// 			rsChan <- errors.New("QUERY_ERR: " + rs.Error.Error())
// 		}
// 		categoriesChan <- tempCategories

// 	}(rsChan)

// 	i := 0
// 	for {
// 		if i > 2 {
// 			break
// 		}
// 		select {
// 		case total = <-totalChan:
// 			i++
// 		case categories = <-categoriesChan:
// 			i++
// 		case err := <-rsChan:
// 			return nil, 0, err
// 		}
// 	}

// 	return categories, total, nil
// }

func (r *CategoryRepo) ListCategories(ctx context.Context, req categorymodel.ListCategoryReq) ([]categorymodel.Category, int64, error) {

	var categories []categorymodel.Category
	var total int64

	db := r.dbCtx.GetMainConnection()

	query := db.Model(&categorymodel.Category{})

	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.Description != "" {
		query = query.Where("description LIKE ?", "%"+req.Description+"%")
	}
	query = query.Where("status in (?)", []string{sharedModel.StatusActive})

	rs := query.Count(&total).Limit(req.Paging.Limit).Offset(req.Paging.Page).Find(&categories)
	if rs.Error != nil {
		return nil, 0, errors.WithStack(rs.Error)
	}

	return categories, total, nil
}

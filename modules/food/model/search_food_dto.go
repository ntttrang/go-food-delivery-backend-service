package foodmodel

type SearchFoodDto struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}

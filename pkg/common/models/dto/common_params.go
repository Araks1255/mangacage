package dto

type CommonParams struct {
	Sort  string  `form:"sort"`
	Query *string `form:"query"`
	Order string  `form:"order"`
	Page  int     `form:"page,default=1"`
	Limit uint    `form:"limit,default=60" binding:"max=1000"`
}

package order

type Order struct {
	Id       int
	UserId   int
	Products []OrderProduct
}

type OrderProduct struct {
	ProductId     int `json:"product_id" binding:"required"`
	ProductAmount int `json:"product_amount" binding:"required"`
}

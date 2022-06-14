package order

type OrderResponse struct {
	Id         int                     `json:"id"`
	User       OrderUserResponse       `json:"user"`
	Products   []OrderProductsResponse `json:"products"`
	OrderPrice int                     `json:"order_price"`
}

type OrderUserResponse struct {
	UserId   *int    `json:"user_id"`
	UserName *string `json:"user_name"`
}

type OrderProductsResponse struct {
	ProductId     *int    `json:"product_id"`
	ProductName   *string `json:"product_name"`
	ProductPrice  *int    `json:"product_price"`
	ProductAmount *int    `json:"product_amount"`
}

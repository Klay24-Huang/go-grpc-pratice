package order

type OrderRequest struct {
	UserId       int            `json:"user_id" binding:"required"`
	OrderProduct []OrderProduct `json:"products" binding:"required"`
}

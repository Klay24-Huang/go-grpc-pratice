package user

type UserRequest struct {
	Name    string `json:"name" binding:"required"`
	Account string `json:"account" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
}

package user

type UserResponse struct {
	Id      int    `json:"id"`
	Account string `json:"account"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
}

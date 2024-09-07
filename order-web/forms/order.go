package forms

type CreateOrderForm struct {
	Address string `json:"address" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required,Phone"`
	Post    string `json:"post" binding:"required"`
}

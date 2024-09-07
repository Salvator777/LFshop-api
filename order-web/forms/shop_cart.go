package forms

type ShopCartItemForm struct {
	GoodsId int32 `json:"goods" binding:"required"`
	Nums    int32 `json:"nums" binding:"required,min=1"`
}

type ShopCartItemUpdateForm struct {
	Nums int32 `json:"nums" binding:"required,min=1"`
	// check不是必填的，所以设置为指针类型
	// checked 字段通常用于表示该商品是否被选中，尤其是在用户进行批量操作时
	Checked *bool `json:"checked"`
}

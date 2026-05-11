package dto

type Cart struct {
	Key   string `json:"Key"`
	Items []Item `json:"Items"`
}

type Item struct {
	Item     uint64 `json:"Item"`
	Quantity uint   `json:"Quantity"`
}

type UpdateItemQuantity struct {
	Item     uint64 `json:"Item" binding:"required"`
	Quantity *uint  `json:"Quantity" binding:"required"`
}

type CartResponse struct {
	Message      string              `json:"Message"`
	Cart         Cart                `json:"Cart"`
	RemovedItems []Item              `json:"RemovedItems"`
	Bill         *CreateOrderBillDTO `json:"Bill"`
}

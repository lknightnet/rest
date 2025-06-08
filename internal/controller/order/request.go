package order

type OrderRequest struct {
	InstrumentationQuantity int    `json:"instrumentation_quantity"`
	IsDelivery              bool   `json:"is_delivery"`
	PaymentMethod           bool   `json:"payment_method"`
	City                    string `json:"city"`
	Bonuses                 int    `json:"bonuses"`
	Comment                 string `json:"comment"`
}

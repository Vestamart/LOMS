package domain

type OrderStatus int

const (
	New OrderStatus = iota
	AwaitingPayment
	Failed
	Payed
	Cancelled
)

type Order struct {
	UserID int64
	Status OrderStatus
	Items  []Item
}

type Item struct {
	Sku   uint32 `json:"sku"`
	Count uint32 `json:"count"`
}

type StocksItem struct {
	TotalCount uint32 `json:"total_count"`
	Reserved   uint32 `json:"reserved"`
}

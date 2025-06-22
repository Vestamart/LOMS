package domain

import (
	"time"
)

type OrderStatus int

const (
	New OrderStatus = iota
	AwaitingPayment
	Failed
	Payed
	Cancelled
)

type EventType string

const (
	OrderCreated         EventType = "new"
	OrderFailed          EventType = "failed"
	OrderAwaitingPayment EventType = "awaiting_payment"
	OrderPayed           EventType = "payed"
	OrderCancelled       EventType = "cancelled"
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

type OrderEvent struct {
	OrderID   int64     `json:"order_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Info      string    `json:"info"`
}

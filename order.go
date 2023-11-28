package rbt_orderbook

import "github.com/shopspring/decimal"

// Single Order in an order book, as a node in a LimitOrder FIFO queue
type Order struct {
	Id       int
	Volume   decimal.Decimal
	Next     *Order
	Prev     *Order
	Limit    *LimitOrder
	BidOrAsk bool
}

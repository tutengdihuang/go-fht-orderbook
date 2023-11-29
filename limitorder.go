package rbt_orderbook

import (
	"github.com/shopspring/decimal"
)

// Limit price orders combined as a FIFO queue
type LimitOrder struct {
	Price decimal.Decimal

	orders      *ordersQueue
	totalVolume decimal.Decimal
}

func NewLimitOrder(price decimal.Decimal) LimitOrder {
	q := NewOrdersQueue()
	return LimitOrder{
		Price:       price,
		orders:      &q,
		totalVolume: decimal.NewFromFloat(0.0),
	}
}

func (this *LimitOrder) TotalVolume() decimal.Decimal {
	return this.totalVolume
}

func (this *LimitOrder) Size() int {
	return this.orders.Size()
}

func (this *LimitOrder) Enqueue(o *Order) {
	this.orders.Enqueue(o)
	o.Limit = this
	this.totalVolume = this.totalVolume.Add(o.Volume)
}

func (this *LimitOrder) Dequeue() *Order {
	if this.orders.IsEmpty() {
		return nil
	}

	o := this.orders.Dequeue()
	this.totalVolume = this.totalVolume.Sub(o.Volume)
	return o
}

func (this *LimitOrder) Delete(o *Order) {
	if o.Limit != this {
		panic("order does not belong to the limit")
	}

	this.orders.Delete(o)
	o.Limit = nil
	this.totalVolume = this.totalVolume.Sub(o.Volume)
}

func (this *LimitOrder) Clear() {
	q := NewOrdersQueue()
	this.orders = &q
	this.totalVolume = decimal.Zero
}

package rbt_orderbook

import (
	"fmt"
	"github.com/shopspring/decimal"
	"sync"
)

// maximum limits per orderbook side to pre-allocate memory
const MaxLimitsNum int = 10000

type Orderbook struct {
	Bids *redBlackBST
	Asks *redBlackBST

	bidLimitsCache map[decimal.Decimal]*LimitOrder
	askLimitsCache map[decimal.Decimal]*LimitOrder
	pool           *sync.Pool
}

func NewOrderbook() Orderbook {
	bids := NewRedBlackBST()
	asks := NewRedBlackBST()
	return Orderbook{
		Bids: &bids,
		Asks: &asks,

		bidLimitsCache: make(map[decimal.Decimal]*LimitOrder, MaxLimitsNum),
		askLimitsCache: make(map[decimal.Decimal]*LimitOrder, MaxLimitsNum),
		pool: &sync.Pool{
			New: func() interface{} {
				limit := NewLimitOrder(decimal.NewFromFloat(0.0))
				return &limit
			},
		},
	}
}

func (this *Orderbook) Add(price decimal.Decimal, o *Order) {
	var limit *LimitOrder

	if o.BidOrAsk {
		limit = this.bidLimitsCache[price]
	} else {
		limit = this.askLimitsCache[price]
	}

	if limit == nil {
		// getting a new limit from pool
		limit = this.pool.Get().(*LimitOrder)
		limit.Price = price

		// insert into the corresponding BST and cache
		if o.BidOrAsk {
			this.Bids.Put(price, limit)
			this.bidLimitsCache[price] = limit
		} else {
			this.Asks.Put(price, limit)
			this.askLimitsCache[price] = limit
		}
	}

	// add order to the limit
	limit.Enqueue(o)
}

func (this *Orderbook) Cancel(o *Order) {
	limit := o.Limit
	limit.Delete(o)

	if limit.Size() == 0 {
		// remove the limit if there are no orders
		if o.BidOrAsk {
			this.Bids.Delete(limit.Price)
			delete(this.bidLimitsCache, limit.Price)
		} else {
			this.Asks.Delete(limit.Price)
			delete(this.askLimitsCache, limit.Price)
		}

		// put it back to the pool
		this.pool.Put(limit)
	}
}

func (this *Orderbook) ClearBidLimit(price decimal.Decimal) {
	this.clearLimit(price, true)
}

func (this *Orderbook) ClearAskLimit(price decimal.Decimal) {
	this.clearLimit(price, false)
}

func (this *Orderbook) clearLimit(price decimal.Decimal, bidOrAsk bool) {
	var limit *LimitOrder
	if bidOrAsk {
		limit = this.bidLimitsCache[price]
	} else {
		limit = this.askLimitsCache[price]
	}

	if limit == nil {
		panic(fmt.Sprintf("there is no such price limit %+v", price))
	}

	limit.Clear()
}

func (this *Orderbook) DeleteBidLimit(price decimal.Decimal) {
	limit := this.bidLimitsCache[price]
	if limit == nil {
		return
	}

	this.deleteLimit(price, true)
	delete(this.bidLimitsCache, price)

	// put limit back to the pool
	limit.Clear()
	this.pool.Put(limit)

}

func (this *Orderbook) DeleteAskLimit(price decimal.Decimal) {
	limit := this.askLimitsCache[price]
	if limit == nil {
		return
	}

	this.deleteLimit(price, false)
	delete(this.askLimitsCache, price)

	// put limit back to the pool
	limit.Clear()
	this.pool.Put(limit)
}

func (this *Orderbook) deleteLimit(price decimal.Decimal, bidOrAsk bool) {
	if bidOrAsk {
		this.Bids.Delete(price)
	} else {
		this.Asks.Delete(price)
	}
}

func (this *Orderbook) GetVolumeAtBidLimit(price decimal.Decimal) decimal.Decimal {
	limit := this.bidLimitsCache[price]
	if limit == nil {
		return decimal.Zero
	}
	return limit.TotalVolume()
}

func (this *Orderbook) GetVolumeAtAskLimit(price decimal.Decimal) decimal.Decimal {
	limit := this.askLimitsCache[price]
	if limit == nil {
		return decimal.Zero
	}
	return limit.TotalVolume()
}

func (this *Orderbook) GetBestBid() decimal.Decimal {
	return this.Bids.Max()
}

func (this *Orderbook) GetBestOffer() decimal.Decimal {
	return this.Asks.Min()
}

func (this *Orderbook) BLength() int {
	return len(this.bidLimitsCache)
}

func (this *Orderbook) ALength() int {
	return len(this.askLimitsCache)
}

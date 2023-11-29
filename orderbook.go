package rbt_orderbook

import (
	"fmt"
	"github.com/shopspring/decimal"
	"sync"
)

// maximum limits per orderbook side to pre-allocate memory
const MaxLimitsNum int = 10000

type Orderbook struct {
	Bids           *redBlackBST
	Asks           *redBlackBST
	bidLimtRwLock  sync.RWMutex
	bidLimitsCache map[decimal.Decimal]*LimitOrder
	askLimtRwLock  sync.RWMutex
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
func (this *Orderbook) getBidLimitsCacheByPrice(price decimal.Decimal) *LimitOrder {
	var limit *LimitOrder
	for k := range this.bidLimitsCache {
		if k.Equal(price) {
			limit = this.bidLimitsCache[k]
		}
	}
	return limit
}

func (this *Orderbook) getAskLimitsCacheByPrice(price decimal.Decimal) *LimitOrder {
	var limit *LimitOrder
	for k := range this.askLimitsCache {
		if k.Equal(price) {
			limit = this.askLimitsCache[k]
		}
	}
	return limit
}

func (this *Orderbook) setBidLimitsCache(limit *LimitOrder, price decimal.Decimal) {
	this.bidLimtRwLock.Lock()
	defer this.bidLimtRwLock.Unlock()
	this.bidLimitsCache[price] = limit
}
func (this *Orderbook) setAskLimitsCache(limit *LimitOrder, price decimal.Decimal) {
	this.askLimtRwLock.Lock()
	defer this.askLimtRwLock.Unlock()
	this.askLimitsCache[price] = limit
}

func (this *Orderbook) deleteBidLimitsCache(price decimal.Decimal) {
	for k := range this.bidLimitsCache {
		if k.Equal(price) {
			this.bidLimtRwLock.Lock()
			delete(this.bidLimitsCache, k)
			this.bidLimtRwLock.Unlock()
		}
	}
}
func (this *Orderbook) deleteAskLimitsCache(price decimal.Decimal) {
	for k := range this.askLimitsCache {
		if k.Equal(price) {
			this.askLimtRwLock.Lock()
			delete(this.askLimitsCache, k)
			this.askLimtRwLock.Unlock()
		}
	}
}

func (this *Orderbook) Add(price decimal.Decimal, o *Order) {
	var limit *LimitOrder

	if o.BidOrAsk {
		limit = this.getBidLimitsCacheByPrice(price)
	} else {
		limit = this.getAskLimitsCacheByPrice(price)
	}

	if limit == nil {
		// getting a new limit from pool
		limit = this.pool.Get().(*LimitOrder)
		limit.Price = price

		// insert into the corresponding BST and cache
		if o.BidOrAsk {
			this.Bids.Put(price, limit)
			this.setBidLimitsCache(limit, price)
		} else {
			this.Asks.Put(price, limit)
			this.setAskLimitsCache(limit, price)
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
			this.deleteBidLimitsCache(limit.Price)
		} else {
			this.Asks.Delete(limit.Price)
			this.deleteAskLimitsCache(limit.Price)
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
		limit = this.getBidLimitsCacheByPrice(price)
	} else {
		limit = this.getAskLimitsCacheByPrice(price)
	}

	if limit == nil {
		panic(fmt.Sprintf("there is no such price limit %+v", price))
	}

	limit.Clear()
}

func (this *Orderbook) DeleteBidLimit(price decimal.Decimal) {
	limit := this.getBidLimitsCacheByPrice(price)
	if limit == nil {
		return
	}

	this.deleteLimit(price, true)
	this.deleteBidLimitsCache(price)

	// put limit back to the pool
	limit.Clear()
	this.pool.Put(limit)

}

func (this *Orderbook) DeleteAskLimit(price decimal.Decimal) {
	limit := this.getAskLimitsCacheByPrice(price)
	if limit == nil {
		return
	}

	this.deleteLimit(price, false)
	this.deleteAskLimitsCache(price)

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
	limit := this.getBidLimitsCacheByPrice(price)
	if limit == nil {
		return decimal.Zero
	}
	return limit.TotalVolume()
}

func (this *Orderbook) GetVolumeAtAskLimit(price decimal.Decimal) decimal.Decimal {
	limit := this.getAskLimitsCacheByPrice(price)
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

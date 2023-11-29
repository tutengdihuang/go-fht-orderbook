package rbt_orderbook

import (
	"github.com/shopspring/decimal"
	"math/rand"
	"testing"
)

func TestLimitOrderEmpty(t *testing.T) {
	price := decimal.NewFromFloat(3.141593)
	l := NewLimitOrder(price)
	if l.Price.Cmp(price) != 0 || l.TotalVolume().Cmp(decimal.Zero) != 0 {
		t.Errorf("limit order initialization error")
	}
}

func TestLimitOrderAddOrder(t *testing.T) {
	price := decimal.NewFromFloat(3.141593)
	volume := decimal.NewFromFloat(25.0)
	l := NewLimitOrder(price)
	o := &Order{Volume: volume}
	l.Enqueue(o)

	if l.TotalVolume().Cmp(volume) != 0 {
		t.Errorf("total volume counted incorrectly")
		t.Errorf("tl.TotalVolume():%+v", l.TotalVolume())
		t.Errorf("volume:%+v", volume)
	}
	if l.Size() != 1 {
		t.Errorf("it should have size = 1")
	}
	if o.Limit != &l {
		t.Errorf("Parent Limit link should be set for an order")
	}
}

func TestLimitOrderAddMultipleOrders(t *testing.T) {
	price := decimal.NewFromFloat(3.141593)
	volume := decimal.Zero
	l := NewLimitOrder(price)
	n := 100
	for i := 0; i < n; i += 1 {
		o := &Order{Id: i, Volume: decimal.NewFromFloat(rand.Float64())}
		volume = volume.Add(o.Volume)
		l.Enqueue(o)
	}
	if volume.Cmp(l.TotalVolume()) != 0 {
		t.Errorf("total volume calculated incorrectly")
	}

	if l.Size() != n {
		t.Errorf("total count calculated incorrectly")
	}

	o := l.Dequeue()
	if l.TotalVolume().Cmp(volume.Sub(o.Volume)) != 0 {
		t.Errorf("total volume calculated incorrectly")
	}
	if l.Size() != n-1 {
		t.Errorf("total count calculated incorrectly")
	}
}

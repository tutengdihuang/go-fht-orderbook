package rbt_orderbook

import (
	"github.com/shopspring/decimal"
	"math/rand"
	"testing"
	//"fmt"
)

func TestIndexMinPQOne(t *testing.T) {
	minpq := NewIndexMinPQ(10)
	minpq.Insert(0, decimal.NewFromFloat(5.0))
	res := minpq.Top()

	expected := decimal.NewFromFloat(5.0)
	if res.Cmp(expected) != 0 {
		t.Errorf("actual %+v != expected %+v", res, expected)
	}
}

func TestIndexMinPQTwo(t *testing.T) {
	minpq := NewIndexMinPQ(10)
	minpq.Insert(0, decimal.NewFromFloat(6.0))
	minpq.Insert(1, decimal.NewFromFloat(5.0))

	res := [2]decimal.Decimal{}
	res[0] = minpq.Top()
	minpq.DelTop()
	res[1] = minpq.Top()

	exp := [2]decimal.Decimal{decimal.NewFromFloat(5.0), decimal.NewFromFloat(6.0)}
	if res != exp {
		t.Errorf("actual %+v != expected %+v", res, exp)
	}
}

func TestIndexMinPQThree(t *testing.T) {
	minpq := NewIndexMinPQ(10)
	minpq.Insert(0, decimal.NewFromFloat(6.0))
	minpq.Insert(1, decimal.NewFromFloat(5.0))
	minpq.Insert(2, decimal.NewFromFloat(4.0))
	res := [3]decimal.Decimal{}
	res[0] = minpq.Top()
	minpq.DelTop()
	res[1] = minpq.Top()
	minpq.DelTop()
	res[2] = minpq.Top()
	minpq.DelTop()

	exp := [3]decimal.Decimal{decimal.NewFromFloat(4.0), decimal.NewFromFloat(5.0), decimal.NewFromFloat(6.0)}
	if res != exp {
		t.Errorf("actual %+v != expected %+v", res, exp)
	}

	if !minpq.IsEmpty() {
		t.Errorf("pq should be empty")
	}
}

func TestIndexMinPQRandom(t *testing.T) {
	minpq := NewIndexMinPQ(100)
	emptyindex := 0
	for i := 0; i < 1000; i += 1 {
		emptyindex = i
		if minpq.Size() == 100 {
			emptyindex = minpq.DelTop()
		}
		minpq.Insert(emptyindex, decimal.NewFromInt(int64(rand.Intn(100))))
	}

	res := make([]decimal.Decimal, 100)
	for i := range res {
		res[i] = minpq.Top()
		minpq.DelTop()
	}

	if !minpq.IsEmpty() {
		t.Errorf("pq should be empty after all")
	}

	for i := 1; i < 100; i += 1 {
		if res[i].LessThan(res[i-1]) {
			t.Errorf("invalid order")
		}
	}
}

func BenchmarkIndexMinPQLimitedRandomInsertWithCaching(b *testing.B) {
	pq := NewIndexMinPQ(10000)

	// maximum number of levels in average is 10k
	limitslist := make([]decimal.Decimal, 10000)
	for i := range limitslist {
		limitslist[i] = decimal.NewFromFloat(rand.Float64())
	}

	// preallocate empty orders
	orders := make([]*Order, 0, b.N)
	for i := 0; i < b.N; i += 1 {
		orders = append(orders, &Order{})
	}

	// measure insertion time
	b.ResetTimer()

	limitscache := make(map[decimal.Decimal]*LimitOrder)
	for i := 0; i < b.N; i += 1 {
		// create a new order
		o := orders[i]
		o.Id = i
		o.Volume = decimal.NewFromFloat(rand.Float64())
		// o := &Order{
		// 	Id: i,
		// 	Volume: decimal.NewFromFloat(rand.Float64()),
		// }

		// set the price
		price := limitslist[rand.Intn(len(limitslist))]

		// append order to the limit price
		if limitscache[price] != nil {
			// append to the existing limit in cache
			limitscache[price].Enqueue(o)
		} else {
			// new limit
			l := NewLimitOrder(price)
			l.Enqueue(o)

			// caching limit
			limitscache[price] = &l

			// inserting into heap
			pq.Insert(len(limitscache)-1, price)
		}
	}
}

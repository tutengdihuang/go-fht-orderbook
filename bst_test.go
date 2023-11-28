package rbt_orderbook

import (
	"github.com/shopspring/decimal"
	"math/rand"
	"testing"
	//"fmt"
)

func TestBSTEmpty(t *testing.T) {
	rb := NewBST()
	if rb.Size() != 0 || !rb.IsEmpty() {
		t.Errorf("BST should be empty")
	}
}

func TestBSTBasic(t *testing.T) {
	st := NewBST()
	keys := make([]decimal.Decimal, 0)
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromFloat(rand.Float64())
		keys = append(keys, k)
		st.Put(k, nil)
	}

	if st.Size() != 10 {
		t.Errorf("size should equals 10, got %d", st.Size())
	}
	if st.IsEmpty() {
		t.Errorf("st should not be empty")
	}

	for _, k := range keys {
		if !st.Contains(k) {
			t.Errorf("st should contain the key %+v", k)
		}
	}
}

func TestBSTHeight(t *testing.T) {
	st := NewBST()
	n := 100000
	for i := 0; i < n; i += 1 {
		k := decimal.NewFromFloat(rand.Float64())
		st.Put(k, nil)
	}

	if st.Size() != n {
		t.Errorf("size should equals %d, got %d", n, st.Size())
	}
	if st.IsEmpty() {
		t.Errorf("st should not be empty")
	}

	height := st.Height()
	if height < 17 || height > 51 {
		t.Errorf("red black bst height should be in range lgN <= height <= 3*lgN, in our case from 17 to 51, but we got %d", height)
	}
}

func TestBSTMinMax(t *testing.T) {
	st := NewBST()
	for i := 0; i < 10; i += 1 {
		st.Put(decimal.NewFromInt(int64(10-i)), nil)
	}

	min := decimal.NewFromInt(1.0)
	if !st.Min().Equals(min) {
		t.Errorf("min %+v != %+v", st.Min(), min)
	}

	max := decimal.NewFromFloat(10.0)
	if !st.Max().Equal(max) {
		t.Errorf("min %+v != %+v", st.Max().Round(8), max.Round(8))
	}
}

func TestBSTMinMaxCachedOnDelete(t *testing.T) {
	st := NewBST()
	for i := 0; i < 100; i += 1 {
		st.Put(decimal.NewFromInt(int64(100-i)), nil)
	}

	min := decimal.NewFromInt(1.0)
	if !st.Min().Equals(min) {
		t.Errorf("min %+v != %+v", st.Min(), min)
	}

	max := decimal.NewFromFloat(100.0)
	if !st.Max().Equal(max) {
		t.Errorf("min %+v != %+v", st.Max(), max)
	}
	st.Delete(max)
	max = decimal.NewFromFloat(99.0)
	if !st.Max().Equal(max) {
		t.Errorf("min %+v != %+v", st.Max(), max)
	}
	if st.Size() != 99 {
		t.Errorf("size should be 99")
	}

	for i := 1; i < 20; i += 1 {
		st.Delete(decimal.NewFromInt(int64(i)))
	}
	for i := 99; i > 70; i -= 1 {
		st.Delete(decimal.NewFromInt(int64(i)))
	}

	min = decimal.NewFromFloat(20.0)
	if !st.Min().Equals(min) {
		t.Errorf("min %+v != %+v", st.Min(), min)
	}

	max = decimal.NewFromFloat(70.0)
	if !st.Max().Equal(max) {
		t.Errorf("max %+v != %+v", st.Max(), max)
	}
}

func TestBSTFloor(t *testing.T) {
	st := NewBST()
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromInt(int64(20 - 2*i))
		st.Put(k, nil)
	}

	keymiss := decimal.NewFromFloat(3.0)
	flmiss := decimal.NewFromFloat(2.0)
	if st.Floor(keymiss) != flmiss {
		t.Errorf("floor != %s", flmiss.String())
	}

	keyhit := decimal.NewFromFloat(10.0)
	if st.Floor(keyhit) != keyhit {
		t.Errorf("floor != %s", keyhit.String())
	}
}

func TestBSTCeiling(t *testing.T) {
	st := NewBST()
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromFloat(20 - 2*float64(i))
		st.Put(k, nil)
	}

	keymiss := decimal.NewFromFloat(3.0)
	clmiss := decimal.NewFromFloat(4.0)
	if st.Ceiling(keymiss) != clmiss {
		t.Errorf("ceiling != %s", clmiss.String())
	}

	keyhit := decimal.NewFromFloat(10.0)
	if st.Ceiling(keyhit) != keyhit {
		t.Errorf("ceiling != %s", keyhit.String())
	}
}

func TestBSTSelect(t *testing.T) {
	st := NewBST()
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromFloat(10 - float64(i))
		st.Put(k, nil)
	}

	key := decimal.NewFromFloat(3.0)
	if st.Select(2) != key {
		t.Errorf("element with rank=2 should be %s", key.String())
	}

	key = decimal.NewFromFloat(10.0)
	if st.Select(9) != key {
		t.Errorf("element with rank=9 should be %s", key.String())
	}
}

func TestBSTRank(t *testing.T) {
	st := NewBST()
	keys := make([]decimal.Decimal, 0)
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromFloat(10 - float64(i))
		keys = append(keys, k)
		st.Put(k, nil)
	}

	for i := range keys {
		k := st.Select(i)
		if st.Rank(k) != i {
			t.Errorf("rank of %s != %d", k.String(), i)
		}
	}

	newMax := decimal.NewFromFloat(11.0)
	if st.Rank(newMax) != len(keys) {
		t.Errorf("rank of new maximum should equal to the number of nodes in the tree")
	}

	if st.Rank(newMax) != st.Rank(decimal.NewFromFloat(12.0)) {
		t.Errorf("rank of new maximum should not depend on the new maximum concrete value")
	}
}

func TestBSTKeys(t *testing.T) {
	st := NewBST()
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromFloat(10 - float64(i))
		st.Put(k, nil)
	}

	lo := decimal.NewFromFloat(3.0)
	hi := decimal.NewFromFloat(6.0)
	keys := st.Keys(lo, hi)
	if len(keys) != 4 {
		t.Errorf("keys len should equal 4, got %+v", keys)
	}

	if !keys[0].Equal(lo) {
		t.Errorf("first key should be %s, got %s", lo.String(), keys[0].String())
	}

	if !keys[len(keys)-1].Equal(hi) {
		t.Errorf("last key should be %s, got %s", hi.String(), keys[len(keys)-1].String())
	}

	for i := 1; i < len(keys); i += 1 {
		if keys[i].LessThan(keys[i-1]) {
			t.Errorf("non-decreasing keys order validation failed")
		}
	}
}
func TestBSTDelete(t *testing.T) {
	st := NewBST()
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromInt(int64(i))
		st.Put(k, nil)
	}

	key := decimal.NewFromFloat(5.0)
	st.Delete(key)
	if st.Size() != 9 {
		t.Errorf("tree size should shrink")
	}

	if st.Contains(key) {
		t.Errorf("minimum element should be removed from the tree")
	}
}

func TestBSTPutLinkedListOrder(t *testing.T) {
	st := NewBST()
	for i := 0; i < 100; i += 1 {
		k := decimal.NewFromFloat(rand.Float64())
		st.Put(k, nil)
	}

	min := st.MinPointer()
	for p := min; p != nil && p.Next != nil; p = p.Next {
		if p.Next.Key.LessThan(p.Key) {
			t.Errorf("incorrect keys order")
			break
		}
	}
}

func TestBSTPutDeleteLinkedListOrder(t *testing.T) {
	st := NewBST()
	n := 1000
	for i := 0; i < int(n); i += 1 {
		k := decimal.NewFromFloat(rand.Float64())
		st.Put(k, nil)
	}

	// deleting random 90% of the nodes
	k := int(decimal.NewFromInt(int64(n)).Mul(decimal.NewFromFloat(0.9)).IntPart())
	for i := 0; i < k; i += 1 {
		k := st.Select(rand.Intn(st.Size()))
		st.Delete(k)
	}

	if st.Size() != n-k {
		t.Errorf("incorrect tree size %d", st.Size())
	}

	min := st.MinPointer()
	for p := min; p != nil && p.Next != nil; p = p.Next {
		if p.Next.Key.LessThan(p.Key) {
			t.Errorf("incorrect keys order")
			break
		}
	}
}

func BenchmarkBSTLimitedRandomInsertWithCaching(b *testing.B) {
	st := NewBST()

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

			// inserting into tree
			st.Put(l.Price, &l)
		}
	}
}

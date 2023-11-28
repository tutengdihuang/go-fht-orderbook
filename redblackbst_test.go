package rbt_orderbook

import (
	"github.com/shopspring/decimal"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	//"fmt"
)

func TestRedBlackEmpty(t *testing.T) {
	rb := NewRedBlackBST()
	if rb.Size() != 0 || !rb.IsEmpty() {
		t.Errorf("Red Black BST should be empty")
	}
}

func TestRedBlackBasic(t *testing.T) {
	st := NewRedBlackBST()
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

func TestRedBlackHeight(t *testing.T) {
	st := NewRedBlackBST()
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
	if height < 17 || height > 34 {
		t.Errorf("red black bst height should be in range lgN <= height <= 2*lgN, in our case from 17 to 34, but we got %d", height)
	}
}

func TestRedBlackMinMax(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		st.Put(decimal.NewFromInt(int64(10-i)), nil)
	}

	min := decimal.NewFromInt(1)
	if !st.Min().Equals(min) {
		t.Errorf("min %s != %s", st.Min().String(), min.String())
	}

	max := decimal.NewFromInt(10)
	if !st.Max().Equals(max) {
		t.Errorf("max %s != %s", st.Max().String(), max.String())
	}
}

func TestRedBlackMinMaxCachedOnDelete(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 100; i += 1 {
		st.Put(decimal.NewFromInt(int64(100-i)), nil)
	}

	min := decimal.NewFromInt(1)
	if !st.Min().Equals(min) {
		t.Errorf("min %s != %s", st.Min().String(), min.String())
	}

	max := decimal.NewFromInt(100)
	if !st.Max().Equals(max) {
		t.Errorf("max %s != %s", st.Max().String(), max.String())
	}

	st.DeleteMin()
	st.DeleteMin()
	for i := 3; i < 20; i += 1 {
		st.Delete(decimal.NewFromInt(int64(i)))
	}
	st.DeleteMax()
	st.DeleteMax()
	for i := 98; i > 70; i -= 1 {
		st.Delete(decimal.NewFromInt(int64(i)))
	}

	min = decimal.NewFromInt(20)
	if !st.Min().Equals(min) {
		t.Errorf("min %s != %s", st.Min().String(), min.String())
	}

	max = decimal.NewFromInt(70)
	if !st.Max().Equals(max) {
		t.Errorf("max %s != %s", st.Max().String(), max.String())
	}
}

func TestRedBlackFloor(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromInt(int64(20 - 2*i))
		st.Put(k, nil)
	}

	keymiss := decimal.NewFromFloat(3.0)
	flmiss := decimal.NewFromFloat(2.0)
	if !st.Floor(keymiss).Equals(flmiss) {
		t.Errorf("floor != %s", st.Floor(keymiss).String())
	}

	keyhit := decimal.NewFromFloat(10.0)
	if !st.Floor(keyhit).Equals(keyhit) {
		t.Errorf("floor != %s", st.Floor(keyhit).String())
	}
}

func TestRedBlackCeiling(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromInt(int64(20 - 2*i))
		st.Put(k, nil)
	}

	keymiss := decimal.NewFromFloat(3.0)
	clmiss := decimal.NewFromFloat(4.0)
	if !st.Ceiling(keymiss).Equals(clmiss) {
		t.Errorf("ceiling != %s", st.Ceiling(keymiss).String())
	}

	keyhit := decimal.NewFromFloat(10.0)
	if !st.Ceiling(keyhit).Equals(keyhit) {
		t.Errorf("ceiling != %s", st.Ceiling(keyhit).String())
	}
}

func TestRedBlackSelect(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromInt(int64(10 - i))
		st.Put(k, nil)
	}

	key := decimal.NewFromFloat(3.0)
	if !st.Select(2).Equals(key) {
		t.Errorf("element with rank=2 should be %s", key.String())
	}

	key = decimal.NewFromFloat(10.0)
	if !st.Select(9).Equals(key) {
		t.Errorf("element with rank=9 should be %s", key.String())
	}
}

func TestRedBlackRank(t *testing.T) {
	st := NewRedBlackBST()
	keys := make([]decimal.Decimal, 0)
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromInt(int64(10 - i))
		keys = append(keys, k)
		st.Put(k, nil)
	}

	for i := range keys {
		k := st.Select(i)
		if st.Rank(k) != i {
			t.Errorf("rank of %s != %d", k.String(), i)
		}
	}

	newMax := decimal.NewFromInt(11)
	if st.Rank(newMax) != len(keys) {
		t.Errorf("rank of new maximum should equal to the number of nodes in the tree")
	}

	if st.Rank(newMax) != st.Rank(decimal.NewFromInt(12)) {
		t.Errorf("rank of new maximum should not depend on the new maximum concrete value")
	}
}

func TestRedBlackKeys(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromInt(int64(10 - i))
		st.Put(k, nil)
	}

	lo := decimal.NewFromFloat(3.0)
	hi := decimal.NewFromFloat(6.0)
	keys := st.Keys(lo, hi)
	if len(keys) != 4 {
		t.Errorf("keys len should equal 4, %+v", keys)
	}

	if !keys[0].Equals(lo) {
		t.Errorf("first key should be %s", lo.String())
	}

	if !keys[len(keys)-1].Equals(hi) {
		t.Errorf("last key should be %s", hi.String())
	}

	for i := 1; i < len(keys); i += 1 {
		if keys[i].LessThan(keys[i-1]) {
			t.Errorf("non-decreasing keys order validation failed")
		}
	}
}

func TestRedBlackDeleteMin(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromInt(int64(10 - i))
		st.Put(k, nil)
	}

	st.DeleteMin()
	if st.Size() != 9 {
		t.Errorf("tree size should shrink")
	}

	if st.Contains(decimal.NewFromInt(1)) {
		t.Errorf("minimum element should be removed from the tree")
	}

	if !st.IsRedBlack() {
		t.Errorf("certification failed")
	}
}

func TestRedBlackDeleteMax(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := decimal.NewFromInt(int64(i))
		st.Put(k, nil)
	}

	st.DeleteMax()
	if st.Size() != 9 {
		t.Errorf("tree size should shrink")
	}

	if st.Contains(decimal.NewFromInt(9)) {
		t.Errorf("maximum element should be removed from the tree")
	}

	if !st.IsRedBlack() {
		t.Errorf("certification failed")
	}
}

func TestRedBlackDelete(t *testing.T) {
	st := NewRedBlackBST()
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
		t.Errorf("element should be removed from the tree")
	}

	if !st.IsRedBlack() {
		t.Errorf("certification failed")
	}
}

func TestRedBlackPutLinkedListOrder(t *testing.T) {
	st := NewRedBlackBST()
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

func TestRedBlackPutDeleteLinkedListOrder(t *testing.T) {
	st := NewRedBlackBST()
	n := 1000
	for i := 0; i < n; i += 1 {
		k := decimal.NewFromFloat(rand.Float64())
		st.Put(k, nil)
	}

	// deleting from both ends and in the middle 90% of the nodes
	k := int(decimal.NewFromInt(int64(n)).Mul(decimal.NewFromFloat(0.3)).IntPart())
	for i := 0; i < k; i += 1 {
		st.DeleteMin()
		k := st.Select(rand.Intn(st.Size()))
		st.Delete(k)
		st.DeleteMax()
	}

	if st.Size() != n-3*k {
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

func benchmarkRedBlackLimitedRandomInsertWithCaching(n int, b *testing.B) {
	st := NewRedBlackBST()

	// maximum number of levels in average is 10k
	limitslist := make([]decimal.Decimal, n)
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

func BenchmarkRedBlack5kLevelsRandomInsertWithCaching(b *testing.B) {
	benchmarkRedBlackLimitedRandomInsertWithCaching(5000, b)
}

func BenchmarkRedBlack10kLevelsRandomInsertWithCaching(b *testing.B) {
	benchmarkRedBlackLimitedRandomInsertWithCaching(10000, b)
}

func BenchmarkRedBlack20kLevelsRandomInsertWithCaching(b *testing.B) {
	benchmarkRedBlackLimitedRandomInsertWithCaching(20000, b)
}

func TestLimitOrder_Clear(t *testing.T) {
	type fields struct {
		Price       decimal.Decimal
		orders      *ordersQueue
		totalVolume decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &LimitOrder{
				Price:       tt.fields.Price,
				orders:      tt.fields.orders,
				totalVolume: tt.fields.totalVolume,
			}
			this.Clear()
		})
	}
}

func TestLimitOrder_Delete(t *testing.T) {
	type fields struct {
		Price       decimal.Decimal
		orders      *ordersQueue
		totalVolume decimal.Decimal
	}
	type args struct {
		o *Order
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &LimitOrder{
				Price:       tt.fields.Price,
				orders:      tt.fields.orders,
				totalVolume: tt.fields.totalVolume,
			}
			this.Delete(tt.args.o)
		})
	}
}

func TestLimitOrder_Dequeue(t *testing.T) {
	type fields struct {
		Price       decimal.Decimal
		orders      *ordersQueue
		totalVolume decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   *Order
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &LimitOrder{
				Price:       tt.fields.Price,
				orders:      tt.fields.orders,
				totalVolume: tt.fields.totalVolume,
			}
			if got := this.Dequeue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dequeue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimitOrder_Enqueue(t *testing.T) {
	type fields struct {
		Price       decimal.Decimal
		orders      *ordersQueue
		totalVolume decimal.Decimal
	}
	type args struct {
		o *Order
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &LimitOrder{
				Price:       tt.fields.Price,
				orders:      tt.fields.orders,
				totalVolume: tt.fields.totalVolume,
			}
			this.Enqueue(tt.args.o)
		})
	}
}

func TestLimitOrder_Size(t *testing.T) {
	type fields struct {
		Price       decimal.Decimal
		orders      *ordersQueue
		totalVolume decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &LimitOrder{
				Price:       tt.fields.Price,
				orders:      tt.fields.orders,
				totalVolume: tt.fields.totalVolume,
			}
			if got := this.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimitOrder_TotalVolume(t *testing.T) {
	type fields struct {
		Price       decimal.Decimal
		orders      *ordersQueue
		totalVolume decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &LimitOrder{
				Price:       tt.fields.Price,
				orders:      tt.fields.orders,
				totalVolume: tt.fields.totalVolume,
			}
			if got := this.TotalVolume(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TotalVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBST(t *testing.T) {
	tests := []struct {
		name string
		want bst
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBST(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBST() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIndexMinPQ(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want indexMinPQ
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIndexMinPQ(tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIndexMinPQ() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLimitOrder(t *testing.T) {
	type args struct {
		price decimal.Decimal
	}
	tests := []struct {
		name string
		args args
		want LimitOrder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLimitOrder(tt.args.price); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLimitOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMinPQ(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want minPQ
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMinPQ(tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMinPQ() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOrderbook(t *testing.T) {
	tests := []struct {
		name string
		want Orderbook
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrderbook(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderbook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOrdersQueue(t *testing.T) {
	tests := []struct {
		name string
		want ordersQueue
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrdersQueue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrdersQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRedBlackBST(t *testing.T) {
	tests := []struct {
		name string
		want redBlackBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRedBlackBST(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRedBlackBST() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderbook_ALength(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			if got := this.ALength(); got != tt.want {
				t.Errorf("ALength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderbook_Add(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	type args struct {
		price decimal.Decimal
		o     *Order
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			this.Add(tt.args.price, tt.args.o)
		})
	}
}

func TestOrderbook_BLength(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			if got := this.BLength(); got != tt.want {
				t.Errorf("BLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderbook_Cancel(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	type args struct {
		o *Order
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			this.Cancel(tt.args.o)
		})
	}
}

func TestOrderbook_ClearAskLimit(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	type args struct {
		price decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			this.ClearAskLimit(tt.args.price)
		})
	}
}

func TestOrderbook_ClearBidLimit(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	type args struct {
		price decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			this.ClearBidLimit(tt.args.price)
		})
	}
}

func TestOrderbook_DeleteAskLimit(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	type args struct {
		price decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			this.DeleteAskLimit(tt.args.price)
		})
	}
}

func TestOrderbook_DeleteBidLimit(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	type args struct {
		price decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			this.DeleteBidLimit(tt.args.price)
		})
	}
}

func TestOrderbook_GetBestBid(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	tests := []struct {
		name   string
		fields fields
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			if got := this.GetBestBid(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBestBid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderbook_GetBestOffer(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	tests := []struct {
		name   string
		fields fields
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			if got := this.GetBestOffer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBestOffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderbook_GetVolumeAtAskLimit(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	type args struct {
		price decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			if got := this.GetVolumeAtAskLimit(tt.args.price); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVolumeAtAskLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderbook_GetVolumeAtBidLimit(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	type args struct {
		price decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			if got := this.GetVolumeAtBidLimit(tt.args.price); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVolumeAtBidLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderbook_clearLimit(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	type args struct {
		price    decimal.Decimal
		bidOrAsk bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			this.clearLimit(tt.args.price, tt.args.bidOrAsk)
		})
	}
}

func TestOrderbook_deleteLimit(t *testing.T) {
	type fields struct {
		Bids           *redBlackBST
		Asks           *redBlackBST
		bidLimitsCache map[decimal.Decimal]*LimitOrder
		askLimitsCache map[decimal.Decimal]*LimitOrder
		pool           *sync.Pool
	}
	type args struct {
		price    decimal.Decimal
		bidOrAsk bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Orderbook{
				Bids:           tt.fields.Bids,
				Asks:           tt.fields.Asks,
				bidLimitsCache: tt.fields.bidLimitsCache,
				askLimitsCache: tt.fields.askLimitsCache,
				pool:           tt.fields.pool,
			}
			this.deleteLimit(tt.args.price, tt.args.bidOrAsk)
		})
	}
}

func Test_bst_Ceiling(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Ceiling(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Ceiling() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_Contains(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Contains(tt.args.key); got != tt.want {
				t1.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_Delete(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.Delete(tt.args.key)
		})
	}
}

func Test_bst_Floor(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Floor(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Floor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_Get(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *LimitOrder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_Height(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Height(); got != tt.want {
				t1.Errorf("Height() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_IsEmpty(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.IsEmpty(); got != tt.want {
				t1.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_Keys(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		lo decimal.Decimal
		hi decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Keys(tt.args.lo, tt.args.hi); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_Max(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Max(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_MaxPointer(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		want   *nodeBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.MaxPointer(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("MaxPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_MaxValue(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		want   *LimitOrder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.MaxValue(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("MaxValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_Min(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Min(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_MinPointer(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		want   *nodeBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.MinPointer(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("MinPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_MinValue(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		want   *LimitOrder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.MinValue(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("MinValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_Print(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.Print()
		})
	}
}

func Test_bst_Put(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		key   decimal.Decimal
		value *LimitOrder
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.Put(tt.args.key, tt.args.value)
		})
	}
}

func Test_bst_Rank(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Rank(tt.args.key); got != tt.want {
				t1.Errorf("Rank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_Select(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		k int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Select(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_Size(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Size(); got != tt.want {
				t1.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_ceiling(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n   *nodeBST
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.ceiling(tt.args.n, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("ceiling() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_delete(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n   *nodeBST
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.delete(tt.args.n, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_deleteMin(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.deleteMin(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("deleteMin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_floor(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n   *nodeBST
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.floor(tt.args.n, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("floor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_get(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n   *nodeBST
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.get(tt.args.n, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_height(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.height(tt.args.n); got != tt.want {
				t1.Errorf("height() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_keys(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n  *nodeBST
		lo decimal.Decimal
		hi decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.keys(tt.args.n, tt.args.lo, tt.args.hi); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_max(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.max(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_min(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.min(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_panicIfEmpty(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.panicIfEmpty()
		})
	}
}

func Test_bst_print(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.print(tt.args.n)
		})
	}
}

func Test_bst_put(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n     *nodeBST
		key   decimal.Decimal
		value *LimitOrder
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.put(tt.args.n, tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("put() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_rank(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n   *nodeBST
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.rank(tt.args.n, tt.args.key); got != tt.want {
				t1.Errorf("rank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_selectNode(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n *nodeBST
		k int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeBST
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.selectNode(tt.args.n, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("selectNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bst_size(t1 *testing.T) {
	type fields struct {
		root *nodeBST
		minC *nodeBST
		maxC *nodeBST
	}
	type args struct {
		n *nodeBST
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &bst{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.size(tt.args.n); got != tt.want {
				t1.Errorf("size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexMinPQ_Change(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	type args struct {
		i   int
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			pq.Change(tt.args.i, tt.args.key)
		})
	}
}

func Test_indexMinPQ_Contains(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			if got := pq.Contains(tt.args.i); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexMinPQ_DelTop(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			if got := pq.DelTop(); got != tt.want {
				t.Errorf("DelTop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexMinPQ_Delete(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			pq.Delete(tt.args.i)
		})
	}
}

func Test_indexMinPQ_Insert(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	type args struct {
		i   int
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			pq.Insert(tt.args.i, tt.args.key)
		})
	}
}

func Test_indexMinPQ_IsEmpty(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			if got := pq.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexMinPQ_Size(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			if got := pq.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexMinPQ_Top(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	tests := []struct {
		name   string
		fields fields
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			if got := pq.Top(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Top() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexMinPQ_TopIndex(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			if got := pq.TopIndex(); got != tt.want {
				t.Errorf("TopIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexMinPQ_checkIndex(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			pq.checkIndex(tt.args.i)
		})
	}
}

func Test_indexMinPQ_sink(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			pq.sink(tt.args.i)
		})
	}
}

func Test_indexMinPQ_swim(t *testing.T) {
	type fields struct {
		keys         []decimal.Decimal
		index2offset []int
		offset2index []int
		n            int
	}
	type args struct {
		i int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &indexMinPQ{
				keys:         tt.fields.keys,
				index2offset: tt.fields.index2offset,
				offset2index: tt.fields.offset2index,
				n:            tt.fields.n,
			}
			pq.swim(tt.args.i)
		})
	}
}

func Test_minPQ_DelTop(t *testing.T) {
	type fields struct {
		keys []decimal.Decimal
		n    int
	}
	tests := []struct {
		name   string
		fields fields
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &minPQ{
				keys: tt.fields.keys,
				n:    tt.fields.n,
			}
			if got := pq.DelTop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DelTop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_minPQ_Insert(t *testing.T) {
	type fields struct {
		keys []decimal.Decimal
		n    int
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &minPQ{
				keys: tt.fields.keys,
				n:    tt.fields.n,
			}
			pq.Insert(tt.args.key)
		})
	}
}

func Test_minPQ_IsEmpty(t *testing.T) {
	type fields struct {
		keys []decimal.Decimal
		n    int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &minPQ{
				keys: tt.fields.keys,
				n:    tt.fields.n,
			}
			if got := pq.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_minPQ_Size(t *testing.T) {
	type fields struct {
		keys []decimal.Decimal
		n    int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &minPQ{
				keys: tt.fields.keys,
				n:    tt.fields.n,
			}
			if got := pq.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_minPQ_Top(t *testing.T) {
	type fields struct {
		keys []decimal.Decimal
		n    int
	}
	tests := []struct {
		name   string
		fields fields
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &minPQ{
				keys: tt.fields.keys,
				n:    tt.fields.n,
			}
			if got := pq.Top(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Top() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_minPQ_sink(t *testing.T) {
	type fields struct {
		keys []decimal.Decimal
		n    int
	}
	type args struct {
		k int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &minPQ{
				keys: tt.fields.keys,
				n:    tt.fields.n,
			}
			pq.sink(tt.args.k)
		})
	}
}

func Test_minPQ_swim(t *testing.T) {
	type fields struct {
		keys []decimal.Decimal
		n    int
	}
	type args struct {
		k int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &minPQ{
				keys: tt.fields.keys,
				n:    tt.fields.n,
			}
			pq.swim(tt.args.k)
		})
	}
}

func Test_ordersQueue_Delete(t *testing.T) {
	type fields struct {
		head *Order
		tail *Order
		size int
	}
	type args struct {
		o *Order
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ordersQueue{
				head: tt.fields.head,
				tail: tt.fields.tail,
				size: tt.fields.size,
			}
			this.Delete(tt.args.o)
		})
	}
}

func Test_ordersQueue_Dequeue(t *testing.T) {
	type fields struct {
		head *Order
		tail *Order
		size int
	}
	tests := []struct {
		name   string
		fields fields
		want   *Order
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ordersQueue{
				head: tt.fields.head,
				tail: tt.fields.tail,
				size: tt.fields.size,
			}
			if got := this.Dequeue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dequeue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ordersQueue_Enqueue(t *testing.T) {
	type fields struct {
		head *Order
		tail *Order
		size int
	}
	type args struct {
		o *Order
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ordersQueue{
				head: tt.fields.head,
				tail: tt.fields.tail,
				size: tt.fields.size,
			}
			this.Enqueue(tt.args.o)
		})
	}
}

func Test_ordersQueue_IsEmpty(t *testing.T) {
	type fields struct {
		head *Order
		tail *Order
		size int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ordersQueue{
				head: tt.fields.head,
				tail: tt.fields.tail,
				size: tt.fields.size,
			}
			if got := this.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ordersQueue_Size(t *testing.T) {
	type fields struct {
		head *Order
		tail *Order
		size int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &ordersQueue{
				head: tt.fields.head,
				tail: tt.fields.tail,
				size: tt.fields.size,
			}
			if got := this.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_Ceiling(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Ceiling(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Ceiling() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_Contains(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Contains(tt.args.key); got != tt.want {
				t1.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_Delete(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.Delete(tt.args.key)
		})
	}
}

func Test_redBlackBST_DeleteMax(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.DeleteMax()
		})
	}
}

func Test_redBlackBST_DeleteMin(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.DeleteMin()
		})
	}
}

func Test_redBlackBST_Floor(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Floor(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Floor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_Get(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *LimitOrder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_Height(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Height(); got != tt.want {
				t1.Errorf("Height() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_IsEmpty(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.IsEmpty(); got != tt.want {
				t1.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_IsRedBlack(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.IsRedBlack(); got != tt.want {
				t1.Errorf("IsRedBlack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_Keys(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		lo decimal.Decimal
		hi decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Keys(tt.args.lo, tt.args.hi); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_Max(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Max(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_MaxPointer(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.MaxPointer(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("MaxPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_MaxValue(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		want   *LimitOrder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.MaxValue(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("MaxValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_Min(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Min(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_MinPointer(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.MinPointer(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("MinPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_MinValue(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		want   *LimitOrder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.MinValue(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("MinValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_Print(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.Print()
		})
	}
}

func Test_redBlackBST_Put(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		key   decimal.Decimal
		value *LimitOrder
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.Put(tt.args.key, tt.args.value)
		})
	}
}

func Test_redBlackBST_Rank(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Rank(tt.args.key); got != tt.want {
				t1.Errorf("Rank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_Select(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		k int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Select(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_Size(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.Size(); got != tt.want {
				t1.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_ceiling(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n   *nodeRedBlack
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.ceiling(tt.args.n, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("ceiling() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_delete(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n   *nodeRedBlack
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.delete(tt.args.n, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_deleteMax(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.deleteMax(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("deleteMax() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_deleteMin(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.deleteMin(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("deleteMin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_flipColors(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.flipColors(tt.args.n)
		})
	}
}

func Test_redBlackBST_floor(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n   *nodeRedBlack
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.floor(tt.args.n, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("floor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_get(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n   *nodeRedBlack
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.get(tt.args.n, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_height(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.height(tt.args.n); got != tt.want {
				t1.Errorf("height() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_is23(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.is23(tt.args.n); got != tt.want {
				t1.Errorf("is23() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_isBalanced(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			got, got1 := t.isBalanced(tt.args.n)
			if got != tt.want {
				t1.Errorf("isBalanced() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t1.Errorf("isBalanced() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_redBlackBST_isRed(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.isRed(tt.args.n); got != tt.want {
				t1.Errorf("isRed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_keys(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n  *nodeRedBlack
		lo decimal.Decimal
		hi decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.keys(tt.args.n, tt.args.lo, tt.args.hi); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_max(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.max(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_min(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.min(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_moveRedLeft(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.moveRedLeft(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("moveRedLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_moveRedRight(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.moveRedRight(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("moveRedRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_panicIfEmpty(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.panicIfEmpty()
		})
	}
}

func Test_redBlackBST_print(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			t.print(tt.args.n)
		})
	}
}

func Test_redBlackBST_put(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n     *nodeRedBlack
		key   decimal.Decimal
		value *LimitOrder
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.put(tt.args.n, tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("put() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_rank(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n   *nodeRedBlack
		key decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.rank(tt.args.n, tt.args.key); got != tt.want {
				t1.Errorf("rank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_rotateLeft(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.rotateLeft(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("rotateLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_rotateRight(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.rotateRight(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("rotateRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_selectNode(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
		k int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *nodeRedBlack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.selectNode(tt.args.n, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("selectNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redBlackBST_size(t1 *testing.T) {
	type fields struct {
		root *nodeRedBlack
		minC *nodeRedBlack
		maxC *nodeRedBlack
	}
	type args struct {
		n *nodeRedBlack
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &redBlackBST{
				root: tt.fields.root,
				minC: tt.fields.minC,
				maxC: tt.fields.maxC,
			}
			if got := t.size(tt.args.n); got != tt.want {
				t1.Errorf("size() = %v, want %v", got, tt.want)
			}
		})
	}
}

package rbt_orderbook

import (
	"github.com/shopspring/decimal"
	"math/rand"
	"testing"
	//"fmt"
)

func TestOrderbookEmpty(t *testing.T) {
	b := NewOrderbook()
	if b.BLength() != 0 {
		t.Errorf("book should be empty")
	}
	if b.ALength() != 0 {
		t.Errorf("book should be empty")
	}
}

func TestOrderbookAddOne(t *testing.T) {
	b := NewOrderbook()
	bid := &Order{
		BidOrAsk: true,
	}
	ask := &Order{
		BidOrAsk: false,
	}
	b.Add(decimal.NewFromFloat(1.0), bid)
	b.Add(decimal.NewFromFloat(2.0), ask)
	if b.BLength() != 1 {
		t.Errorf("book should have 1 bid")
	}
	if b.ALength() != 1 {
		t.Errorf("book should have 1 ask")
	}
}

func TestOrderbookAddMultiple(t *testing.T) {
	b := NewOrderbook()
	for i := 0; i < 100; i += 1 {
		bid := &Order{
			BidOrAsk: true,
		}
		b.Add(decimal.NewFromInt(int64(i)), bid)
	}

	for i := 100; i < 200; i += 1 {
		bid := &Order{
			BidOrAsk: false,
		}
		b.Add(decimal.NewFromInt(int64(i)), bid)
	}

	if b.BLength() != 100 {
		t.Errorf("book should have 100 bids")
	}
	if b.ALength() != 100 {
		t.Errorf("book should have 100 asks")
	}

	expectedBestBid := decimal.NewFromInt(99)
	if !b.GetBestBid().Equal(expectedBestBid) {
		t.Errorf("best bid should be %s", expectedBestBid.String())
	}

	expectedBestOffer := decimal.NewFromInt(100)
	if !b.GetBestOffer().Equal(expectedBestOffer) {
		t.Errorf("best offer should be %s", expectedBestOffer.String())
	}
}

func TestOrderbookAddAndCancel(t *testing.T) {
	b := NewOrderbook()
	bid1 := &Order{
		Id:       1,
		BidOrAsk: true,
	}
	bid2 := &Order{
		Id:       2,
		BidOrAsk: true,
	}
	b.Add(decimal.NewFromFloat(1.0), bid1)
	b.Add(decimal.NewFromFloat(2.0), bid2)
	if b.GetBestBid().Cmp(decimal.NewFromFloat(2.0)) != 0 {
		t.Errorf("best bid should be 2.0")
	}
	b.Cancel(bid2)
	if b.GetBestBid().Cmp(decimal.NewFromFloat(1.0)) != 0 {
		t.Errorf("best bid should be 1.0 now")
	}
}

func TestGetVolumeAtLimit(t *testing.T) {
	b := NewOrderbook()
	bid1 := &Order{
		Id:       1,
		BidOrAsk: true,
		Volume:   decimal.NewFromFloat(0.1),
	}
	bid2 := &Order{
		Id:       2,
		BidOrAsk: true,
		Volume:   decimal.NewFromFloat(0.2),
	}
	b.Add(decimal.NewFromFloat(1.0), bid1)
	b.Add(decimal.NewFromFloat(1.0), bid2)
	if b.GetVolumeAtBidLimit(decimal.NewFromFloat(1.0)).Sub(decimal.NewFromFloat(0.3)).Abs().GreaterThan(decimal.NewFromFloat(0.0000001)) {
		t.Errorf("invalid volume at limit: %+v", b.GetVolumeAtBidLimit(decimal.NewFromFloat(1.0)))
	}
}

func benchmarkOrderbookLimitedRandomInsert(n int, b *testing.B) {
	book := NewOrderbook()

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
	for i := 0; i < b.N; i += 1 {
		price := limitslist[rand.Intn(len(limitslist))]

		// create a new order
		o := orders[i]
		o.Id = i
		o.Volume = decimal.NewFromFloat(rand.Float64())
		o.BidOrAsk = price.LessThan(decimal.NewFromFloat(0.5))

		// add to the book
		book.Add(price, o)
	}

	//fmt.Printf("bid size %d, ask size %d\n", book.BLength(), book.ALength())
}

func BenchmarkOrderbook5kLevelsRandomInsert(b *testing.B) {
	benchmarkOrderbookLimitedRandomInsert(5000, b)
}

func BenchmarkOrderbook10kLevelsRandomInsert(b *testing.B) {
	benchmarkOrderbookLimitedRandomInsert(10000, b)
}

func BenchmarkOrderbook20kLevelsRandomInsert(b *testing.B) {
	benchmarkOrderbookLimitedRandomInsert(20000, b)
}

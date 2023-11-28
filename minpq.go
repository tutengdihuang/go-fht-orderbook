package rbt_orderbook

import "github.com/shopspring/decimal"

// Mininum oriented Priority Queue
type minPQ struct {
	keys []decimal.Decimal
	n    int
}

func NewMinPQ(size int) minPQ {
	return minPQ{
		keys: make([]decimal.Decimal, size+1),
	}
}

func (pq *minPQ) Size() int {
	return pq.n
}

func (pq *minPQ) IsEmpty() bool {
	return pq.n == 0
}

func (pq *minPQ) Insert(key decimal.Decimal) {
	if pq.n+1 == cap(pq.keys) {
		panic("pq is full")
	}

	pq.n++
	pq.keys[pq.n] = key

	// restore order: LogN
	pq.swim(pq.n)
}

func (pq *minPQ) Top() decimal.Decimal {
	if pq.IsEmpty() {
		panic("pq is empty")
	}

	return pq.keys[1]
}

// removes minimal element and returns it
func (pq *minPQ) DelTop() decimal.Decimal {
	if pq.IsEmpty() {
		panic("pq is empty")
	}

	top := pq.keys[1]
	pq.keys[1] = pq.keys[pq.n]
	pq.n--

	// restore order: logN
	pq.sink(1)

	return top
}

func (pq *minPQ) swim(k int) {
	for k > 1 && pq.keys[k].LessThan(pq.keys[k/2]) {
		// swap
		pq.keys[k], pq.keys[k/2] = pq.keys[k/2], pq.keys[k]
		k = k / 2
	}
}

func (pq *minPQ) sink(k int) {
	for 2*k <= pq.n {
		c := 2 * k
		// select minimum of two children
		if c < pq.n && pq.keys[c+1].LessThan(pq.keys[c]) {
			c++
		}

		if pq.keys[c].LessThan(pq.keys[k]) {
			// swap
			pq.keys[c], pq.keys[k] = pq.keys[k], pq.keys[c]
			k = c
		} else {
			break
		}
	}
}

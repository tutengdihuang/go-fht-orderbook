package rbt_orderbook

import (
	"fmt"
	"github.com/shopspring/decimal"
)

// Simple Binary Search Tree, not self-balancing, good for random input

type nodeBST struct {
	Key   decimal.Decimal
	Value *LimitOrder
	Next  *nodeBST
	Prev  *nodeBST

	left  *nodeBST
	right *nodeBST
	size  int
}

type bst struct {
	root *nodeBST
	minC *nodeBST // cached min/max keys for O(1) access
	maxC *nodeBST
}

func NewBST() bst {
	return bst{}
}

func (t *bst) Size() int {
	return t.size(t.root)
}

func (t *bst) size(n *nodeBST) int {
	if n == nil {
		return 0
	}

	return n.size
}

func (t *bst) IsEmpty() bool {
	return t.size(t.root) == 0
}

func (t *bst) panicIfEmpty() {
	if t.IsEmpty() {
		panic("BST is empty")
	}
}

func (t *bst) Contains(key decimal.Decimal) bool {
	return t.get(t.root, key) != nil
}

func (t *bst) Get(key decimal.Decimal) *LimitOrder {
	t.panicIfEmpty()

	x := t.get(t.root, key)
	if x == nil {
		panic(fmt.Sprintf("key %+v does not exist", key))
	}

	return x.Value
}

func (t *bst) get(n *nodeBST, key decimal.Decimal) *nodeBST {
	if n == nil {
		return nil
	}

	if n.Key == key {
		return n
	}

	if n.Key.GreaterThan(key) {
		return t.get(n.left, key)
	} else {
		return t.get(n.right, key)
	}
}

func (t *bst) Put(key decimal.Decimal, value *LimitOrder) {
	t.root = t.put(t.root, key, value)
}

func (t *bst) put(n *nodeBST, key decimal.Decimal, value *LimitOrder) *nodeBST {
	if n == nil {
		// search miss, creating a new node
		n := &nodeBST{
			Value: value,
			Key:   key,
			size:  1,
		}

		if t.minC == nil || key.LessThan(t.minC.Key) {
			// new min
			t.minC = n
		}
		if t.maxC == nil || key.GreaterThan(t.maxC.Key) {
			// new max
			t.maxC = n
		}

		return n
	}

	if n.Key == key {
		// search hit, updating the value
		n.Value = value
		return n
	}

	if n.Key.GreaterThan(key) {
		left := n.left
		n.left = t.put(n.left, key, value)
		if left == nil {
			// new node has been just inserted to the left
			prev := n.Prev
			if prev != nil {
				prev.Next = n.left
			}
			n.left.Prev = prev
			n.left.Next = n
			n.Prev = n.left
		}
	} else {
		right := n.right
		n.right = t.put(n.right, key, value)
		if right == nil {
			// new node has been just inserted to the right
			next := n.Next
			if next != nil {
				next.Prev = n.right
			}
			n.right.Next = next
			n.right.Prev = n
			n.Next = n.right
		}
	}

	// re-calc size
	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *bst) Height() int {
	if t.IsEmpty() {
		return 0
	}

	return t.height(t.root)
}

func (t *bst) height(n *nodeBST) int {
	if n == nil {
		return 0
	}

	lheight := t.height(n.left)
	rheight := t.height(n.right)

	height := lheight
	if rheight > lheight {
		height = rheight
	}

	return height + 1
}

func (t *bst) Min() decimal.Decimal {
	t.panicIfEmpty()
	return t.minC.Key
}

func (t *bst) MinValue() *LimitOrder {
	t.panicIfEmpty()
	return t.minC.Value
}

func (t *bst) MinPointer() *nodeBST {
	t.panicIfEmpty()
	return t.minC
}

func (t *bst) min(n *nodeBST) *nodeBST {
	if n.left == nil {
		return n
	}

	return t.min(n.left)
}

func (t *bst) Max() decimal.Decimal {
	t.panicIfEmpty()
	return t.maxC.Key
}

func (t *bst) MaxValue() *LimitOrder {
	t.panicIfEmpty()
	return t.maxC.Value
}

func (t *bst) MaxPointer() *nodeBST {
	t.panicIfEmpty()
	return t.maxC
}

func (t *bst) max(n *nodeBST) *nodeBST {
	if n.right == nil {
		return n
	}

	return t.max(n.right)
}

func (t *bst) Floor(key decimal.Decimal) decimal.Decimal {
	t.panicIfEmpty()

	floor := t.floor(t.root, key)
	if floor == nil {
		panic(fmt.Sprintf("there are no keys <= %+v", key))
	}

	return floor.Key
}

func (t *bst) floor(n *nodeBST, key decimal.Decimal) *nodeBST {
	if n == nil {
		// search miss
		return nil
	}

	if n.Key == key {
		// search hit
		return n
	}

	if n.Key.GreaterThan(key) {
		// floor must be in the left sub-tree
		return t.floor(n.left, key)
	}

	// key could be in the right sub-tree, if not, using current root
	floor := t.floor(n.right, key)
	if floor != nil {
		return floor
	}

	return n
}

func (t *bst) Ceiling(key decimal.Decimal) decimal.Decimal {
	t.panicIfEmpty()

	ceiling := t.ceiling(t.root, key)
	if ceiling == nil {
		panic(fmt.Sprintf("there are no keys >= %+v", key))
	}

	return ceiling.Key
}

func (t *bst) ceiling(n *nodeBST, key decimal.Decimal) *nodeBST {
	if n == nil {
		// search miss
		return nil
	}

	if n.Key == key {
		// search hit
		return n
	}

	if n.Key.LessThan(key) {
		// ceiling must be in the right sub-tree
		return t.ceiling(n.right, key)
	}

	// the key could be in the left sub-tree, if not, using current root
	ceiling := t.ceiling(n.left, key)
	if ceiling != nil {
		return ceiling
	}

	return n
}

func (t *bst) Select(k int) decimal.Decimal {
	if k < 0 || k >= t.Size() {
		panic("index out of range")
	}

	return t.selectNode(t.root, k).Key
}

func (t *bst) selectNode(n *nodeBST, k int) *nodeBST {
	if t.size(n.left) == k {
		return n
	}

	if t.size(n.left) > k {
		return t.selectNode(n.left, k)
	}

	k = k - t.size(n.left) - 1
	return t.selectNode(n.right, k)
}

func (t *bst) Rank(key decimal.Decimal) int {
	t.panicIfEmpty()
	return t.rank(t.root, key)
}

func (t *bst) rank(n *nodeBST, key decimal.Decimal) int {
	if n == nil {
		return 0
	}

	if n.Key == key {
		return t.size(n.left)
	}

	if n.Key.GreaterThan(key) {
		return t.rank(n.left, key)
	}

	return t.size(n.left) + 1 + t.rank(n.right, key)
}

func (t *bst) deleteMin(n *nodeBST) *nodeBST {
	if n == nil {
		return nil
	}

	if n.left == nil {
		// we've reached the least leave of the tree
		next := n.Next
		prev := n.Prev
		if prev != nil {
			prev.Next = next
		}
		if next != nil {
			next.Prev = prev
		}
		n.Next = nil
		n.Prev = nil

		// updating global min
		if t.minC == n {
			t.minC = next
		}

		return n.right
	}

	n.left = t.deleteMin(n.left)

	// update size
	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *bst) Delete(key decimal.Decimal) {
	t.panicIfEmpty()

	t.root = t.delete(t.root, key)
}

func (t *bst) delete(n *nodeBST, key decimal.Decimal) *nodeBST {
	if n == nil {
		return nil
	}

	if n.Key.Equal(key) {
		// search hit

		// updating linked list
		next := n.Next
		prev := n.Prev
		if prev != nil {
			prev.Next = next
		}
		if next != nil {
			next.Prev = prev
		}
		n.Next = nil
		n.Prev = nil

		// updating global min and max
		if t.minC == n {
			t.minC = next
		}
		if t.maxC == n {
			t.maxC = prev
		}

		// replacing by successor (we can do similar with precedessor)
		if n.left == nil {
			return n.right
		} else if n.right == nil {
			return n.left
		}

		newn := t.min(n.right)
		newn.right = t.deleteMin(n.right)
		newn.left = n.left
		n = newn
	} else if n.Key.GreaterThan(key) {
		n.left = t.delete(n.left, key)
	} else {
		n.right = t.delete(n.right, key)
	}

	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *bst) Keys(lo, hi decimal.Decimal) []decimal.Decimal {
	if lo.LessThan(t.Min()) || hi.GreaterThan(t.Max()) {
		panic("keys out of range")
	}

	return t.keys(t.root, lo, hi)
}

func (t *bst) keys(n *nodeBST, lo, hi decimal.Decimal) []decimal.Decimal {
	if n == nil {
		return nil
	}

	if n.Key.LessThan(lo) {
		return t.keys(n.right, lo, hi)
	} else if n.Key.GreaterThan(hi) {
		return t.keys(n.left, lo, hi)
	}

	l := t.keys(n.left, lo, hi)
	r := t.keys(n.right, lo, hi)

	keys := make([]decimal.Decimal, 0)
	if l != nil {
		keys = append(keys, l...)
	}
	keys = append(keys, n.Key)
	if r != nil {
		keys = append(keys, r...)
	}

	return keys
}

func (t *bst) Print() {
	fmt.Println()
	t.print(t.root)
	fmt.Println()
}

func (t *bst) print(n *nodeBST) {
	if n == nil {
		return
	}

	fmt.Printf("%+v ", n.Key)

	t.print(n.left)
	t.print(n.right)
}

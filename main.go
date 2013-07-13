package lockfreehash

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type node struct {
	rawkey Key
	value  interface{}
	key    uint32
	next   *node
}
type Hash struct {
	array []*node
	bits  uint32
	num   uint32
	ratio float32
}

type Key interface {
	GetHash() uint32
	Equal(Key) bool
}

func New() *Hash {
	ret := new(Hash)
	ret.bits = 5
	ret.array = make([]*node, 1<<ret.bits)
	ret.array[0] = newSentry(0)
	ret.num = 0
	ret.ratio = 0.6
	return ret
}

func (h *Hash) initBucket(idx uint32) {
	parent := idx >> 1
	if h.array[parent] == nil {
		h.initBucket(parent)
	}
	fmt.Println("dead lock in initBucket??")
	sentry := newSentry(idx)
	h.array[idx] = listInsert(h.array[parent], sentry)
}

// insert p to the sorted list , which begins with head
// if the list already contains a node with p's key, return that node and discard p
func listInsert(head, p *node) *node {
	for {
		next := head
		prev := head
		for {
			if next == nil || next.key > p.key {
				break
			}
			if next.key == p.key {
				// insert a sentry and it's already in the list
				if p.rawkey == nil && next.rawkey == nil {
					return next
				}

				// update a key-value node
				if next.rawkey.Equal(p.rawkey) {
					next.value = p.value
					return next
				}

				// hash conflict... two thing hava the same hash val
				// move on ... append the node to the end of these
			}
			prev = next
			next = next.next
		}

		// insert p between prev and next
		tmp := unsafe.Pointer(prev.next)
		if atomic.CompareAndSwapPointer(&tmp, unsafe.Pointer(next), unsafe.Pointer(p)) == true {
			return p
		} else {
			head = prev
		}
	}
}

func (h *Hash) Put(rawkey Key, value interface{}) {
	hash := rawkey.GetHash()
	mask := uint32((1 << h.bits) - 1)
	idx := hash & mask
	key := bitReverse(hash) & 1

	n := &node{
		key:    key,
		value:  value,
		rawkey: rawkey,
	}

	if h.array[idx] == nil {
		// initialize this slot. it doesn't matter that different thread insert the sentry simentanlly, just one will success
		h.initBucket(idx)
	}
	listInsert(h.array[idx], n)
	h.num++

	caps := 1 << h.bits
	if float32(h.num) > h.ratio*float32(caps) {
		reHash(h)
	}
}

func reHash(h *Hash) {

}

func bitReverse(v uint32) uint32 {
	mask := []uint32{
		0x55555555, //...010101010101
		0xaaaaaaaa, //...101010101010
		0x33333333, //...001100110011
		0xcccccccc, //...110011001100
		0x0f0f0f0f, //...00001111
		0xf0f0f0f0, //...11110000
		0x00ff00ff, //...0000000011111111
		0xff00ff00, //...1111111100000000
		0x0000ffff,
		0xffff0000,
	}

	for i := 0; i < 5; i++ {
		tmp1 := v & mask[2*i]
		tmp2 := v & mask[2*i+1]
		v = tmp1 | tmp2
	}
	return v
}

func newSentry(idx uint32) *node {
	tmp := bitReverse(idx)
	sentry := &node{
		key: tmp,
	}
	return sentry
}

func (h *Hash) Delete(rawkey Key) {

}

func (h *Hash) Get(rawkey Key) (interface{}, bool) {
	hash := rawkey.GetHash()
	mask := uint32((1 << h.bits) - 1)
	key := bitReverse(hash) & 1
	idx := hash & mask

	if h.array[idx] == nil {
		return nil, false
	}

	sentry := h.array[idx]
	for n := sentry.next; n != nil && (n.key&1) != 0; n = n.next {
		if n.key == key && rawkey.Equal(n.rawkey) {
			return n.value, true
		}
	}
	return nil, false
}

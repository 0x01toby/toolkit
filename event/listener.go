package event

import (
	"fmt"
	"sort"
)

type ListenerFunc func(e Event) error

func (fn ListenerFunc) Handle(e Event) error {
	return fn(e)
}

type ListenerItem struct {
	Priority int
	Listener Listener
}

type ListenerQueue struct {
	items []*ListenerItem
}

func (lq *ListenerQueue) Len() int {
	return len(lq.items)
}

func (lq *ListenerQueue) IsEmpty() bool {
	return len(lq.items) == 0
}

func (lq *ListenerQueue) Push(item *ListenerItem) *ListenerQueue {
	lq.items = append(lq.items, item)
	return lq
}

func (lq *ListenerQueue) Sort() *ListenerQueue {
	items := SortListenerItems(lq.items)
	if !sort.IsSorted(items) {
		sort.Sort(items)
	}
	lq.items = items
	return lq
}

func (lq *ListenerQueue) Items() []*ListenerItem {
	return lq.items
}

func (lq *ListenerQueue) Remove(listener Listener) {
	if listener == nil {
		return
	}

	ptrVal := fmt.Sprintf("%p", listener)
	var newItems []*ListenerItem
	for _, item := range lq.items {
		itemPtrVal := fmt.Sprintf("%p", item.Listener)
		if itemPtrVal == ptrVal {
			continue
		}
		newItems = append(newItems, item)
	}
	lq.items = newItems
}

func (lq *ListenerQueue) Clear() {
	lq.items = make([]*ListenerItem, 0)
}

package memoryM

type heapStruct struct {
	store []*element
}

func (h *heapStruct) Len() int {
	return len(h.store)
}

func (h *heapStruct) Less(i, j int) bool {
	return h.store[i].ttl < h.store[j].ttl
}

func (h *heapStruct) Swap(i, j int) {
	h.store[i], h.store[j] = h.store[j], h.store[i]
	h.store[i].pos, h.store[j].pos = j, j
}

func (h *heapStruct) Push(e interface{}) {
	h.store = append(h.store, e.(*element))
}

func (h *heapStruct) Pop() interface{} {
	n := len(h.store)
	he := h.store[n-1]
	h.store = h.store[0 : n-1]
	return he
}

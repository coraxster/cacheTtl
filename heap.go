package memoryM

type heapStruct struct {
	store []*heapEl
	dict  map[string]int // we need a dict to quickly find element in heap
}

type heapEl struct {
	key string
	ttl int64
}

func (h *heapStruct) Len() int {
	return len(h.store)
}

func (h *heapStruct) Less(i, j int) bool {
	return h.store[i].ttl < h.store[j].ttl
}

func (h *heapStruct) Swap(i, j int) {
	h.store[i], h.store[j] = h.store[j], h.store[i]
	h.dict[h.store[j].key] = j
	h.dict[h.store[i].key] = i
}

func (h *heapStruct) Push(e interface{}) {
	el := e.(*heapEl)
	h.dict[el.key] = len(h.dict)
	h.store = append(h.store, el)
}

func (h *heapStruct) Pop() interface{} {
	n := len(h.store)
	he := h.store[n-1]
	h.store = h.store[0 : n-1]
	return he
}

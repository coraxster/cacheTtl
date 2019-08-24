package cacheTtl

type heapStruct []*element

func (h heapStruct) Len() int {
	return len(h)
}

func (h heapStruct) Less(i, j int) bool {
	return h[i].ttl.Before(h[j].ttl)
}

func (h heapStruct) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index, h[j].index = i, j
}

func (h *heapStruct) Push(e interface{}) {
	el := e.(*element)
	el.index = len(*h)
	*h = append(*h, el)
}

func (h *heapStruct) Pop() interface{} {
	old := *h
	n := len(old)
	el := old[n-1]
	*h = old[0 : n-1]
	el.index = -1
	return el
}

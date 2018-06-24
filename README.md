# memMan

Simple KV. With TTL.

Fast GC in top of heap struct.

```
2000000 Set:        556 ns/op
30000000 Get:       44.3 ns/op
2000000 SetGet:     615 ns/op
5000000 SimpleGC:   484.402861ms
5000000 AdvancedGC: 21.079µs
```
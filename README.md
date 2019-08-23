# cacheTtl

Simple KV. With TTL.

Fast GC on top of heap struct.

```
2000000 Set:        372 ns/op
30000000 Get:       43.6 ns/op
2000000 SetGet:     411 ns/op
5000000 SimpleGC:   469.402861ms
5000000 AdvancedGC: 20.079µs
```
# Byte Slice Interval Search Tree

Thread safe interval search tree with byte slice keys optimized for point query read performance.

[![Go Reference](https://godoc.org/github.com/logbn/byteinterval?status.svg)](https://godoc.org/github.com/logbn/byteinterval)
[![Go Report Card](https://goreportcard.com/badge/github.com/logbn/byteinterval?1)](https://goreportcard.com/report/github.com/logbn/byteinterval)
[![Go Coverage](https://github.com/logbn/byteinterval/wiki/coverage.svg)](https://raw.githack.com/wiki/logbn/byteinterval/coverage.html)

This is a convenience wrapper around a red/black interval tree [github.com/rdleal/intervalst/interval](https://pkg.go.dev/github.com/rdleal/intervalst/interval)

```go
tree := bytesinterval.New[int]()
tree.Insert([]byte(`alpha`), []byte(`bravo`), 100)
tree.Insert([]byte(`bravo`), []byte(`charlie`), 200)
c := tree.Insert([]byte(`bravo`), []byte(`delta`), 300)

items := tree.Find([]byte(`bravo`))
require.Equal(t, items, []int{200, 300})

items = tree.Find([]byte(`charlie`))
require.Equal(t, items, []int{300})

items = tree.FindAny([]byte(`alpha`), []byte(`bravo`), []byte(`charlie`))
require.Equal(t, items, []int{100, 200, 300})

c.Remove()

items = tree.Find([]byte(`bravo`))
require.Equal(t, items, []int{200})

items = tree.Find([]byte(`charlie`))
require.Equal(t, items, []int(nil))
```

## Performance

```
> make bench

goos: linux
goarch: amd64
pkg: github.com/logbn/byteinterval
cpu: 12th Gen Intel(R) Core(TM) i5-12600K
BenchmarkInsert/Overlap_Small/256-16         16165           74602   ns/op         35473 B/op        768 allocs/op
BenchmarkInsert/Overlap_Small/4096-16          685         1764461   ns/op        485062 B/op      12305 allocs/op
BenchmarkInsert/Overlap_Small/65536-16          20        51068256   ns/op       7916763 B/op     212992 allocs/op
BenchmarkInsert/Overlap_Big/256-16           16857           78283   ns/op         35026 B/op        768 allocs/op
BenchmarkInsert/Overlap_Big/4096-16            655         1757425   ns/op        489269 B/op      12306 allocs/op
BenchmarkInsert/Overlap_Big/65536-16            21        51194526   ns/op       7839382 B/op     212212 allocs/op
BenchmarkFind/256-16                       1683768             713.5 ns/op            40 B/op          3 allocs/op
BenchmarkFind/4096-16                      1688766             706.9 ns/op            40 B/op          3 allocs/op
BenchmarkFind/65536-16                     1674928             726.0 ns/op            40 B/op          3 allocs/op
BenchmarkFindAny/256-16                     276043            4214   ns/op          1385 B/op         53 allocs/op
BenchmarkFindAny/4096-16                    191158            6378   ns/op          1385 B/op         53 allocs/op
BenchmarkFindAny/65536-16                   102261           11433   ns/op          1384 B/op         53 allocs/op
PASS
ok      github.com/logbn/byteinterval   14.620s
```


## Random stuff

An fp-growth and improved fp-growth algorithm written in Go.

Test:

```
export GOGC=off

go test -v -run=TestBenchmarkMin* -timeout 99999s .

cd cmd
chokidar --initial "*.go" -c "go test ."


chokidar --initial "*.go" -c "go test -v -run=TestBenchmarkMin* -timeout 99999s ."
```


go test -run=TestBench* -v

Run profile:

```
go test -c
cmd.test.exe -test.memprofile=mem.prof -test.run=TestBench* -test.v
cmd.test.exe -test.cpuprofile=cpu.prof -test.run=TestBenchMarkMining* -test.v
go tool pprof -text --alloc_space cmd.test.exe mem.prof
go tool pprof -text cmd.test.exe cpu.prof
```

osx:

```
go test -c
./cmd.test -test.memprofile=mem.prof -test.run=TestBenchmarkMining* -test.v
./cmd.test -test.cpuprofile=cpu.prof -test.run=TestBenchmarkMining* -test.v
go tool pprof -text --alloc_objects cmd.test mem.prof
go tool pprof -text cmd.test cpu.prof
go tool pprof -svg --alloc_objects cmd.test mem.prof > mem.svg
go tool pprof -svg cmd.test cpu.prof > cpu.svg
go tool pprof -png --alloc_objects cmd.test mem.prof > mem.png
go tool pprof -png cmd.test cpu.prof > cpu.png
```


## Results 17.1.2017


### Observation: Results are random

```
"minsup=0.000001/transactions=10000000/itemspertx=5/minItems=10": 2.5707416,
"minsup=0.000010/transactions=10000000/itemspertx=5/minItems=100": 2.6412529,
"minsup=0.000100/transactions=10000000/itemspertx=5/minItems=1000": 2.550684,
"minsup=0.000500/transactions=10000000/itemspertx=5/minItems=5000": 2.5707453,
"minsup=0.001000/transactions=10000000/itemspertx=5/minItems=10000": 2.753271,
"minsup=0.005000/transactions=10000000/itemspertx=5/minItems=50000": 3.7924038,
"minsup=0.010000/transactions=10000000/itemspertx=5/minItems=100000": 2.345501,
"minsup=0.020000/transactions=10000000/itemspertx=5/minItems=200000": 2.7109971,
"minsup=0.030000/transactions=10000000/itemspertx=5/minItems=300000": 2.6995032,
"minsup=0.400000/transactions=10000000/itemspertx=5/minItems=4000000": 2.4863225,
"minsup=0.500000/transactions=10000000/itemspertx=5/minItems=5000000": 2.4159714
```

Most probable cause: too much memory allocation going on.

#### Trace

The database is taking up the most space. And ordering that database is costing the most space allocation too:

```
go tool pprof -text --alloc_space cmd.test.exe mem.prof
6.52GB of 6.52GB total (  100%)
Dropped 18 nodes (cum <= 0.03GB)
      flat  flat%   sum%        cum   cum%
    3.27GB 50.06% 50.06%     3.27GB 50.06%  github.com/arekkas/fp-growth/cmd.generateDatabase
    3.26GB 49.92%   100%     3.26GB 49.92%  github.com/arekkas/fp-growth/cmd.OrderItems
         0     0%   100%     3.26GB 49.92%  github.com/arekkas/fp-growth/cmd.Mine
         0     0%   100%     3.27GB 50.06%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining
         0     0%   100%     3.26GB 49.92%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining.func1
         0     0%   100%     6.52GB   100%  runtime.goexit
         0     0%   100%     6.52GB   100%  testing.tRunner
```

Apparently, OrderItems also takes up the most cpu complexitym, and specifically `growslice`:

```
go tool pprof -text cmd.test.exe cpu.prof
27.76s of 28.68s total (96.79%)
Dropped 70 nodes (cum <= 0.14s)
      flat  flat%   sum%        cum   cum%
     3.84s 13.39% 13.39%      4.20s 14.64%  runtime.mallocgc
     3.31s 11.54% 24.93%      7.13s 24.86%  github.com/arekkas/fp-growth/cmd.OrderItems
     2.41s  8.40% 33.33%      7.33s 25.56%  runtime.growslice
     2.24s  7.81% 41.14%     10.58s 36.89%  github.com/arekkas/fp-growth/cmd.generateDatabase
     1.85s  6.45% 47.59%      3.71s 12.94%  runtime.scanobject
     1.53s  5.33% 52.93%      1.53s  5.33%  math/rand.(*rngSource).Int63
     1.48s  5.16% 58.09%      2.50s  8.72%  runtime.mapassign1
     1.24s  4.32% 62.41%      1.24s  4.32%  github.com/arekkas/fp-growth/cmd.buildTree
     1.10s  3.84% 66.25%      4.62s 16.11%  github.com/arekkas/fp-growth/cmd.NewHeadTable
     1.01s  3.52% 69.77%      1.02s  3.56%  runtime.mapaccess1_fast64
     0.99s  3.45% 73.22%      0.99s  3.45%  runtime.memmove
     0.89s  3.10% 76.32%      3.58s 12.48%  math/rand.(*Rand).Int31n
     0.72s  2.51% 78.84%      0.72s  2.51%  runtime.heapBitsForObject
     0.67s  2.34% 81.17%      4.25s 14.82%  math/rand.(*Rand).Intn
     0.62s  2.16% 83.33%      0.62s  2.16%  runtime/internal/atomic.Or8
     0.61s  2.13% 85.46%      1.19s  4.15%  runtime.greyobject
     0.58s  2.02% 87.48%      2.69s  9.38%  math/rand.(*Rand).Int31
     0.58s  2.02% 89.50%      2.11s  7.36%  math/rand.(*Rand).Int63
     0.48s  1.67% 91.18%      0.48s  1.67%  runtime.aeshash64
     0.41s  1.43% 92.61%      0.41s  1.43%  runtime.memclr
     0.22s  0.77% 93.38%      0.36s  1.26%  runtime.typedmemmove
     0.18s  0.63% 94.00%      1.42s  4.95%  github.com/arekkas/fp-growth/cmd.NewFPTree
     0.18s  0.63% 94.63%      0.18s  0.63%  runtime.memequal64
     0.18s  0.63% 95.26%      0.39s  1.36%  runtime.newobject
     0.16s  0.56% 95.82%      0.16s  0.56%  runtime.heapBitsSetType
     0.07s  0.24% 96.06%      0.53s  1.85%  runtime.(*mcentral).cacheSpan
     0.06s  0.21% 96.27%      0.24s  0.84%  runtime.gopreempt_m
     0.06s  0.21% 96.48%      0.93s  3.24%  runtime.systemstack
     0.03s   0.1% 96.58%      0.31s  1.08%  runtime.newstack
     0.02s  0.07% 96.65%      0.24s  0.84%  runtime.(*mheap).alloc_m
     0.01s 0.035% 96.69%      0.17s  0.59%  runtime.(*mspan).sweep
     0.01s 0.035% 96.72%      0.18s  0.63%  runtime.goschedImpl
     0.01s 0.035% 96.76%      0.32s  1.12%  runtime.morestack
     0.01s 0.035% 96.79%      0.20s   0.7%  runtime.sweepone
         0     0% 96.79%     13.17s 45.92%  github.com/arekkas/fp-growth/cmd.Mine
         0     0% 96.79%     10.58s 36.89%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining
         0     0% 96.79%     13.17s 45.92%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining.func1
         0     0% 96.79%      0.53s  1.85%  runtime.(*mcache).nextFree.func1
         0     0% 96.79%      0.53s  1.85%  runtime.(*mcache).refill
         0     0% 96.79%      0.35s  1.22%  runtime.(*mcentral).grow
         0     0% 96.79%      0.37s  1.29%  runtime.(*mheap).alloc
         0     0% 96.79%      0.24s  0.84%  runtime.(*mheap).alloc.func1
         0     0% 96.79%      3.68s 12.83%  runtime.gcBgMarkWorker
         0     0% 96.79%      3.68s 12.83%  runtime.gcDrain
         0     0% 96.79%     27.43s 95.64%  runtime.goexit
         0     0% 96.79%      0.16s  0.56%  runtime.makeslice
         0     0% 96.79%      0.87s  3.03%  runtime.startTheWorldWithSema
         0     0% 96.79%     23.75s 82.81%  testing.tRunner
```

### Try: Removing append logic

```
func OrderItems(d DataSet, h HeadTable) DataSet {
	ds := make(DataSet, len(d))
	for x, dd := range d {
		items := make([]int, len(dd))
		i := 0
		for _, ih := range h {
			for _, id := range dd {
				if id == ih.Item {
					items[i] = id
					i++
				}
			}
		}
		items = items[:i]
		ds[x] = items
	}
	return ds
}
```

#### Results still random

```
"minsup=0.000001/transactions=5000000/itemspertx=5/minItems=5": 0.8379999,
"minsup=0.000010/transactions=5000000/itemspertx=5/minItems=50": 0.9544992000000001,
"minsup=0.000100/transactions=5000000/itemspertx=5/minItems=500": 0.8164994000000001,
"minsup=0.000500/transactions=5000000/itemspertx=5/minItems=2500": 1.5150001,
"minsup=0.001000/transactions=5000000/itemspertx=5/minItems=5000": 0.9224999,
"minsup=0.005000/transactions=5000000/itemspertx=5/minItems=25000": 0.9906296000000001,
"minsup=0.010000/transactions=5000000/itemspertx=5/minItems=50000": 0.9079985,
"minsup=0.020000/transactions=5000000/itemspertx=5/minItems=100000": 0.8584997000000001,
"minsup=0.030000/transactions=5000000/itemspertx=5/minItems=150000": 0.9295004,
"minsup=0.400000/transactions=5000000/itemspertx=5/minItems=2000000": 0.8699992000000001,
"minsup=0.500000/transactions=5000000/itemspertx=5/minItems=2500000": 0.8775008000000001
```

```
5.58GB of 5.58GB total (99.94%)
Dropped 27 nodes (cum <= 0.03GB)
      flat  flat%   sum%        cum   cum%
    3.19GB 57.13% 57.13%     3.19GB 57.13%  github.com/arekkas/fp-growth/cmd.generateDatabase
    2.39GB 42.81% 99.94%     2.39GB 42.81%  github.com/arekkas/fp-growth/cmd.OrderItems
         0     0% 99.94%     2.39GB 42.85%  github.com/arekkas/fp-growth/cmd.Mine
         0     0% 99.94%     3.19GB 57.13%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining
         0     0% 99.94%     2.39GB 42.85%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining.func1
         0     0% 99.94%     5.58GB   100%  runtime.goexit
         0     0% 99.94%     5.58GB   100%  testing.tRunner
         
23.51s of 24.49s total (96.00%)
Dropped 68 nodes (cum <= 0.12s)
      flat  flat%   sum%        cum   cum%
     2.57s 10.49% 10.49%      2.92s 11.92%  runtime.mallocgc
     2.13s  8.70% 19.19%      4.13s 16.86%  runtime.scanobject
     2.01s  8.21% 27.40%      9.93s 40.55%  github.com/arekkas/fp-growth/cmd.generateDatabase
     1.98s  8.08% 35.48%      3.15s 12.86%  github.com/arekkas/fp-growth/cmd.OrderItems
     1.52s  6.21% 41.69%      1.52s  6.21%  math/rand.(*rngSource).Int63
     1.42s  5.80% 47.49%      4.59s 18.74%  github.com/arekkas/fp-growth/cmd.NewHeadTable
     1.40s  5.72% 53.21%      2.35s  9.60%  runtime.mapassign1
     1.23s  5.02% 58.23%      1.23s  5.02%  github.com/arekkas/fp-growth/cmd.buildTree
     0.86s  3.51% 61.74%      0.87s  3.55%  runtime.heapBitsForObject
     0.85s  3.47% 65.21%      3.16s 12.90%  runtime.growslice
     0.82s  3.35% 68.56%      0.82s  3.35%  runtime.mapaccess1_fast64
     0.81s  3.31% 71.87%      3.67s 14.99%  math/rand.(*Rand).Int31n
     0.73s  2.98% 74.85%      4.40s 17.97%  math/rand.(*Rand).Intn
     0.68s  2.78% 77.62%      2.86s 11.68%  math/rand.(*Rand).Int31
     0.66s  2.69% 80.32%      2.18s  8.90%  math/rand.(*Rand).Int63
     0.61s  2.49% 82.81%      1.18s  4.82%  runtime.greyobject
     0.61s  2.49% 85.30%      0.61s  2.49%  runtime/internal/atomic.Or8
     0.57s  2.33% 87.63%      0.57s  2.33%  runtime.memmove
     0.51s  2.08% 89.71%      0.51s  2.08%  runtime.aeshash64
     0.36s  1.47% 91.18%      0.36s  1.47%  runtime.memclr
     0.19s  0.78% 91.96%      0.31s  1.27%  runtime.typedmemmove
     0.14s  0.57% 92.53%      1.37s  5.59%  github.com/arekkas/fp-growth/cmd.NewFPTree
     0.14s  0.57% 93.10%      0.26s  1.06%  runtime.newobject
     0.13s  0.53% 93.63%      0.13s  0.53%  runtime.heapBitsSetType
     0.13s  0.53% 94.16%      0.13s  0.53%  runtime.memequal64
     0.09s  0.37% 94.53%      1.22s  4.98%  runtime.makeslice
     0.07s  0.29% 94.81%      0.14s  0.57%  runtime.shade
     0.06s  0.24% 95.06%      0.28s  1.14%  runtime.newstack
     0.05s   0.2% 95.26%      0.92s  3.76%  runtime.systemstack
     0.03s  0.12% 95.39%      0.20s  0.82%  runtime.(*mspan).sweep
     0.03s  0.12% 95.51%      0.17s  0.69%  runtime.gcmarkwb_m
     0.03s  0.12% 95.63%      0.26s  1.06%  runtime.sweepone
     0.02s 0.082% 95.71%      0.46s  1.88%  runtime.(*mcentral).cacheSpan
     0.02s 0.082% 95.79%      0.19s  0.78%  runtime.goschedImpl
     0.02s 0.082% 95.88%      0.13s  0.53%  runtime.schedule
     0.01s 0.041% 95.92%      0.18s  0.73%  runtime.gopreempt_m
     0.01s 0.041% 95.96%      0.29s  1.18%  runtime.morestack
     0.01s 0.041% 96.00%      0.18s  0.73%  runtime.writebarrierptr_nostore1.func1
         0     0% 96.00%      9.11s 37.20%  github.com/arekkas/fp-growth/cmd.Mine
         0     0% 96.00%      9.94s 40.59%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining
         0     0% 96.00%      9.11s 37.20%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining.func1
         0     0% 96.00%      0.46s  1.88%  runtime.(*mcache).nextFree.func1
         0     0% 96.00%      0.46s  1.88%  runtime.(*mcache).refill
         0     0% 96.00%      0.34s  1.39%  runtime.(*mcentral).grow
         0     0% 96.00%      0.29s  1.18%  runtime.(*mheap).alloc
         0     0% 96.00%      0.24s  0.98%  runtime.(*mheap).alloc.func1
         0     0% 96.00%      0.23s  0.94%  runtime.(*mheap).alloc_m
         0     0% 96.00%      4.21s 17.19%  runtime.gcBgMarkWorker
         0     0% 96.00%      4.21s 17.19%  runtime.gcDrain
         0     0% 96.00%     23.26s 94.98%  runtime.goexit
         0     0% 96.00%      0.15s  0.61%  runtime.gosweepone.func1
         0     0% 96.00%      0.87s  3.55%  runtime.startTheWorldWithSema
         0     0% 96.00%     19.05s 77.79%  testing.tRunner
```

### Try: Replaced with proper sort solution

Next, a proper sort solution was used (the one provided by golang), not allocating a new database. This results in better
mem / cpu profiles but has equally unsatisfying results

```
"minsup=0.000001/transactions=5000000/itemspertx=5/minItems=5": 2.034684,
"minsup=0.000010/transactions=5000000/itemspertx=5/minItems=50": 2.2116977,
"minsup=0.000100/transactions=5000000/itemspertx=5/minItems=500": 2.8525002,
"minsup=0.000500/transactions=5000000/itemspertx=5/minItems=2500": 2.0490009,
"minsup=0.001000/transactions=5000000/itemspertx=5/minItems=5000": 2.1565008,
"minsup=0.005000/transactions=5000000/itemspertx=5/minItems=25000": 3.0685119,
"minsup=0.010000/transactions=5000000/itemspertx=5/minItems=50000": 2.5010006000000002,
"minsup=0.020000/transactions=5000000/itemspertx=5/minItems=100000": 2.3720007,
"minsup=0.030000/transactions=5000000/itemspertx=5/minItems=150000": 2.2644998,
"minsup=0.400000/transactions=5000000/itemspertx=5/minItems=2000000": 2.5505017,
"minsup=0.500000/transactions=5000000/itemspertx=5/minItems=2500000": 2.4024993
```

```

11.96GB of 11.96GB total (  100%)
Dropped 5 nodes (cum <= 0.06GB)
      flat  flat%   sum%        cum   cum%
    4.12GB 34.43% 34.43%     4.12GB 34.43%  github.com/arekkas/fp-growth/cmd.HeadTable.Get
    3.31GB 27.71% 62.14%     3.31GB 27.71%  github.com/arekkas/fp-growth/cmd.generateDatabase
    2.47GB 20.69% 82.82%     8.64GB 72.28%  github.com/arekkas/fp-growth/cmd.OrderItems
    2.05GB 17.16%   100%     6.17GB 51.59%  github.com/arekkas/fp-growth/cmd.filterItems
         0     0%   100%     8.64GB 72.28%  github.com/arekkas/fp-growth/cmd.Mine
         0     0%   100%     4.12GB 34.43%  github.com/arekkas/fp-growth/cmd.OrderItems.func1
         0     0%   100%     3.31GB 27.71%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining
         0     0%   100%     8.65GB 72.29%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining.func1
         0     0%   100%    11.96GB   100%  runtime.goexit
         0     0%   100%    11.96GB   100%  testing.tRunner


40.89s of 43.13s total (94.81%)
Dropped 88 nodes (cum <= 0.22s)
      flat  flat%   sum%        cum   cum%
     7.14s 16.55% 16.55%      9.52s 22.07%  runtime.mallocgc
     2.77s  6.42% 22.98%      6.71s 15.56%  runtime.scanobject
     2.27s  5.26% 28.24%     10.02s 23.23%  github.com/arekkas/fp-growth/cmd.generateDatabase
     2.06s  4.78% 33.02%      7.51s 17.41%  runtime.growslice
     1.73s  4.01% 37.03%      1.73s  4.01%  runtime.heapBitsSetType
     1.71s  3.96% 40.99%      1.71s  3.96%  runtime.heapBitsForObject
     1.67s  3.87% 44.86%      1.67s  3.87%  github.com/arekkas/fp-growth/cmd.buildTree
     1.64s  3.80% 48.67%      4.95s 11.48%  github.com/arekkas/fp-growth/cmd.NewHeadTable
     1.60s  3.71% 52.38%      2.48s  5.75%  runtime.mapassign1
     1.53s  3.55% 55.92%      5.66s 13.12%  github.com/arekkas/fp-growth/cmd.HeadTable.Get
     1.48s  3.43% 59.36%      1.48s  3.43%  math/rand.(*rngSource).Int63
     1.26s  2.92% 62.28%      1.26s  2.92%  runtime.memmove
     1.19s  2.76% 65.04%      2.25s  5.22%  runtime.greyobject
     1.07s  2.48% 67.52%      1.07s  2.48%  runtime/internal/atomic.Or8
     0.98s  2.27% 69.79%      0.98s  2.27%  runtime.memclr
     0.96s  2.23% 72.01%      1.63s  3.78%  github.com/arekkas/fp-growth/cmd.OrderItemsByHeaderTableWrapper.Less
     0.92s  2.13% 74.15%     11.04s 25.60%  github.com/arekkas/fp-growth/cmd.filterItems
     0.83s  1.92% 76.07%      0.83s  1.92%  runtime.mapaccess1_fast64
     0.80s  1.85% 77.93%      3.42s  7.93%  math/rand.(*Rand).Int31n
     0.67s  1.55% 79.48%      0.67s  1.55%  github.com/arekkas/fp-growth/cmd.HeadTable.GetPosition
     0.62s  1.44% 80.92%      4.04s  9.37%  math/rand.(*Rand).Intn
     0.59s  1.37% 82.29%      2.07s  4.80%  math/rand.(*Rand).Int63
     0.55s  1.28% 83.56%      2.62s  6.07%  math/rand.(*Rand).Int31
     0.55s  1.28% 84.84%      5.66s 13.12%  runtime.newobject
     0.48s  1.11% 85.95%     17.17s 39.81%  github.com/arekkas/fp-growth/cmd.OrderItems
     0.43s     1% 86.95%      2.13s  4.94%  github.com/arekkas/fp-growth/cmd.(*OrderItemsByHeaderTableWrapper).Less
     0.42s  0.97% 87.92%      2.78s  6.45%  sort.insertionSort
     0.39s   0.9% 88.82%      0.39s   0.9%  runtime.aeshash64
     0.36s  0.83% 89.66%      6.02s 13.96%  github.com/arekkas/fp-growth/cmd.OrderItems.func1
     0.34s  0.79% 90.45%      3.30s  7.65%  sort.Sort
     0.32s  0.74% 91.19%      0.32s  0.74%  runtime.duffcopy
     0.25s  0.58% 91.77%      0.90s  2.09%  runtime.typedmemmove
     0.23s  0.53% 92.30%      0.23s  0.53%  github.com/arekkas/fp-growth/cmd.(*OrderItemsByHeaderTableWrapper).Swap
     0.21s  0.49% 92.79%      1.88s  4.36%  github.com/arekkas/fp-growth/cmd.NewFPTree
     0.12s  0.28% 93.07%      0.24s  0.56%  runtime.unlock
     0.11s  0.26% 93.32%      0.39s   0.9%  runtime.schedule
     0.11s  0.26% 93.58%      2.89s  6.70%  sort.quickSort
     0.10s  0.23% 93.81%      2.10s  4.87%  runtime.convT2I
     0.09s  0.21% 94.02%      0.24s  0.56%  runtime.lock
     0.07s  0.16% 94.18%      1.82s  4.22%  runtime.systemstack
     0.06s  0.14% 94.32%      0.56s  1.30%  runtime.newstack
     0.06s  0.14% 94.46%      0.45s  1.04%  runtime.sweepone
     0.05s  0.12% 94.57%      0.95s  2.20%  runtime.(*mcentral).cacheSpan
     0.02s 0.046% 94.62%      0.35s  0.81%  runtime.(*mheap).alloc_m
     0.02s 0.046% 94.67%      0.52s  1.21%  runtime.goschedImpl
     0.01s 0.023% 94.69%      0.61s  1.41%  runtime.(*mcentral).grow
     0.01s 0.023% 94.71%      0.36s  0.83%  runtime.(*mheap).alloc.func1
     0.01s 0.023% 94.74%      0.38s  0.88%  runtime.(*mspan).sweep
     0.01s 0.023% 94.76%      6.56s 15.21%  runtime.gcDrain
     0.01s 0.023% 94.78%      0.41s  0.95%  runtime.gopreempt_m
     0.01s 0.023% 94.81%      0.57s  1.32%  runtime.morestack
         0     0% 94.81%     24.01s 55.67%  github.com/arekkas/fp-growth/cmd.Mine
         0     0% 94.81%     10.03s 23.26%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining
         0     0% 94.81%     24.02s 55.69%  github.com/arekkas/fp-growth/cmd.TestBenchmarkMining.func1

```

### Observation: min_sup = 1 is very memory heavy

... even with low # of transactions

... the number of items pers transaction has a high impact on memory allocation

### Observation: support counts are iterated with too large gaps

manually iterating minItems helped and finally yielded ok-ish results!

### FINALLY RESULTS

```
"minsup=0.030000/transactions=1000/itemspertx=8/minItems=30": 58.5780143,
"minsup=0.031000/transactions=1000/itemspertx=8/minItems=31": 41.0704476,
"minsup=0.032000/transactions=1000/itemspertx=8/minItems=32": 20.0553367,
"minsup=0.033000/transactions=1000/itemspertx=8/minItems=33": 7.2605179,
"minsup=0.034000/transactions=1000/itemspertx=8/minItems=34": 5.5139975,
"minsup=0.035000/transactions=1000/itemspertx=8/minItems=35": 3.7444995,
"minsup=0.036000/transactions=1000/itemspertx=8/minItems=36": 2.4385017,
"minsup=0.037000/transactions=1000/itemspertx=8/minItems=37": 1.1875009,
"minsup=0.038000/transactions=1000/itemspertx=8/minItems=38": 0.9679995,
"minsup=0.039000/transactions=1000/itemspertx=8/minItems=39": 1.1464987,
"minsup=0.040000/transactions=1000/itemspertx=8/minItems=40": 0.38200080000000003,
"minsup=0.045000/transactions=1000/itemspertx=8/minItems=45": 0.1854994,
"minsup=0.050000/transactions=1000/itemspertx=8/minItems=50": 0.054498500000000005,
"minsup=0.100000/transactions=1000/itemspertx=8/minItems=100": 0.0015021000000000001,
"minsup=0.500000/transactions=1000/itemspertx=8/minItems=500": 0.001
```

### 18.1.2017

The basic "Support Count Two-dimensional Table" algorithm is working. It's further improved by leabing out 0 fields, but
it does not filter infrequent rows

### Results, at last!

Yup, that thing has definitely a bad performance! Why? Because it's super memory intensive

```
"algo=improved/minsup=0.050000/transactions=1000/items=8/minItems=50": 8.451999,
"algo=improved/minsup=0.070000/transactions=1000/items=8/minItems=70": 0.26749860000000003,
"algo=improved/minsup=0.080000/transactions=1000/items=8/minItems=80": 0.3689983,
"algo=improved/minsup=0.100000/transactions=1000/items=8/minItems=100": 0.3499987,
"algo=improved/minsup=0.200000/transactions=1000/items=8/minItems=200": 0.002,
"algo=improved/minsup=0.500000/transactions=1000/items=8/minItems=500": 0.0004998,
"algo=original/minsup=0.050000/transactions=1000/items=8/minItems=50": 0.045996,
"algo=original/minsup=0.070000/transactions=1000/items=8/minItems=70": 0.0430006,
"algo=original/minsup=0.080000/transactions=1000/items=8/minItems=80": 0.048999100000000004,
"algo=original/minsup=0.100000/transactions=1000/items=8/minItems=100": 0.0020008,
"algo=original/minsup=0.200000/transactions=1000/items=8/minItems=200": 0.0009987000000000002,
"algo=original/minsup=0.500000/transactions=1000/items=8/minItems=500": 0.0004997000000000001
```

Macbook pro:

```
```


### the garbage collector is taking a lot of time
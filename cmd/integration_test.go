package cmd

import (
	"testing"
	"fmt"
	"encoding/json"
	"time"
)

func TestMining(t *testing.T) {
	minSup := 2
	db := exDB
	t.Logf("%v", Mine(db, minSup))
}


func TestBenchmarkMining(t *testing.T) {
	// t.SkipNow()
	times := map[string]float64{}
	minSups := []int{
		//0.4,
		//0.2,
		//0.1,
		//0.01,
		//0.0001,
		//0.000001,
		500,
		100,
		50,
		45,
		40,
		39,
		38,
		37,
		36,
		35,
		34,
		33,
		32,
		31,
		30,
	}
	txss := []int{
		1000,
	}
	uss := []int{
		8,
	}

	var dbs = map[int]map[int]map[int]DataSet{}
	dbstart := time.Now()
	for _, minSup := range minSups {
		for _, txs := range txss {
			for _, us := range uss {
				dbs[minSup] = make(map[int]map[int]DataSet)
				dbs[minSup][txs] = make(map[int]DataSet)
				dbs[minSup][txs][us], _ = generateDatabase(txs, us)
			}
		}
	}
	dbend := time.Now()

	for _, minSup := range minSups {
		for _, txs := range txss {
			for _, us := range uss {
				//var minItems = int(float32(txs) * minSup)
				var procSup = float32(minSup) / float32(txs)
				var minItems = minSup
				db := dbs[minSup][txs][us]
				// var patterns []Pattern
				d := fmt.Sprintf("minsup=%f/transactions=%d/items=%d/minItems=%d", procSup, txs, us, minItems)
				t.Run(d, func(t *testing.T) {
					if minItems < 2 {
						t.SkipNow()
					}
					start := time.Now()
					_ = Mine(db, minItems)

					end := time.Now()
					t.Logf("Iteration took %f seconds, %d ns", end.Sub(start).Seconds(), end.Sub(start).Nanoseconds())
					times[d] = end.Sub(start).Seconds()
				})
			}
		}
	}
	out, _ := json.MarshalIndent(times, "", "\t")
	t.Logf("Database generation took %f seconds", dbend.Sub(dbstart).Seconds())

	t.Logf("%s", string(out))
}

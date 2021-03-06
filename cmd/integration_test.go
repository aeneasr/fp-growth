package cmd

import (
	"testing"
	"fmt"
	"encoding/json"
	"time"
	"runtime"
	"runtime/debug"
)

func TestMining(t *testing.T) {
	minSup := 2
	db := exDB
	t.Logf("%v", Mine(db, minSup))
}
func runner(improved bool, dd string, db DataSet, minItems int, times map[string]float64) func (t *testing.T) {
	return func(t *testing.T) {
		if minItems < 2 {
			t.SkipNow()
		}
		start := time.Now()
		if !improved {
			_ = Mine(db, minItems)
		} else {
			_ = MineImproved(db, minItems)
		}

		end := time.Now()
		pp := fmt.Sprintf("Iteration took %f seconds, %d ns", end.Sub(start).Seconds(), end.Sub(start).Nanoseconds())
		fmt.Println(pp)
		t.Log(pp)
		times[dd] = end.Sub(start).Seconds()


		start = time.Now()
		runtime.GC()
		end = time.Now()
		fmt.Printf("GC took %f seconds\n", end.Sub(start).Seconds())
	}
}

func TestBenchmarkMining(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	debug.SetGCPercent(-1)
	times := map[string]float64{}
	minSups := []int{
		800,
		700,
		600,
		500,
		400,
		300,
		200,
		100,
		90,
		80,
		70,
		60,
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
		29,
		28,
		27,26,25,24,23,22,21,20,19,18,17,16,15,
	}
	txss := []int{		1000	}
	uss := []int{		40	}

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
				var procSup = float32(minSup) / float32(txs)
				var minItems = minSup
				d := fmt.Sprintf("algo=original/minsup=%f/transactions=%d/items=%d/minItems=%d", procSup, txs, us, minItems)
				t.Run(d, runner(false, d, dbs[minSup][txs][us], minItems, times))
				//if minSup >= 45{
					d = fmt.Sprintf("algo=improved/minsup=%f/transactions=%d/items=%d/minItems=%d", procSup, txs, us, minItems)
					t.Run(d, runner(true, d, dbs[minSup][txs][us], minItems, times))
				//}
			}
		}
	}
	out, _ := json.MarshalIndent(times, "", "\t")
	t.Logf("Database generation took %f seconds", dbend.Sub(dbstart).Seconds())

	t.Logf("%s", string(out))
}


func TestBenchmarkOriginalMining(t *testing.T) {
	times := map[string]float64{}
	txss := []int{		1000,	}
	uss := []int{		8,	}
	minSups := []int{
		500,
		200,
		100,
		80,
		70,
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
				var procSup = float32(minSup) / float32(txs)
				var minItems = minSup
				d := fmt.Sprintf("algo=original/minsup=%f/transactions=%d/items=%d/minItems=%d", procSup, txs, us, minItems)
				t.Run(d,runner(false, d, dbs[minSup][txs][us], minItems, times))
			}
		}
	}
	out, _ := json.MarshalIndent(times, "", "\t")
	t.Logf("Database generation took %f seconds", dbend.Sub(dbstart).Seconds())

	t.Logf("%s", string(out))
}

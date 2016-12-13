package cmd

import (
	"math/rand"
	"time"
)

func generateDatabase(txs, uniques int) (DataSet, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	db := make([]Items, txs)
	for i := 0; i < txs; i++ {
		db[i] = []int{}
		var should int
		for j := 1; j <= uniques; j++ {
			should = r.Intn(2)
			if should > 0 {
				db[i] = append(db[i], j)
			}
		}
		if len(db[i]) == 0 {
			db[i] = []int{r.Intn(uniques)+1}
		}
	}
	return db, nil
}